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
	Type string
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
	Line 	 int
}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Right.String())
}

// Binary struct to represent binary expressions (e.g., 16 * 38)
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
	Line	 int
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
}

func (v *VarStmt) String() string {
	return fmt.Sprintf("var %s = %v", v.Name, v.Initializer)
}

// Identifier represents a variable being used in an expression
type Identifier struct {
	Name string
}

func (i *Identifier) String() string {
	return i.Name
}