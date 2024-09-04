// ast.go
package main

import "fmt"

type Expr interface {
	String() string
}

type Literal struct {
	Value interface{}
}

// String method converts the literal value to its string representation
func (l *Literal) String() string {
	// Special case for nil
	if l.Value == nil {
		return "nil"
	}
	
	return fmt.Sprintf("%v", l.Value)
}
