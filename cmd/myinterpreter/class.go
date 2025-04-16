package main

import (
	"fmt"
	"os"
)

type LoxClass struct {
	Name    string
	Methods map[string]*UserFunction
}

// String returns the class name when printed.
func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) Call(env *Environment, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Klass:  c,
		Fields: make(map[string]interface{}),
		// A fields table could be added here later.
	}
	// If an initializer method were defined (commonly named "init"),
	// you would look it up and call it here.
	return instance
}

type LoxInstance struct {
	Klass  *LoxClass
	Fields map[string]interface{}
	// In a more advanced implementation, instance fields would be stored here.
}

// String returns a string representation of an instance.
func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.Klass.Name)
}

func (inst *LoxInstance) Get(name Token) (interface{}, bool) {
	if value, ok := inst.Fields[name.Lexeme]; ok {
		return value, true
	}

	if method, ok := inst.Klass.Methods[name.Lexeme]; ok {
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
