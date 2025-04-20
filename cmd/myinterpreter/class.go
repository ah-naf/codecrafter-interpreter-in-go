package main

import (
	"fmt"
	"os"
)

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]*UserFunction
}

// String returns the class name when printed.
func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) Arity() int {
	if initializer, ok := c.Methods["init"]; ok {
		return initializer.Arity()
	}
	return 0
}

func (c *LoxClass) Call(env *Environment, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Klass:  c,
		Fields: make(map[string]interface{}),
	}

	initializer := c.findInitializer()
	if initializer != nil {
		boundInit := initializer.Bind(instance)
		boundInit.Call(env, arguments)
	}

	return instance
}

func (c *LoxClass) findInitializer() *UserFunction {
	if init, ok := c.Methods["init"]; ok {
		return init
	}
	if c.Superclass != nil {
		return c.Superclass.findInitializer()
	}
	return nil
}

func (c *LoxClass) FindMethod(name string) (*UserFunction, bool) {
	if method, ok := c.Methods[name]; ok {
		return method, true
	}
	if c.Superclass != nil {
		return c.Superclass.FindMethod(name)
	}
	return nil, false
}

type LoxInstance struct {
	Klass  *LoxClass
	Fields map[string]interface{}
}

// String returns a string representation of an instance.
func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.Klass.Name)
}

func (inst *LoxInstance) Get(name Token) (interface{}, bool) {
	if value, ok := inst.Fields[name.Lexeme]; ok {
		return value, true
	}

	if method, ok := inst.Klass.FindMethod(name.Lexeme); ok {
		// Bind the method to this instance.
		return method.Bind(inst), true
	}
	fmt.Fprintf(os.Stderr, "Undefined property '%s'.\n", name.Lexeme)
	os.Exit(70)
	return nil, false
}

func (inst *LoxInstance) Set(name Token, value interface{}) {
	inst.Fields[name.Lexeme] = value
}
