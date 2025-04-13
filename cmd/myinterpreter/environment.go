package main

import "fmt"

type Environment struct {
	Values map[string]interface{}
	Parent *Environment
}

// Creates a new environment
func NewEnvironment() *Environment {
	return &Environment{Values: make(map[string]interface{})}
}

// NewEnvironmentWithParent creates a new environment with a reference to a parent environment
func NewEnvironmentWithParent(parent *Environment) *Environment {
	return &Environment{
		Values: make(map[string]interface{}),
		Parent: parent,
	}
}

func (e *Environment) Assign(name string, value interface{}) error {
	if _, exists := e.Values[name]; exists {
		e.Values[name] = value
		return nil
	}
	if e.Parent != nil {
		return e.Parent.Assign(name, value)
	}
	return fmt.Errorf("undefined variable '%s'", name)
}

// Define a new variable in environment
func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

// Get the value of a variable, checking parent scopes if necessary
func (e *Environment) Get(name string) (interface{}, error) {
	if value, exists := e.Values[name]; exists {
		return value, nil
	}

	if e.Parent != nil {
		return e.Parent.Get(name)
	}

	return nil, fmt.Errorf("undefined variable '%s'", name)
}
