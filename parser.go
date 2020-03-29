// Turbo Pascal recursive descent parser

/*
TODO:
- allow: Char(labelPtr^) := #39;
*/

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
			var typ TypeSpec
			if p.tok == COLON {
				p.next()
				typ = p.typeSpec()
			}
			p.expect(EQUALS)
			value := p.constDeclValue()
			p.expect(SEMICOLON)
			decls = append(decls, &ConstDecl{name, typ, value})
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
			spec := p.typeSpecWithFuncProc()
			p.expect(SEMICOLON)
			defs = append(defs, &TypeDef{name, spec})
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
			typ := p.typeSpec()
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

		if p.tok == INTERRUPT {
			p.next()
			p.expect(SEMICOLON)
		}

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

func (p *parser) typeIdent() *TypeIdent {
	if p.matches(CHAR, BOOLEAN, INTEGER, REAL, STRING) {
		tok := p.tok
		p.next()
		return &TypeIdent{"", tok}
	}
	name := p.val
	p.expect(IDENT)
	return &TypeIdent{name, ILLEGAL}
}

// typeSpec: type | functionType | procedureType
func (p *parser) typeSpecWithFuncProc() TypeSpec {
	switch p.tok {
	case PROCEDURE:
		p.next()
		params := p.optionalParamList()
		return &ProcSpec{params}
	case FUNCTION:
		p.next()
		params := p.optionalParamList()
		p.expect(COLON)
		result := p.typeIdent()
		return &FuncSpec{params, result}
	default:
		return p.typeSpec()
	}
}

func (p *parser) typeSpec() TypeSpec {
	switch p.tok {
	case LPAREN:
		p.next()
		names := p.identList()
		p.expect(RPAREN)
		return &ScalarSpec{names}
	case STRING:
		p.next()
		if p.tok != LBRACKET {
			return &IdentSpec{&TypeIdent{"", STRING}}
		}
		p.expect(LBRACKET)
		size, err := strconv.Atoi(p.val)
		if err != nil {
			panic(p.error("expected integer"))
		}
		p.expect(NUM)
		p.expect(RBRACKET)
		return &StringSpec{size}
	case POINTER:
		p.next()
		typ := p.typeIdent()
		return &PointerSpec{typ}
	case ARRAY:
		p.next()
		p.expect(LBRACKET)
		min := p.expr() // much looser grammar than needed here
		p.expect(DOT_DOT)
		max := p.expr()
		p.expect(RBRACKET)
		p.expect(OF)
		ofType := p.typeSpec()
		return &ArraySpec{min, max, ofType}
	case RECORD:
		p.next()
		sections := []*RecordSection{p.recordSection()}
		for p.tok != END && p.tok != EOF {
			sections = append(sections, p.recordSection())
		}
		p.expect(END)
		return &RecordSpec{sections}
	case FILE:
		p.next()
		var ofType TypeSpec
		if p.tok == OF {
			p.next()
			ofType = p.typeSpec()
		}
		return &FileSpec{ofType}
	default:
		ident := p.typeIdent()
		return &IdentSpec{ident}
	}
}

func (p *parser) recordSection() *RecordSection {
	names := p.identList()
	p.expect(COLON)
	typ := p.typeSpec()
	p.expect(SEMICOLON)
	return &RecordSection{names, typ}
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
	case IDENT, AT, CHAR, BOOLEAN, INTEGER, REAL, STRING:
		var convType Token
		if p.matches(CHAR, BOOLEAN, INTEGER, REAL, STRING) {
			convType = p.tok
			p.next()
			p.expect(LPAREN)
		}
		varExpr := p.varExpr()
		if convType != ILLEGAL {
			p.expect(RPAREN)
		}

		switch p.tok {
		case ASSIGN:
			p.next()
			value := p.expr()
			return &AssignStmt{convType, varExpr, value}
		case COLON:
			if !varExpr.IsNameOnly() || convType != ILLEGAL {
				panic(p.error("label must be a simple identifier"))
			}
			if !allowLabel {
				panic(p.error("unexpected label"))
			}
			p.next()
			stmt := p.labelledStmt(false)
			return &LabelledStmt{varExpr.Name, stmt}
		case LPAREN:
			if convType != ILLEGAL {
				panic(p.error("can't have type conversion in procedure call"))
			}
			p.next()
			var args []Expr
			if strings.ToLower(varExpr.Name) == "str" {
				// Special case: Str(expr:width, str);
				args = append(args, p.expr())
				if p.tok == COLON {
					p.next()
					args = append(args, p.constant())
				}
				p.expect(COMMA)
				args = append(args, p.expr())
			} else {
				args = p.argList()
			}
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
		p.expect(UNTIL)
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
	consts := []Expr{p.constantOrRange()}
	for p.tok == COMMA {
		p.next()
		consts = append(consts, p.constantOrRange())
	}
	p.expect(COLON)
	return &CaseElement{consts, p.stmt()}
}

func (p *parser) constantOrRange() Expr {
	expr := p.constant()
	if p.tok == DOT_DOT {
		p.next()
		return &RangeExpr{expr, p.constant()}
	}
	return expr
}

func (p *parser) constant() Expr {
	return p.signedFactor()
}

func (p *parser) constDeclValue() Expr {
	switch p.tok {
	case LPAREN:
		p.next()
		first := p.constant()
		if p.tok == COLON { // record constant
			varExpr, ok := first.(*VarExpr)
			if !ok || !varExpr.IsNameOnly() {
				panic(p.error("expected record field: 'name: value'"))
			}
			p.expect(COLON)
			value := p.expr()
			fields := []*ConstField{&ConstField{varExpr.Name, value}}
			for p.tok == SEMICOLON {
				p.next()
				name := p.val
				p.expect(IDENT)
				p.expect(COLON)
				value = p.expr()
				fields = append(fields, &ConstField{name, value})
			}
			p.expect(RPAREN)
			return &ConstRecordExpr{fields}
		} else { // array constant
			consts := []Expr{first}
			for p.tok == COMMA {
				p.next()
				consts = append(consts, p.constant())
			}
			p.expect(RPAREN)
			return &ConstArrayExpr{consts}
		}
	default:
		return p.constant()
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
	return p.binaryExpr(p.term, p.simpleExpr, PLUS, MINUS, OR, XOR)
}

// term: signedFactor (multiplicativeOp term)?
func (p *parser) term() Expr {
	return p.binaryExpr(p.signedFactor, p.term, STAR, SLASH, DIV, MOD, AND, SHL, SHR)
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
func (p *parser) factor() Expr {
	switch p.tok {
	case LPAREN:
		p.next()
		expr := p.expr()
		p.expect(RPAREN)
		return expr
	case LBRACKET:
		p.next()
		consts := []Expr{p.constantOrRange()}
		for p.tok == COMMA {
			p.next()
			consts = append(consts, p.constantOrRange())
		}
		p.expect(RBRACKET)
		return &SetExpr{consts}
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
	case NIL:
		p.next()
		return &ConstExpr{nil, false}
	case CHAR, BOOLEAN, INTEGER, REAL, STRING:
		typ := p.tok
		p.next()
		p.expect(LPAREN)
		expr := p.expr()
		p.expect(RPAREN)
		return &TypeConvExpr{typ, expr}
	case IDENT, AT:
		varExpr := p.varExpr()
		switch p.tok {
		case LPAREN:
			p.next()
			args := p.argList()
			p.expect(RPAREN)
			var expr Expr = &FuncExpr{varExpr.Name, args}
			if p.tok == POINTER {
				p.next()
				expr = &PointerExpr{expr}
			}
			return expr
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
