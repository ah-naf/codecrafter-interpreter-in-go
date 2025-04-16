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

func (fn *UserFunction) Bind(instance *LoxInstance) *UserFunction {
    env := NewEnvironmentWithParent(fn.Closure)
    env.Define("this", instance)
    return &UserFunction{
        Declaration: fn.Declaration,
        Closure:     env,
    }
}

func (uf *UserFunction) Call(env *Environment, arguments []interface{}) interface{} {
	// Create a new environment that uses the closure (the defining environment) as parent.
	localEnv := NewEnvironmentWithParent(uf.Closure)
	for i, param := range uf.Declaration.Params {
		localEnv.Define(param, arguments[i])
	}

	var returnValue interface{} = nil
	func() {
		defer func() {
			if r := recover(); r != nil {
				if rv, ok := r.(ReturnValue); ok {
					returnValue = rv.Value
				} else {
					// Re-panic if it's not a return value.
					panic(r)
				}
			}
		}()
		// Evaluate the function body.
		for _, stmt := range uf.Declaration.Body.Statements {
			stmt.Eval(localEnv)
		}
	}()

	return returnValue
}
