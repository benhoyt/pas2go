// (Try to!) convert a Turbo Pascal AST into Go code

/*
ISSUES:
- distinguishing string constants vs char, eg: pArg[1] == "/"
- uses operator precedence rather than ParenExpr
- "exit" -> break or return (should EXIT be a keyword in lexer?)
- handle array[Max-Max+1] thing better
*/

package main

import (
	"fmt"
	"io"
	"strings"
)

func Convert(file File, w io.Writer) {
	c := &converter{w}
	switch file := file.(type) {
	case *Program:
		c.program(file)
	case *Unit:
		c.unit(file)
	default:
		panic(fmt.Sprintf("unhandled File type: %T", file))
	}
}

type converter struct {
	w io.Writer
}

func (c *converter) program(program *Program) {
	fmt.Fprint(c.w, "package main\n\n")
	if program.Uses != nil {
		fmt.Fprintf(c.w, "// uses: %s\n\n", strings.Join(program.Uses, ", "))
	}
	c.decls(program.Decls, true)
	fmt.Fprint(c.w, "func main() {\n")
	c.stmts(program.Stmt.Stmts)
	fmt.Fprint(c.w, "}\n")
}

func (c *converter) unit(unit *Unit) {
	fmt.Fprintf(c.w, "package main // unit: %s\n\n", unit.Name)
	if unit.InterfaceUses != nil {
		fmt.Fprintf(c.w, "// interface uses: %s\n\n", strings.Join(unit.InterfaceUses, ", "))
	}
	c.decls(unit.Interface, true)
	if unit.ImplementationUses != nil {
		fmt.Fprintf(c.w, "\n// implementation uses: %s\n\n", strings.Join(unit.ImplementationUses, ", "))
	}
	c.decls(unit.Implementation, true)
	fmt.Fprint(c.w, "func init() {\n")
	c.stmts(unit.Init.Stmts)
	fmt.Fprint(c.w, "}\n")
}

func (c *converter) decls(decls []DeclPart, isMain bool) {
	for _, decl := range decls {
		c.decl(decl, isMain)
	}
}

func (c *converter) decl(decl DeclPart, isMain bool) {
	switch decl := decl.(type) {
	case *ConstDecls:
		fmt.Fprint(c.w, "const (\n")
		for _, d := range decl.Decls {
			fmt.Fprintf(c.w, "%s", d.Name)
			if d.Type != nil {
				fmt.Fprint(c.w, " ")
				c.typeSpec(d.Type)
			}
			fmt.Fprint(c.w, " = ")
			c.expr(d.Value)
			fmt.Fprint(c.w, "\n")
		}
		fmt.Fprint(c.w, ")\n")
	case *FuncDecl:
		if decl.Stmt == nil {
			return
		}
		if isMain {
			fmt.Fprintf(c.w, "func %s(", decl.Name)
		} else {
			fmt.Fprintf(c.w, "%s := func(", decl.Name)
		}
		c.params(decl.Params)
		fmt.Fprintf(c.w, ") (%s ", decl.Name)
		c.typeIdent(decl.Result)
		fmt.Fprint(c.w, ") {\n")
		c.decls(decl.Decls, false)
		c.stmts(decl.Stmt.Stmts)
		fmt.Fprint(c.w, "return\n}\n\n")
	case *LabelDecls:
		// not needed
	case *ProcDecl:
		if decl.Stmt == nil {
			return
		}
		if isMain {
			fmt.Fprintf(c.w, "func %s(", decl.Name)
		} else {
			fmt.Fprintf(c.w, "%s := func(", decl.Name)
		}
		c.params(decl.Params)
		fmt.Fprint(c.w, ") {\n")
		c.decls(decl.Decls, false)
		c.stmts(decl.Stmt.Stmts)
		fmt.Fprint(c.w, "}\n\n")
	case *TypeDefs:
		fmt.Fprint(c.w, "type (\n")
		for _, d := range decl.Defs {
			fmt.Fprintf(c.w, "%s ", d.Name)
			c.typeSpec(d.Type)
			fmt.Fprint(c.w, "\n")
		}
		fmt.Fprint(c.w, ")\n")
	case *VarDecls:
		fmt.Fprint(c.w, "var (\n")
		for _, d := range decl.Decls {
			fmt.Fprintf(c.w, "%s ", strings.Join(d.Names, ", "))
			c.typeSpec(d.Type)
			fmt.Fprint(c.w, "\n")
		}
		fmt.Fprint(c.w, ")\n")
	default:
		panic(fmt.Sprintf("unhandled DeclPart type: %T", decl))
	}
}

func (c *converter) params(params []*ParamGroup) {
	for i, param := range params {
		if i > 0 {
			fmt.Fprint(c.w, ", ")
		}
		fmt.Fprint(c.w, strings.Join(param.Names, ", "), " ")
		switch param.Prefix {
		case VAR:
			fmt.Fprint(c.w, "*")
		case ILLEGAL:
			// no prefix
		default:
			panic(fmt.Sprintf("unhandled ParamGroup.Prefix: %s", param.Prefix))
		}
		c.typeIdent(param.Type)
	}
}

func (c *converter) typeIdent(typ *TypeIdent) {
	var s string
	switch typ.Builtin {
	case CHAR:
		s = "byte"
	case BOOLEAN:
		s = "bool"
	case INTEGER:
		s = "int"
	case REAL:
		s = "float64"
	case STRING:
		s = "string"
	default:
		s = typ.Name
		if strings.ToLower(s) == "pointer" {
			s = "uintptr" // TODO: hmmm
		}
	}
	fmt.Fprint(c.w, s)
}

func (c *converter) stmts(stmts []Stmt) {
	for _, stmt := range stmts {
		c.stmt(stmt)
	}
}

func (c *converter) stmtNoBraces(stmt Stmt) {
	switch stmt := stmt.(type) {
	case *CompoundStmt:
		c.stmts(stmt.Stmts)
	default:
		c.stmt(stmt)
	}
}

func (c *converter) stmt(stmt Stmt) {
	switch stmt := stmt.(type) {
	case *AssignStmt:
		// TODO: handle TypeConv?
		c.expr(stmt.Var)
		fmt.Fprint(c.w, " = ")
		c.expr(stmt.Value)
	case *CaseStmt:
		fmt.Fprint(c.w, "switch ")
		c.expr(stmt.Selector)
		fmt.Fprint(c.w, " {\n")
		for _, cas := range stmt.Cases {
			fmt.Fprint(c.w, "case ")
			c.exprs(cas.Consts) // TODO: handle RangeExpr with: case 1, 2, 3:
			fmt.Fprint(c.w, ":\n")
			c.stmtNoBraces(cas.Stmt)
		}
		if stmt.Else != nil {
			fmt.Fprint(c.w, "default:\n")
			c.stmts(stmt.Else)
		}
		fmt.Fprint(c.w, "}")
	case *CompoundStmt:
		fmt.Fprint(c.w, "{\n")
		c.stmts(stmt.Stmts)
		fmt.Fprint(c.w, "}")
	case *EmptyStmt:
		return
	case *ForStmt:
		fmt.Fprintf(c.w, "for %s := ", stmt.Var)
		c.expr(stmt.Initial)
		if stmt.Down {
			fmt.Fprintf(c.w, "; %s >= ", stmt.Var)
			c.expr(stmt.Final)
			fmt.Fprintf(c.w, "; %s-- {\n", stmt.Var)
		} else {
			fmt.Fprintf(c.w, "; %s <= ", stmt.Var)
			c.expr(stmt.Final)
			fmt.Fprintf(c.w, "; %s++ {\n", stmt.Var)
		}
		c.stmtNoBraces(stmt.Stmt)
		fmt.Fprint(c.w, "}")
	case *GotoStmt:
		fmt.Fprintf(c.w, "goto %s", stmt.Label)
	case *IfStmt:
		fmt.Fprint(c.w, "if ")
		c.expr(stmt.Cond)
		fmt.Fprint(c.w, " {\n")
		c.stmtNoBraces(stmt.Then)
		fmt.Fprint(c.w, "}")
		if stmt.Else != nil {
			innerIf, isElseIf := stmt.Else.(*IfStmt)
			if isElseIf {
				fmt.Fprint(c.w, " else ")
				c.stmtNoBraces(innerIf)
			} else {
				fmt.Fprint(c.w, " else {\n")
				c.stmtNoBraces(stmt.Else)
				fmt.Fprint(c.w, "}")
			}
		}
	case *LabelledStmt:
		fmt.Fprintf(c.w, "%s:\n", stmt.Label)
		c.stmt(stmt.Stmt)
	case *ProcStmt:
		c.expr(stmt.Proc)
		fmt.Fprint(c.w, "(")
		c.exprs(stmt.Args)
		fmt.Fprint(c.w, ")")
	case *RepeatStmt:
		fmt.Fprint(c.w, "for {\n")
		c.stmts(stmt.Stmts)
		fmt.Fprint(c.w, "if ")
		c.expr(stmt.Cond)
		fmt.Fprint(c.w, " {\nbreak\n}\n}")
	case *WhileStmt:
		fmt.Fprint(c.w, "for ")
		c.expr(stmt.Cond)
		fmt.Fprint(c.w, " {\n")
		c.stmtNoBraces(stmt.Stmt)
		fmt.Fprint(c.w, "}")
	case *WithStmt:
		// TODO: fix this; drop multi Vars support?
		fmt.Fprint(c.w, "// WITH temp = ")
		c.expr(stmt.Vars[0])
		fmt.Fprint(c.w, "\n")
		c.stmtNoBraces(stmt.Stmt)
	default:
		panic(fmt.Sprintf("unhandled Stmt: %T", stmt))
	}
	fmt.Fprint(c.w, "\n")
}

func (c *converter) exprs(exprs []Expr) {
	for i, expr := range exprs {
		if i > 0 {
			fmt.Fprint(c.w, ", ")
		}
		c.expr(expr)
	}
}

func (c *converter) expr(expr Expr) {
	switch expr := expr.(type) {
	case *BinaryExpr:
		// TODO: handle precedence instead of ParenExpr?
		if expr.Op == IN {
			c.inExpr(expr)
			return
		}
		c.expr(expr.Left)
		fmt.Fprintf(c.w, " %s ", operatorStr(expr.Op))
		c.expr(expr.Right)
	case *ConstExpr:
		switch value := expr.Value.(type) {
		case string:
			fmt.Fprintf(c.w, "%q", value)
		case float64:
			s := fmt.Sprintf("%g", value)
			if !strings.Contains(s, ".") {
				s += ".0"
			}
			fmt.Fprint(c.w, s)
		case nil:
			fmt.Fprint(c.w, "nil")
		default:
			if expr.IsHex {
				fmt.Fprintf(c.w, "0x%02X", value)
			} else {
				fmt.Fprintf(c.w, "%v", value)
			}
		}
	case *ConstArrayExpr:
		fmt.Fprint(c.w, "[...]string{") // TODO: not necessarily string
		c.exprs(expr.Values)
		fmt.Fprint(c.w, "}")
	case *ConstRecordExpr:
		// TODO: need type of const expr here
		fmt.Fprint(c.w, "{")
		for i, field := range expr.Fields {
			if i > 0 {
				fmt.Fprint(c.w, ", ")
			}
			fmt.Fprint(c.w, field.Name)
			fmt.Fprint(c.w, ": ")
			c.expr(field.Value)
		}
		fmt.Fprint(c.w, "}")
	case *FuncExpr:
		c.expr(expr.Func)
		fmt.Fprint(c.w, "(")
		c.exprs(expr.Args)
		fmt.Fprint(c.w, ")")
	case *ParenExpr:
		fmt.Fprint(c.w, "(")
		c.expr(expr.Expr)
		fmt.Fprint(c.w, ")")
	case *PointerExpr:
		fmt.Fprint(c.w, "&")
		c.expr(expr.Expr)
	// case *RangeExpr:
	// case *SetExpr:
	case *TypeConvExpr:
		c.typeIdent(&TypeIdent{"", expr.Type})
		fmt.Fprint(c.w, "(")
		c.expr(expr.Expr)
		fmt.Fprint(c.w, ")")
	case *UnaryExpr:
		fmt.Fprint(c.w, operatorStr(expr.Op))
		c.expr(expr.Expr)
	case *VarExpr:
		if expr.HasAt {
			fmt.Fprintf(c.w, "*")
		}
		fmt.Fprintf(c.w, expr.Name)
		for _, suffix := range expr.Suffixes {
			switch suffix := suffix.(type) {
			case *IndexSuffix:
				fmt.Fprint(c.w, "[")
				c.exprs(suffix.Indexes)
				fmt.Fprint(c.w, "]")
			case *DotSuffix:
				fmt.Fprint(c.w, ".", suffix.Field)
			case *PointerSuffix:
			default:
				panic(fmt.Sprintf("unhandled VarSuffix: %T", suffix))
			}
		}
	case *WidthExpr:
		c.expr(expr.Expr)
		// TODO: handle Width
	default:
		fmt.Fprintf(c.w, "%s", expr)
		//TODO		panic(fmt.Sprintf("unhandled Stmt: %T", stmt))
	}
}

func (c *converter) inExpr(expr *BinaryExpr) {
	fmt.Fprint(c.w, "(")
	values := expr.Right.(*SetExpr)
	for i, value := range values.Values {
		if i > 0 {
			fmt.Fprint(c.w, " || ")
		}
		if rangeExpr, ok := value.(*RangeExpr); ok {
			c.expr(expr.Left)
			fmt.Fprint(c.w, ">=")
			c.expr(rangeExpr.Min)
			fmt.Fprint(c.w, " && ")
			c.expr(expr.Left)
			fmt.Fprint(c.w, "<=")
			c.expr(rangeExpr.Max)
		} else {
			c.expr(expr.Left)
			fmt.Fprint(c.w, "==")
			c.expr(value)
		}
	}
	fmt.Fprint(c.w, ")")
}

func (c *converter) typeSpec(spec TypeSpec) {
	switch spec := spec.(type) {
	case *FuncSpec:
		fmt.Fprint(c.w, "func(")
		c.params(spec.Params)
		fmt.Fprint(c.w, ") ")
		c.typeIdent(spec.Result)
	case *ProcSpec:
		fmt.Fprint(c.w, "func(")
		c.params(spec.Params)
		fmt.Fprint(c.w, ")")
	case *ScalarSpec:
		// TODO: also define constants, see EDITOR.PAS: TDrawMode = (DrawingOff, DrawingOn, TextEntry);
		fmt.Fprint(c.w, "uint8")
	case *IdentSpec:
		c.typeIdent(spec.TypeIdent)
	case *StringSpec:
		// TODO: how to handle string sizes? should we use [Size]byte
		fmt.Fprint(c.w, "string")
	case *ArraySpec:
		// TODO: tidy up size if constants
		// TODO: how to deal with Min?
		fmt.Fprint(c.w, "[")
		c.expr(spec.Max)
		fmt.Fprint(c.w, "-")
		c.expr(spec.Min)
		fmt.Fprint(c.w, "+1]")
		c.typeSpec(spec.Of)
	case *RecordSpec:
		fmt.Fprint(c.w, "struct {\n")
		for _, section := range spec.Sections {
			fmt.Fprint(c.w, strings.Join(section.Names, ", "), " ")
			c.typeSpec(section.Type)
			fmt.Fprint(c.w, "\n")
		}
		fmt.Fprint(c.w, "}")
	case *FileSpec:
		// TODO: handle Of
		fmt.Fprint(c.w, "FILE")
	case *PointerSpec:
		fmt.Fprint(c.w, "*")
		c.typeIdent(spec.Type)
	default:
		fmt.Fprintf(c.w, "%s", spec)
		//TODO		panic(fmt.Sprintf("unhandled TypeSpec: %T", spec))
	}
}

func operatorStr(op Token) string {
	switch op {
	case EQUALS:
		return "=="
	case NOT_EQUALS:
		return "!="
	case OR:
		return "||"
	case XOR:
		return "^" // TODO: note, only for integers in Go
	case DIV:
		return "/"
	case MOD:
		return "%"
	case AND:
		return "&&"
	case SHL:
		return "<<"
	case SHR:
		return ">>"
	case NOT:
		return "!"
	default:
		// same as in Pascal
		return op.String()
	}
}
