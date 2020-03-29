// Turbo Pascal abstract syntax tree (AST) type

package main

import (
	"fmt"
	"strings"
)

type File interface {
	file()
	String() string
}

func (f *Program) file() {}
func (f *Unit) file()    {}

type Program struct {
	Name  string
	Uses  []string
	Decls []DeclPart
	Stmt  *CompoundStmt
}

func (p *Program) String() string {
	return fmt.Sprintf(`program %s;%s

%s%s.
`,
		p.Name, formatUses(p.Uses), formatDecls(p.Decls), p.Stmt)
}

func formatUses(uses []string) string {
	str := ""
	if uses != nil {
		str = "\nuses " + strings.Join(uses, ", ") + ";"
	}
	return str
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

type Unit struct {
	Name               string
	InterfaceUses      []string
	Interface          []DeclPart
	ImplementationUses []string
	Implementation     []DeclPart
	Init               *CompoundStmt
}

func (u *Unit) String() string {
	return fmt.Sprintf(`unit %s;

interface%s
	%s

implementation%s
    %s%s.
`,
		u.Name, formatUses(u.InterfaceUses), formatDecls(u.Interface),
		formatUses(u.ImplementationUses), formatDecls(u.Implementation), u.Init)
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
	Type  TypeSpec
	Value Expr
}

func (d *ConstDecl) String() string {
	typeStr := ""
	if d.Type != nil {
		typeStr = fmt.Sprintf(": %s", d.Type)
	}
	return fmt.Sprintf("%s%s = %s", d.Name, typeStr, d.Value)
}

type FuncDecl struct {
	Name   string
	Params []*ParamGroup
	Result *TypeIdent
	Decls  []DeclPart
	Stmt   *CompoundStmt
}

func (d *FuncDecl) String() string {
	stmtStr := ""
	if d.Stmt != nil {
		stmtStr = d.Stmt.String() + ";\n"
	}
	return fmt.Sprintf("function %s%s: %s;\n%s%s",
		d.Name, formatParams(d.Params), d.Result, formatDecls(d.Decls), stmtStr)
}

type TypeIdent struct {
	Name    string
	Builtin Token
}

func (t *TypeIdent) String() string {
	if t.Name == "" {
		return strings.ToLower(t.Builtin.String())
	}
	return t.Name
}

type ParamGroup struct {
	Prefix Token // ILLEGAL (meaning no prefix), VAR, FUNCTION, PROCEDURE
	Names  []string
	Type   *TypeIdent
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
	stmtStr := ""
	if d.Stmt != nil {
		stmtStr = d.Stmt.String() + ";\n"
	}
	return fmt.Sprintf("procedure %s%s;\n%s%s",
		d.Name, formatParams(d.Params), formatDecls(d.Decls), stmtStr)
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
	Type TypeSpec
}

func (d *TypeDef) String() string {
	return fmt.Sprintf("%s = %s", d.Name, d.Type)
}

type TypeSpec interface {
	typeSpec()
	String() string
}

func (s *FuncSpec) typeSpec()    {}
func (s *ProcSpec) typeSpec()    {}
func (s *ScalarSpec) typeSpec()  {}
func (s *IdentSpec) typeSpec()   {}
func (s *StringSpec) typeSpec()  {}
func (s *ArraySpec) typeSpec()   {}
func (s *RecordSpec) typeSpec()  {}
func (s *FileSpec) typeSpec()    {}
func (s *PointerSpec) typeSpec() {}

type FuncSpec struct {
	Params []*ParamGroup
	Result *TypeIdent
}

func (s *FuncSpec) String() string {
	return fmt.Sprintf("%s: %s", formatParams(s.Params), s.Result)
}

type ProcSpec struct {
	Params []*ParamGroup
}

func (s *ProcSpec) String() string {
	return formatParams(s.Params)
}

type ScalarSpec struct {
	Names []string
}

func (s *ScalarSpec) String() string {
	return "(" + strings.Join(s.Names, ", ") + ")"
}

type IdentSpec struct {
	*TypeIdent
}

type StringSpec struct {
	Size int
}

func (s *StringSpec) String() string {
	return fmt.Sprintf("string[%d]", s.Size)
}

type ArraySpec struct {
	Min Expr
	Max Expr
	Of  TypeSpec
}

func (s *ArraySpec) String() string {
	return fmt.Sprintf("array[%s .. %s] of %s", s.Min, s.Max, s.Of)
}

type RecordSpec struct {
	Sections []*RecordSection
}

func (s *RecordSpec) String() string {
	strs := []string{"record\n"}
	for _, section := range s.Sections {
		strs = append(strs, "    "+section.String()+";\n")
	}
	strs = append(strs, "end")
	return strings.Join(strs, "")
}

type RecordSection struct {
	Names []string
	Type  TypeSpec
}

func (s *RecordSection) String() string {
	return fmt.Sprintf("%s: %s", strings.Join(s.Names, ", "), s.Type)
}

type FileSpec struct {
	Of TypeSpec
}

func (s *FileSpec) String() string {
	if s.Of != nil {
		return fmt.Sprintf("file of %s", s.Of)
	}
	return "file"
}

type PointerSpec struct {
	Type *TypeIdent
}

func (s *PointerSpec) String() string {
	return "^" + s.Type.String()
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
	Type  TypeSpec
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
	TypeConv Token
	Var      *VarExpr
	Value    Expr
}

func (s *AssignStmt) String() string {
	if s.TypeConv != ILLEGAL {
		typeStr := s.TypeConv.String()
		typeStr = string(typeStr[0]) + strings.ToLower(typeStr[1:])
		return fmt.Sprintf("%s(%s) := %s", typeStr, s.Var, s.Value)
	}
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
		caseStrs[i] = "    " + c.String() + ";\n"
	}
	elseStr := ""
	if s.Else != nil {
		elseStr = "else\n" + indentStmts(s.Else)
	}
	return fmt.Sprintf("case %s of\n%s%send",
		s.Selector, strings.Join(caseStrs, ""), elseStr)
}

type CaseElement struct {
	Consts []Expr
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
	return fmt.Sprintf("repeat\n%suntil %s", indentStmts(s.Stmts), s.Cond)
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

func (e *BinaryExpr) expr()      {}
func (e *ConstExpr) expr()       {}
func (e *ConstArrayExpr) expr()  {}
func (e *ConstRecordExpr) expr() {}
func (e *FuncExpr) expr()        {}
func (e *PointerExpr) expr()     {}
func (e *RangeExpr) expr()       {}
func (e *SetExpr) expr()         {}
func (e *TypeConvExpr) expr()    {}
func (e *UnaryExpr) expr()       {}
func (e *VarExpr) expr()         {}

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
	IsHex bool
}

func (e *ConstExpr) String() string {
	switch value := e.Value.(type) {
	case string:
		return escapeString(value)
	case nil:
		return "nil"
	default:
		if e.IsHex {
			return fmt.Sprintf("$%02X", value)
		}
		return fmt.Sprintf("%v", value)
	}
}

func escapeString(s string) string {
	if s == "" {
		return "''"
	}
	out := make([]byte, 0, len(s))
	out = append(out, '\'')
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 32 || c == 39 || c > 126 {
			if out[len(out)-1] == '\'' {
				out = out[:len(out)-1]
				out = append(out, []byte(fmt.Sprintf("#%d'", c))...)
			} else {
				out = append(out, []byte(fmt.Sprintf("'#%d'", c))...)
			}
		} else {
			out = append(out, c)
		}
	}
	out = append(out, '\'')
	t := string(out)
	t = strings.TrimPrefix(t, "''")
	t = strings.TrimSuffix(t, "''")
	return t
}

type ConstArrayExpr struct {
	Values []Expr
}

func (e *ConstArrayExpr) String() string {
	strs := make([]string, len(e.Values))
	for i, v := range e.Values {
		strs[i] = v.String()
	}
	return "(" + strings.Join(strs, ", ") + ")"
}

type ConstRecordExpr struct {
	Fields []*ConstField
}

func (e *ConstRecordExpr) String() string {
	strs := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		strs[i] = f.String()
	}
	return "(" + strings.Join(strs, "; ") + ")"
}

type ConstField struct {
	Name  string
	Value Expr
}

func (f *ConstField) String() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Value)
}

type FuncExpr struct {
	Func string
	Args []Expr
}

func (e *FuncExpr) String() string {
	return e.Func + formatArgList(e.Args)
}

type PointerExpr struct {
	Expr Expr
}

func (e *PointerExpr) String() string {
	return e.Expr.String() + "^"
}

type RangeExpr struct {
	Min Expr
	Max Expr
}

func (e *RangeExpr) String() string {
	return fmt.Sprintf("%s .. %s", e.Min, e.Max)
}

type SetExpr struct {
	Values []Expr
}

func (e *SetExpr) String() string {
	strs := make([]string, len(e.Values))
	for i, v := range e.Values {
		strs[i] = v.String()
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

type TypeConvExpr struct {
	Type Token
	Expr Expr
}

func (e *TypeConvExpr) String() string {
	typeStr := e.Type.String()
	typeStr = string(typeStr[0]) + strings.ToLower(typeStr[1:])
	return fmt.Sprintf("%s(%s)", typeStr, e.Expr)
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
