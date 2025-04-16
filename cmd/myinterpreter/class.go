package main

import "fmt"

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

type LoxInstance struct {
	Klass *LoxClass
	// In a more advanced implementation, instance fields would be stored here.
}

// String returns a string representation of an instance.
func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.Klass.Name)
}


func (c *LoxClass) Call(env *Environment, arguments []interface{}) interface{} {
    instance := &LoxInstance{
        Klass: c,
        // A fields table could be added here later.
    }
    // If an initializer method were defined (commonly named "init"),
    // you would look it up and call it here.
    return instance
}