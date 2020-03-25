// Turbo Pascal abstract syntax tree (AST) type

package main

import (
	"fmt"
	"strings"
)

type Program struct {
	Name  string
	Uses  []string
	Decls []DeclPart
	Stmt  *CompoundStmt
}

func (p *Program) String() string {
	usesStr := ""
	if p.Uses != nil {
		usesStr = "\nuses " + strings.Join(p.Uses, ", ") + ";"
	}

	declStrs := make([]string, len(p.Decls))
	for i, decl := range p.Decls {
		declStrs[i] = decl.String() + "\n\n"
	}

	return fmt.Sprintf("program %s;%s\n\n%s%s.\n",
		p.Name, usesStr, strings.Join(declStrs, ""), p.Stmt)
}

type DeclPart interface {
	declPart()
	String() string
}

func (p *ConstDecls) declPart() {}
func (p *FuncDecl) declPart()   {}
func (p *LabelDecls) declPart() {}
func (p *ProcDecl) declPart()   {}
func (p *TypeDefs) declPart()   {}
func (p *VarDecls) declPart()   {}

type ConstDecls struct {
	Decls []*ConstDecl
}

func (d *ConstDecls) String() string {
	strs := make([]string, len(d.Decls))
	for i, decl := range d.Decls {
		strs[i] = "    " + decl.String() + ";"
	}
	return "const\n" + strings.Join(strs, "\n")
}

type ConstDecl struct {
	Name  string
	Value *ConstExpr
}

func (d *ConstDecl) String() string {
	return fmt.Sprintf("%s = %s", d.Name, d.Value)
}

type FuncDecl struct {
	Name   string
	Params []*ParamGroup
	Result string
	Decls  []DeclPart
	Stmt   *CompoundStmt
}

func formatDecls(decls []DeclPart) string {
	declsStr := ""
	if decls != nil {
		strs := make([]string, len(decls))
		for i, decl := range decls {
			strs[i] = decl.String() + "\n"
		}
		declsStr = strings.Join(strs, "")
	}
	return declsStr
}

func (d *FuncDecl) String() string {
	return fmt.Sprintf("function %s%s: %s;\n%s%s;",
		d.Name, formatParams(d.Params), d.Result, formatDecls(d.Decls), d.Stmt)
}

type ParamGroup struct {
	Prefix Token
	Names  []string
	Type   string
}

func (g *ParamGroup) String() string {
	prefix := ""
	if g.Prefix != ILLEGAL {
		prefix = strings.ToLower(g.Prefix.String()) + " "
	}
	return fmt.Sprintf("%s%s: %s", prefix, strings.Join(g.Names, ", "), g.Type)
}

type LabelDecls struct {
	Labels []string
}

func (d *LabelDecls) String() string {
	return "label " + strings.Join(d.Labels, ", ") + ";"
}

type ProcDecl struct {
	Name   string
	Params []*ParamGroup
	Decls  []DeclPart
	Stmt   *CompoundStmt
}

func formatParams(params []*ParamGroup) string {
	str := ""
	if params != nil {
		strs := make([]string, len(params))
		for i, group := range params {
			strs[i] = group.String()
		}
		str = "(" + strings.Join(strs, "; ") + ")"
	}
	return str
}

func (d *ProcDecl) String() string {
	return fmt.Sprintf("procedure %s%s;\n%s%s;",
		d.Name, formatParams(d.Params), formatDecls(d.Decls), d.Stmt)
}

type TypeDefs struct {
	Defs []*TypeDef
}

func (d *TypeDefs) String() string {
	strs := make([]string, len(d.Defs))
	for i, def := range d.Defs {
		strs[i] = "    " + def.String() + ";"
	}
	return "type\n" + strings.Join(strs, "\n")
}

type TypeDef struct {
	Name string
	Type string // TODO: Type
}

func (d *TypeDef) String() string {
	return fmt.Sprintf("%s = %s", d.Name, d.Type)
}

type VarDecls struct {
	Decls []*VarDecl
}

func (d *VarDecls) String() string {
	strs := make([]string, len(d.Decls))
	for i, decl := range d.Decls {
		strs[i] = "    " + decl.String() + ";"
	}
	return "var\n" + strings.Join(strs, "\n")
}

type VarDecl struct {
	Names []string
	Type  string // TODO: Type
}

func (d *VarDecl) String() string {
	return fmt.Sprintf("%s: %s", strings.Join(d.Names, ", "), d.Type)
}

// Statements

type Stmt interface {
	stmt()
	String() string
}

func (s *AssignStmt) stmt()   {}
func (s *CaseStmt) stmt()     {}
func (s *CompoundStmt) stmt() {}
func (s *EmptyStmt) stmt()    {}
func (s *ForStmt) stmt()      {}
func (s *GotoStmt) stmt()     {}
func (s *IfStmt) stmt()       {}
func (s *LabelledStmt) stmt() {}
func (s *ProcStmt) stmt()     {}
func (s *RepeatStmt) stmt()   {}
func (s *WhileStmt) stmt()    {}
func (s *WithStmt) stmt()     {}

type AssignStmt struct {
	Var   *VarExpr
	Value Expr
}

func (s *AssignStmt) String() string {
	return fmt.Sprintf("%s := %s", s.Var, s.Value)
}

type CaseStmt struct {
	Selector Expr
	Cases    []*CaseElement
	Else     []Stmt
}

func (s *CaseStmt) String() string {
	caseStrs := make([]string, len(s.Cases))
	for i, c := range s.Cases {
		caseStrs[i] = "    " + c.String() + ";\n" // TODO: indentation
	}
	elseStr := ""
	if s.Else != nil {
		elseStr = "else\n" + indentStmts(s.Else)
	}
	return fmt.Sprintf("case %s of\n%s%send",
		s.Selector, strings.Join(caseStrs, ""), elseStr)
}

type CaseElement struct {
	Consts []Expr // TODO: should be []Constant
	Stmt   Stmt
}

func (e *CaseElement) String() string {
	constStrs := make([]string, len(e.Consts))
	for i, c := range e.Consts {
		constStrs[i] = c.String()
	}
	return fmt.Sprintf("%s: %s", strings.Join(constStrs, ", "), e.Stmt)
}

type CompoundStmt struct {
	Stmts []Stmt
}

func indentStmts(stmts []Stmt) string {
	lines := []string{}
	for _, stmt := range stmts {
		subLines := strings.Split(stmt.String()+";", "\n")
		for _, sl := range subLines {
			lines = append(lines, "    "+sl+"\n")
		}
	}
	// TODO: better way to do this with Join?
	if len(lines) > 0 && lines[len(lines)-1] == "    ;\n" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "")
}

func (s *CompoundStmt) String() string {
	return "begin\n" + indentStmts(s.Stmts) + "end"
}

type EmptyStmt struct{}

func (s *EmptyStmt) String() string {
	return ""
}

type ForStmt struct {
	Var     string // TODO
	Initial Expr
	Down    bool
	Final   Expr
	Stmt    Stmt
}

func (s *ForStmt) String() string {
	toStr := "to"
	if s.Down {
		toStr = "downto"
	}
	return fmt.Sprintf("for %s := %s %s %s do %s", s.Var, s.Initial, toStr, s.Final, s.Stmt)
}

type GotoStmt struct {
	Label string
}

func (s *GotoStmt) String() string {
	return fmt.Sprintf("goto %s", s.Label)
}

type IfStmt struct {
	Cond Expr
	Then Stmt
	Else Stmt
}

func (s *IfStmt) String() string {
	str := fmt.Sprintf("if %s then %s", s.Cond, s.Then)
	if s.Else != nil {
		str += fmt.Sprintf(" else %s", s.Else)
	}
	return str
}

type LabelledStmt struct {
	Label string
	Stmt  Stmt
}

func (s *LabelledStmt) String() string {
	return fmt.Sprintf("%s:\n%s", s.Label, s.Stmt)
}

type ProcStmt struct {
	Proc string
	Args []Expr
}

func formatArgList(args []Expr) string {
	strs := make([]string, len(args))
	for i, arg := range args {
		strs[i] = arg.String()
	}
	return "(" + strings.Join(strs, ", ") + ")"
}

func (s *ProcStmt) String() string {
	str := s.Proc
	if s.Args != nil {
		str += formatArgList(s.Args)
	}
	return str
}

type RepeatStmt struct {
	Stmts []Stmt
	Cond  Expr
}

func (s *RepeatStmt) String() string {
	return fmt.Sprintf("begin\n%suntil %s", indentStmts(s.Stmts), s.Cond)
}

type WhileStmt struct {
	Cond Expr
	Stmt Stmt
}

func (s *WhileStmt) String() string {
	return fmt.Sprintf("while %s do %s", s.Cond, s.Stmt)
}

type WithStmt struct {
	Vars []*VarExpr
	Stmt Stmt
}

func (s *WithStmt) String() string {
	strs := make([]string, len(s.Vars))
	for i, v := range s.Vars {
		strs[i] = v.String()
	}
	return fmt.Sprintf("with %s do %s", strings.Join(strs, ", "), s.Stmt)
}

// Expressions

type Expr interface {
	expr()
	String() string
}

func (e *BinaryExpr) expr() {}
func (e *ConstExpr) expr()  {}
func (e *FuncExpr) expr()   {}
func (e *UnaryExpr) expr()  {}
func (e *VarExpr) expr()    {}

type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left, strings.ToLower(e.Op.String()), e.Right) // TODO: handle precedence
}

type ConstExpr struct {
	Value interface{}
}

func (e *ConstExpr) String() string {
	switch value := e.Value.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "'#39'"))
	default:
		return fmt.Sprintf("%v", value)
	}
}

type FuncExpr struct {
	Func string
	Args []Expr
}

func (e *FuncExpr) String() string {
	return e.Func + formatArgList(e.Args)
}

type UnaryExpr struct {
	Op   Token
	Expr Expr
}

func (e *UnaryExpr) String() string {
	opStr := e.Op.String()
	if opStr[0] >= 'A' && opStr[0] <= 'Z' {
		return fmt.Sprintf("%s %s", strings.ToLower(opStr), e.Expr)
	}
	return fmt.Sprintf("%s%s", e.Op, e.Expr)
}

type VarExpr struct {
	HasAt    bool
	Name     string
	Suffixes []VarSuffix
}

func (v *VarExpr) IsNameOnly() bool {
	return !v.HasAt && v.Suffixes == nil
}

func (v *VarExpr) String() string {
	parts := []string{}
	if v.HasAt {
		parts = append(parts, "@")
	}
	parts = append(parts, v.Name)
	for _, part := range v.Suffixes {
		parts = append(parts, part.String())
	}
	return strings.Join(parts, "")
}

type VarSuffix interface {
	varSuffix()
	String() string
}

func (s *IndexSuffix) varSuffix()   {}
func (s *DotSuffix) varSuffix()     {}
func (s *PointerSuffix) varSuffix() {}

type IndexSuffix struct {
	Indexes []Expr
}

func (s *IndexSuffix) String() string {
	strs := make([]string, len(s.Indexes))
	for i, index := range s.Indexes {
		strs[i] = index.String()
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

type DotSuffix struct {
	Field string
}

func (s *DotSuffix) String() string {
	return "." + s.Field
}

type PointerSuffix struct{}

func (s *PointerSuffix) String() string {
	return "^"
}
