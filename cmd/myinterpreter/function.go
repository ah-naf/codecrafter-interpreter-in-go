package main

import "fmt"

type Callable interface {
	Arity() int
	Call(env *Environment, arguments []interface{}) interface{}
	String() string
}

type NativeFunction struct {
	arity    int
	function func(arguments []interface{}) interface{}
}

func (nf *NativeFunction) Arity() int {
	return nf.arity
}

func (nf *NativeFunction) Call(env *Environment, arguments []interface{}) interface{} {
	return nf.function(arguments)
}

func (nf *NativeFunction) String() string {
	return "<native fn>"
}

type UserFunction struct {
	Declaration *FunctionStmt
	Closure     *Environment
}

func (uf *UserFunction) Arity() int {
	return len(uf.Declaration.Params)
}

func (uf *UserFunction) String() string {
	return fmt.Sprintf("<fn %s>", uf.Declaration.Name)
}

func (uf *UserFunction) Call(env *Environment, arguments []interface{}) interface{} {
	// Create a new environment that uses the closure (the defining environment) as parent.
	localEnv := NewEnvironmentWithParent(uf.Closure)
	// Bind each parameter to the corresponding argument.
	for i, param := range uf.Declaration.Params {
		localEnv.Define(param, arguments[i])
	}
	// Evaluate the function body.
	var result interface{}
	for _, stmt := range uf.Declaration.Body.Statements {
		result = stmt.Eval(localEnv)
	}
	return result
}