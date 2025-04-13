package main

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
