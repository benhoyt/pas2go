// Turbo Pascal recursive descent parser

package main

import (
	"fmt"
	"strconv"
	"strings"
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

// Parse parses a single source file (program or unit), returning the
// File instance or a *ParseError on error.
func Parse(src []byte) (file File, err error) {
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
	return p.file(), nil
}

// Parser state
type parser struct {
	// Lexer instance and current token values
	lexer *Lexer
	pos   Position // position of last token (tok)
	tok   Token    // last lexed token
	val   string   // string value of last token (or "")
}

func (p *parser) file() File {
	switch p.tok {
	case PROGRAM:
		return p.program()
	case UNIT:
		return p.unit()
	default:
		panic(p.error("expected program or unit"))
	}
}

func (p *parser) program() *Program {
	program := &Program{}

	p.expect(PROGRAM)
	program.Name = p.val
	p.expect(IDENT)
	p.expect(SEMICOLON)

	program.Uses = p.optionalUses()

	program.Decls = p.declParts(true, CONST, FUNCTION, LABEL, PROCEDURE, TYPE, VAR)

	program.Stmt = p.compoundStmt()
	p.expect(DOT)
	p.expect(EOF)

	return program
}

func (p *parser) unit() *Unit {
	unit := &Unit{}

	p.expect(UNIT)
	unit.Name = p.val
	p.expect(IDENT)
	p.expect(SEMICOLON)

	p.expect(INTERFACE)
	unit.InterfaceUses = p.optionalUses()
	unit.Interface = p.declParts(false, CONST, FUNCTION, PROCEDURE, TYPE, VAR)

	p.expect(IMPLEMENTATION)
	unit.ImplementationUses = p.optionalUses()
	unit.Implementation = p.declParts(true, CONST, FUNCTION, LABEL, PROCEDURE, TYPE, VAR)

	unit.Init = p.compoundStmt()
	p.expect(DOT)
	p.expect(EOF)

	return unit
}

func (p *parser) optionalUses() []string {
	var usesList []string
	if p.tok == USES {
		p.next()
		usesList = p.identList()
		p.expect(SEMICOLON)
	}
	return usesList
}

func (p *parser) declParts(allowBodies bool, tokens ...Token) []DeclPart {
	var decls []DeclPart
	for p.matches(tokens...) {
		decls = append(decls, p.declPart(allowBodies))
	}
	return decls
}

func (p *parser) identList() []string {
	idents := []string{p.val}
	p.expect(IDENT)
	for p.tok == COMMA {
		p.next()
		idents = append(idents, p.val)
		p.expect(IDENT)
	}
	return idents
}

func (p *parser) declPart(allowBodies bool) DeclPart {
	switch p.tok {
	case LABEL:
		p.next()
		names := p.identList()
		p.expect(SEMICOLON)
		return &LabelDecls{names}
	case CONST:
		p.next()
		decls := []*ConstDecl{}
		for p.tok == IDENT {
			name := p.val
			p.expect(IDENT)
			p.expect(EQUALS)
			value := p.constant()
			p.expect(SEMICOLON)
			decls = append(decls, &ConstDecl{name, value})
		}
		if len(decls) == 0 {
			panic(p.error("expected const declaration"))
		}
		return &ConstDecls{decls}
	case TYPE:
		p.next()
		defs := []*TypeDef{}
		for p.tok == IDENT {
			name := p.val
			p.expect(IDENT)
			p.expect(EQUALS)
			typ := p.typeIdent() // TODO: should be (type | functionType | procedureType) -- much bigger grammar
			p.expect(SEMICOLON)
			defs = append(defs, &TypeDef{name, typ})
		}
		if len(defs) == 0 {
			panic(p.error("expected type definition"))
		}
		return &TypeDefs{defs}
	case VAR:
		p.next()
		decls := []*VarDecl{}
		for p.tok == IDENT {
			names := p.identList()
			p.expect(COLON)
			typ := p.typeIdent() // TODO: should be 'type' -- much bigger grammar
			p.expect(SEMICOLON)
			decls = append(decls, &VarDecl{names, typ})
		}
		if len(decls) == 0 {
			panic(p.error("expected var declaration"))
		}
		return &VarDecls{decls}
	case PROCEDURE:
		p.next()
		name := p.val
		p.expect(IDENT)
		params := p.optionalParamList()
		p.expect(SEMICOLON)

		var decls []DeclPart
		var stmt *CompoundStmt
		if allowBodies {
			decls = p.declParts(allowBodies, CONST, FUNCTION, LABEL, PROCEDURE, TYPE, VAR)
			stmt = p.compoundStmt()
			p.expect(SEMICOLON)
		}

		return &ProcDecl{name, params, decls, stmt}
	case FUNCTION:
		p.next()
		name := p.val
		p.expect(IDENT)
		params := p.optionalParamList()
		p.expect(COLON)
		result := p.typeIdent()
		p.expect(SEMICOLON)

		var decls []DeclPart
		var stmt *CompoundStmt
		if allowBodies {
			decls = p.declParts(allowBodies, CONST, FUNCTION, LABEL, PROCEDURE, TYPE, VAR)
			stmt = p.compoundStmt()
			p.expect(SEMICOLON)
		}

		return &FuncDecl{name, params, result, decls, stmt}
	default:
		panic(p.error("expected declaration instead of %s", p.tok))
	}
}

func (p *parser) optionalParamList() []*ParamGroup {
	var groups []*ParamGroup
	if p.tok == LPAREN {
		p.next()
		groups = append(groups, p.paramGroup())
		for p.tok == SEMICOLON {
			p.next()
			groups = append(groups, p.paramGroup())
		}
		p.expect(RPAREN)
	}
	return groups
}

func (p *parser) paramGroup() *ParamGroup {
	prefix := ILLEGAL
	if p.matches(VAR, FUNCTION, PROCEDURE) {
		prefix = p.tok
		p.next()
	}
	names := p.identList()
	p.expect(COLON)
	typ := p.typeIdent()
	return &ParamGroup{prefix, names, typ}
}

func (p *parser) typeIdent() string {
	if p.matches(CHAR, BOOLEAN, INTEGER, REAL, STRING) {
		tok := p.tok
		p.next()
		return strings.ToLower(tok.String())
	}
	ident := p.val
	p.expect(IDENT)
	return ident
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
	consts := []Expr{p.constant()}
	for p.tok == COMMA {
		p.next()
		consts = append(consts, p.constant())
	}
	p.expect(COLON)
	return &CaseElement{consts, p.stmt()}
}

func (p *parser) constant() *ConstExpr {
	// TODO: should be:
	// : unsignedNumber
	// | sign unsignedNumber
	// | identifier
	// | sign identifier
	// | string
	// | constantChr
	// ;
	expr := p.factor()
	if constExpr, ok := expr.(*ConstExpr); ok {
		return constExpr
	}
	panic(p.error("expected constant"))
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
			return &ConstExpr{f, false}
		}
		return &ConstExpr{i, false}
	case HEX:
		val := p.val
		p.next()
		i, err := strconv.ParseInt(val, 16, 64)
		if err != nil {
			panic(p.error("invalid hex number: %s", err))
		}
		return &ConstExpr{int(i), true}
	case STR:
		s := p.val
		p.next()
		return &ConstExpr{s, false}
	case NOT:
		p.next()
		return &UnaryExpr{NOT, p.factor()}
	case TRUE:
		p.next()
		return &ConstExpr{true, false}
	case FALSE:
		p.next()
		return &ConstExpr{false, false}
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
