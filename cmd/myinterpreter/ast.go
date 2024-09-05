// ast.go
package main

import "fmt"

// Expr interface for all expression nodes, extended to include ExprEvaluator
type Expr interface {
	String() string
	ExprEvaluator // Include evaluation in the expression interface
}

// Literal struct for literal values (booleans, numbers, strings, nil)
type Literal struct {
	Value interface{}
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
}

func (u *Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Right.String())
}

// Binary struct to represent binary expressions (e.g., 16 * 38)
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}
