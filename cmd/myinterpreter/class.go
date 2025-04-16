package main

type LoxClass struct {
	Name    string
	Methods map[string]*UserFunction
}

// String returns the class name when printed.
func (c *LoxClass) String() string {
	return c.Name
}
