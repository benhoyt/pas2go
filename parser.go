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
	stmts := p.stmts()
	p.expect(END)
	return &CompoundStmt{stmts}
}

func (p *parser) stmts() []Stmt {
	stmts := []Stmt{p.stmt()}
	for p.tok == SEMICOLON {
		p.next()
		stmts = append(stmts, p.stmt())
	}
	return stmts
}

func (p *parser) stmt() Stmt {
	return p.labelledStmt(true)
}

func (p *parser) labelledStmt(allowLabel bool) Stmt {
	switch p.tok {
	case IDENT, AT:
		varExpr := p.varExpr()
		switch p.tok {
		case ASSIGN:
			p.next()
			value := p.expr()
			return &AssignStmt{varExpr, value}
		case COLON:
			if !varExpr.IsNameOnly() {
				panic(p.error("label must be a simple identifier"))
			}
			if !allowLabel {
				panic(p.error("unexpected label"))
			}
			p.next()
			stmt := p.labelledStmt(false)
			return &LabelledStmt{varExpr.Name, stmt}
		case LPAREN:
			if !varExpr.IsNameOnly() {
				panic(p.error("procedure name must be a simple identifier"))
			}
			p.next()
			args := p.argList()
			p.expect(RPAREN)
			return &ProcStmt{varExpr.Name, args}
		default:
			return &ProcStmt{varExpr.Name, nil}
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
		p.next()
		selector := p.expr()
		p.expect(OF)
		cases := []*CaseElement{p.caseElement()}
		var elseStmts []Stmt
		// Grammar quirkiness here, but this seems to mimic Turbo Pascal
		for p.tok == SEMICOLON || p.tok == ELSE {
			if p.tok == SEMICOLON {
				p.next()
			}
			if p.tok == END {
				break
			}
			if p.tok == ELSE {
				p.next()
				elseStmts = p.stmts()
				break
			}
			cases = append(cases, p.caseElement())
		}
		p.expect(END)
		return &CaseStmt{selector, cases, elseStmts}
	case WHILE:
		p.next()
		cond := p.expr()
		p.expect(DO)
		stmt := p.stmt()
		return &WhileStmt{cond, stmt}
	case REPEAT:
		p.next()
		stmts := p.stmts()
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
		vars := []*VarExpr{p.varExpr()}
		for p.tok == COMMA {
			p.next()
			vars = append(vars, p.varExpr())
		}
		p.expect(DO)
		stmt := p.stmt()
		return &WithStmt{vars, stmt}
	default:
		return &EmptyStmt{}
	}
}

func (p *parser) caseElement() *CaseElement {
	consts := []Expr{p.expr()} // should be p.constant()
	for p.tok == COMMA {
		p.next()
		consts = append(consts, p.expr())
	}
	p.expect(COLON)
	return &CaseElement{consts, p.stmt()}
}

func (p *parser) argList() []Expr {
	args := []Expr{p.expr()}
	for p.tok == COMMA {
		p.next()
		args = append(args, p.expr())
	}
	return args
}

// variable: (AT identifier | identifier) (LBRACKET expression (COMMA expression)* RBRACKET | DOT identifier | POINTER)*
func (p *parser) varExpr() *VarExpr {
	hasAt := false
	if p.tok == AT {
		p.next()
		hasAt = true
	}
	ident := p.val
	p.expect(IDENT)
	var parts []VarSuffix
	for p.tok == LBRACKET || p.tok == DOT || p.tok == POINTER {
		var part VarSuffix
		switch p.tok {
		case LBRACKET:
			p.next()
			indexes := []Expr{p.expr()}
			for p.tok == COMMA {
				p.next()
				indexes = append(indexes, p.expr())
			}
			p.expect(RBRACKET)
			part = &IndexSuffix{indexes}
		case DOT:
			p.next()
			field := p.val
			p.expect(IDENT)
			part = &DotSuffix{field}
		case POINTER:
			p.next()
			part = &PointerSuffix{}
		}
		parts = append(parts, part)
	}
	return &VarExpr{hasAt, ident, parts}
}

// expr: simpleExpr (relationalOp expr)?
func (p *parser) expr() Expr {
	return p.binaryExpr(p.simpleExpr, p.expr, EQUALS, NOT_EQUALS, LESS, LTE, GREATER, GTE, IN)
}

// simpleExpr: term (additiveOp simpleExpr)?
func (p *parser) simpleExpr() Expr {
	return p.binaryExpr(p.term, p.simpleExpr, PLUS, MINUS, OR)
}

// term: signedFactor (multiplicativeOp term)?
func (p *parser) term() Expr {
	return p.binaryExpr(p.signedFactor, p.term, STAR, SLASH, DIV, MOD, AND)
}

// signedFactor: (PLUS | MINUS)? factor
func (p *parser) signedFactor() Expr {
	if p.tok == PLUS || p.tok == MINUS {
		op := p.tok
		p.next()
		return &UnaryExpr{op, p.factor()}
	}
	return p.factor()
}

// factor: var | LPAREN expr RPAREN | function | constant | NOT factor | TRUE | FALSE
// TODO: handle constantChr
func (p *parser) factor() Expr {
	switch p.tok {
	case LPAREN:
		p.next()
		expr := p.expr()
		p.expect(RPAREN)
		return expr
	case NUM:
		val := p.val
		p.next()
		i, err := strconv.Atoi(val)
		if err != nil {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(p.error("invalid number: %s", err))
			}
			return &ConstExpr{f}
		}
		return &ConstExpr{i}
	case STR:
		s := p.val
		p.next()
		return &ConstExpr{s}
	case NOT:
		p.next()
		return &UnaryExpr{NOT, p.factor()}
	case TRUE:
		p.next()
		return &ConstExpr{true}
	case FALSE:
		p.next()
		return &ConstExpr{false}
	case IDENT, AT:
		varExpr := p.varExpr()
		switch p.tok {
		case LPAREN:
			if !varExpr.IsNameOnly() {
				panic(p.error("function name must be a simple identifier"))
			}
			p.next()
			args := p.argList()
			p.expect(RPAREN)
			return &FuncExpr{varExpr.Name, args}
		default:
			return varExpr
		}
	default:
		panic(p.error("expected factor"))
	}
}

func (p *parser) binaryExpr(left, right func() Expr, ops ...Token) Expr {
	expr := left()
	for p.matches(ops...) {
		op := p.tok
		p.next()
		rightExpr := right()
		expr = &BinaryExpr{expr, op, rightExpr}
	}
	return expr
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

// Return true iff current token matches one of the given operators,
// but don't parse next token.
func (p *parser) matches(operators ...Token) bool {
	for _, operator := range operators {
		if p.tok == operator {
			return true
		}
	}
	return false
}

// Format given string and args with Sprintf and return *ParseError
// with that message and the current position.
func (p *parser) error(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return &ParseError{p.pos, message}
}
