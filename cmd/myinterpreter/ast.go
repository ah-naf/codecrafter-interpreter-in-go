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
