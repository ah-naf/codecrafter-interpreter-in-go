// ast.go
package main

import "fmt"

type Expr interface {
	String() string
}

type Literal struct {
	Value interface{}
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) String() string {
	// fmt.Println(g.Expression)
	return fmt.Sprintf("(group %s)", g.Expression.String())
}

// String method converts the literal value to its string representation
func (l *Literal) String() string {
	// Special case for nil
	if l.Value == nil {
		return "nil"
	}
	
	return fmt.Sprintf("%v", l.Value)
}

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

// String method for Binary to print its content
func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}