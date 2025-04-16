// ast.go
package main

import "fmt"

// ExprEvaluator is an interface for expressions that can be evaluated
type ExprEvaluator interface {
	Eval(env *Environment) interface{} // Method to evaluate the expression
}

// Expr interface for all expression nodes, extended to include ExprEvaluator
type Expr interface {
	String() string
	ExprEvaluator // Include evaluation in the expression interface
}

// Literal struct for literal values (booleans, numbers, strings, nil)
type Literal struct {
	Value interface{}
	Type  string
}

// String method for Literal to print its content
func (l *Literal) String() string {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.Value)
}

// Grouping struct to represent expressions inside parentheses
type Grouping struct {
	Expression Expr
}

func (g *Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.Expression.String())
}

// Unary struct for unary operators
type Unary struct {
	Operator Token
	Right    Expr
	Line     int
}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Right.String())
}

// Binary struct to represent binary expressions (e.g., 16 * 38)
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
	Line     int
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

// Stmt interface for statements
type Stmt interface {
	Expr // Method to evaluate the statement
}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Expression Expr
}

// String method for ExpressionStatement
func (e *ExpressionStatement) String() string {
	return e.Expression.String() // Return string representation of the expression
}

type PrintStatement struct {
	Expression Expr
}

// String method for PrintStatement
func (p *PrintStatement) String() string {
	return fmt.Sprintf("(print %s)", p.Expression.String()) // Return string representation of print statement
}

// VarStmt represents a variable declaration statement
type VarStmt struct {
	Name        string
	Initializer Expr
	VarUsed     bool
	Line        int
}

func (v *VarStmt) String() string {
	return fmt.Sprintf("var %s = %v", v.Name, v.Initializer)
}

// Identifier represents a variable being used in an expression
type Identifier struct {
	Name string
	Line int
}

func (i *Identifier) String() string {
	return i.Name
}

type AssignStmt struct {
	Name  string
	Value Expr
	Line  int
}

func (a *AssignStmt) String() string {
	return fmt.Sprintf("(%s = %s)", a.Name, a.Value.String())
}

type BlockStmt struct {
	Statements []Stmt
}

func (b *BlockStmt) String() string {
	val := fmt.Sprintf("{\n")
	for _, statement := range b.Statements {
		val += fmt.Sprintf("%s\n", statement.String())
	}
	val += fmt.Sprint("}")
	return val
}

type IfStmt struct {
	Condition Stmt
	Body      []Stmt
	Else      []Stmt
}

func (b *IfStmt) String() string {
	val := fmt.Sprintf("if (%v) {\n", b.Condition)
	for _, statement := range b.Body {
		val += fmt.Sprintf("%s\n", statement.String())
	}
	val += "}"
	if b.Else != nil && len(b.Else) > 0 {
		val += " else {\n"
		for _, statement := range b.Else {
			val += fmt.Sprintf("%s\n", statement.String())
		}
		val += "}"
	}
	return val
}

type WhileStmt struct {
	Condition Stmt
	Body      []Stmt
}

func (w *WhileStmt) String() string {
	val := fmt.Sprintf("while (%v) {\n", w.Condition)
	for _, statement := range w.Body {
		val += fmt.Sprintf("%s\n", statement.String())
	}
	val += "}"
	return val
}

type ForStmt struct {
	Initializer Stmt
	Condition   Expr
	Increment   Expr
	Body        Stmt
}

func (f *ForStmt) String() string {
	s := "for ("
	if f.Initializer != nil {
		s += f.Initializer.String() + " "
	}
	s += "; "
	if f.Condition != nil {
		s += f.Condition.String()
	}
	s += "; "
	if f.Increment != nil {
		s += f.Increment.String()
	}
	s += ") " + f.Body.String()
	return s
}

type CallExpr struct {
	Callee    Expr
	Arguments []Expr
}

func (c *CallExpr) String() string {
	args := ""
	for i, arg := range c.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.String()
	}
	return fmt.Sprintf("%s(%s)", c.Callee.String(), args)
}

type FunctionStmt struct {
	Name   string
	Params []string
	Body   *BlockStmt
}

func (f *FunctionStmt) String() string {
	return fmt.Sprintf("fun %s(%v) %s", f.Name, f.Params, f.Body.String())
}

type ReturnStmt struct {
	Keyword Token // The return keyword (to carry line info, etc.)
	Value   Expr  // The expression being returned (can be nil for no value)
}

func (r *ReturnStmt) String() string {
	if r.Value != nil {
		return fmt.Sprintf("return %s", r.Value.String())
	}
	return "return"
}

type ReturnValue struct {
	Value interface{}
}

type ClassStmt struct {
	Name    string
	Methods []*FunctionStmt // For now, the body is just a list of method declarations.
}

// String returns a string representation of the class statement.
func (c *ClassStmt) String() string {
	return fmt.Sprintf("class %s { ... }", c.Name)
}

type Get struct {
	Object Expr  // The instance expression.
	Name   Token // The property name token.
}

func (g *Get) String() string {
	return fmt.Sprintf("(%s.%s)", g.Object.String(), g.Name.Lexeme)
}

type Set struct {
	Object Expr  // The instance expression.
	Name   Token // The property name token.
	Value  Expr  // The value to assign.
}

func (s *Set) String() string {
	return fmt.Sprintf("(%s.%s = %s)", s.Object.String(), s.Name.Lexeme, s.Value.String())
}
