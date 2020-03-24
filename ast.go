// Turbo Pascal abstract syntax tree (AST) type

package main

import (
	"fmt"
	"strings"
)

type Program struct {
	Name  string
	Decls []DeclPart
	Stmt  *CompoundStmt
}

func (p *Program) String() string {
	s := fmt.Sprintf(`program %s;

%s.
`, p.Name, p.Stmt)
	return s
}

type DeclPart interface {
	declPart()
}

func (p *ConstDecls) declPart() {}
func (p *Function) declPart()   {}
func (p *LabelDecls) declPart() {}
func (p *Procedure) declPart()  {}
func (p *TypeDefs) declPart()   {}
func (p *VarDecls) declPart()   {}

type ConstDecls struct {
	Decls []ConstDecl
}

type ConstDecl struct {
	Name  string
	Value Expr
}

type Function struct {
	Name       string
	Params     []ParamGroup
	ReturnType string // TODO: Type
	//TODO:	Block Block
}

type ParamGroup struct {
	IsVar bool
	Names []string
	Type  string
}

type LabelDecls struct {
	Labels []string
}

type Procedure struct {
	Name   string
	Params []ParamGroup
	//TODO:	Block Block
}

type TypeDefs struct {
	Defs []TypeDef
}

type TypeDef struct {
	Name string
	Type string // TODO: Type
}

type VarDecls struct {
	Decls []VarDecl
}

type VarDecl struct {
	Names []string
	Type  string // TODO: Type
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

type CaseStmt struct { // TODO
}

func (s *CaseStmt) String() string {
	return "case TODO"
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
	// TODO: handle precedence
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Op, e.Right)
}

type ConstExpr struct {
	Value interface{}
}

func (e *ConstExpr) String() string {
	switch e.Value.(type) {
	case string:
		return fmt.Sprintf("'%s'", e.Value) // TODO: proper quoting
	default:
		return fmt.Sprintf("%v", e.Value)
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
	c := e.Op.String()[0]
	if c >= 'A' && c <= 'Z' {
		fmt.Sprintf("%s %s", e.Op, e.Expr)
	}
	return fmt.Sprintf("%s%s", e.Op, e.Expr)
}

type VarExpr struct {
	Name string
}

func (v *VarExpr) String() string {
	return v.Name
}
