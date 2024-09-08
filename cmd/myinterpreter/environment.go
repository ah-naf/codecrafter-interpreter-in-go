package main

import "fmt"

type Environment struct {
	Values map[string]interface{}
}

// Creates a new environment
func NewEnvironment() *Environment {
	return &Environment{Values: make(map[string]interface{})}
}

// Define a new variable in environment
func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

// Get the value of environment
func(e *Environment) Get(name string) (interface{}, error) {
	if value, exists := e.Values[name]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("undefined variable '%s'", name)
}