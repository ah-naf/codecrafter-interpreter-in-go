package main

import "fmt"

type Expr interface {
	String() string
}

type Literal struct {
	Value interface{}
}

func (l *Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}
