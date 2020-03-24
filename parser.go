// Turbo Pascal recursive descent parser

package main

import (
	"fmt"
	"strconv"
)

// ParseError (actually *ParseError) is the type of error returned by
// ParseProgram.
type ParseError struct {
	// Source line/column position where the error occurred.
	Position Position
	// Error message.
	Message string
}

// Error returns a formatted version of the error, including the line
// and column numbers.
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %d:%d: %s", e.Position.Line, e.Position.Column, e.Message)
}

// ParseProgram parses a single program, returning the *Program
// abstract syntax tree or a *ParseError on error.
func ParseProgram(src []byte) (prog *Program, err error) {
	defer func() {
		// The parser uses panic with a *ParseError to signal parsing
		// errors internally, and they're caught here. This
		// significantly simplifies the recursive descent calls as
		// we don't have to check errors everywhere.
		if r := recover(); r != nil {
			// Convert to ParseError or re-panic
			err = r.(*ParseError)
		}
	}()
	lexer := NewLexer(src)
	p := parser{lexer: lexer}
	p.next() // initialize p.tok
	return p.program(), nil
}

// Parser state
type parser struct {
	// Lexer instance and current token values
	lexer *Lexer
	pos   Position // position of last token (tok)
	tok   Token    // last lexed token
	val   string   // string value of last token (or "")
}

func (p *parser) program() *Program {
	program := &Program{}

	p.expect(PROGRAM)
	program.Name = p.val
	p.expect(IDENT)
	p.expect(SEMICOLON)

	program.Stmt = p.compoundStmt()
	p.expect(DOT)
	p.expect(EOF)

	return program
}

func (p *parser) compoundStmt() *CompoundStmt {
	p.expect(BEGIN)
	stmts := p.stmts(END)
	p.expect(END)
	return &CompoundStmt{stmts}
}

func (p *parser) stmts(end Token) []Stmt {
	var stmts []Stmt
	expectSemi := false
	for p.tok != end {
		if expectSemi {
			p.expect(SEMICOLON)
		}
		expectSemi = true
		stmts = append(stmts, p.stmt())
	}
	return stmts
}

func (p *parser) stmt() Stmt {
	return p.labelledStmt(true)
}

func (p *parser) labelledStmt(allowLabel bool) Stmt {
	// label: IDENT COLON unlabelledStmt
	// assignment: variable ASSIGN expr
	//   (AT identifier | identifier)
	//   (LBRACK expression (COMMA expression)* RBRACK |
	//    LBRACK2 expression (COMMA expression)* RBRACK2 |
	//    DOT identifier |
	//    POINTER)*
	// call: IDENT LPAREN paramList RPAREN
	// GOTO
	// BEGIN
	// IF
	// CASE
	// WHILE
	// REPEAT
	// FOR
	// WITH
	// empty
	switch p.tok {
	case IDENT, AT:
		// TODO: or use p.variable() as "covering grammar" and type switch
		// TODO: hasAt := p.tok == AT
		if p.tok == AT {
			p.next()
		}
		ident := p.val
		p.expect(IDENT)
		switch p.tok {
		case ASSIGN:
			p.next()
			value := p.expr()
			return &AssignStmt{&Variable{ident}, value}
		case COLON:
			if !allowLabel {
				panic(p.error("unexpected label"))
			}
			p.next()
			stmt := p.labelledStmt(false)
			return &LabelledStmt{ident, stmt}
		case LPAREN:
			p.next()
			args := p.argList()
			p.expect(RPAREN)
			return &ProcStmt{ident, args}
		default:
			return &ProcStmt{ident, nil}
		}
	case GOTO:
		p.next()
		label := p.val
		p.expect(IDENT)
		return &GotoStmt{label}
	case BEGIN:
		return p.compoundStmt()
	case IF:
		p.next()
		cond := p.expr()
		p.expect(THEN)
		then := p.stmt()
		var elseStmt Stmt
		if p.tok == ELSE {
			p.next()
			elseStmt = p.stmt()
		}
		return &IfStmt{cond, then, elseStmt}
	case CASE:
		return &CaseStmt{} // TODO
	case WHILE:
		p.next()
		cond := p.expr()
		p.expect(DO)
		stmt := p.stmt()
		return &WhileStmt{cond, stmt}
	case REPEAT:
		p.next()
		stmts := p.stmts(UNTIL)
		p.next()
		cond := p.expr()
		return &RepeatStmt{stmts, cond}
	case FOR:
		p.next()
		ident := p.val
		p.expect(IDENT)
		p.expect(ASSIGN)
		initial := p.expr()
		if p.tok != TO && p.tok != DOWNTO {
			panic(p.error("expected 'to' or 'downto'"))
		}
		down := p.tok == DOWNTO
		p.next()
		final := p.expr()
		p.expect(DO)
		stmt := p.stmt()
		return &ForStmt{ident, initial, down, final, stmt}
	case WITH:
		p.next()
		vars := []*Variable{p.variable()}
		for p.tok == COMMA {
			p.next()
			vars = append(vars, p.variable())
		}
		p.expect(DO)
		stmt := p.stmt()
		return &WithStmt{vars, stmt}
	default:
		return &EmptyStmt{}
	}
}

func (p *parser) argList() []Expr {
	args := []Expr{p.expr()}
	for p.tok == COMMA {
		p.next()
		args = append(args, p.expr())
	}
	return args
}

func (p *parser) variable() *Variable {
	ident := p.val
	p.expect(IDENT)
	return &Variable{ident}
}

func (p *parser) expr() Expr {
	val := p.val
	p.expect(NUM)
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(p.error("expected integer"))
	}
	return &IntegerExpr{i}
}

// Parse next token into p.tok (and set p.pos and p.val).
func (p *parser) next() {
	p.pos, p.tok, p.val = p.lexer.Scan()
	if p.tok == ILLEGAL {
		panic(p.error("%s", p.val))
	}
}

// Ensure current token is tok, and parse next token into p.tok.
func (p *parser) expect(tok Token) {
	if p.tok != tok {
		panic(p.error("expected %s instead of %s", tok, p.tok))
	}
	p.next()
}

// Format given string and args with Sprintf and return *ParseError
// with that message and the current position.
func (p *parser) error(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return &ParseError{p.pos, message}
}
