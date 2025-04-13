// evaluator.go
package main

import (
	"fmt"
	"os"
	"strconv"
)

func (c *CallExpr) Eval(env *Environment) interface{} {
	calleval := c.Callee.Eval(env)

	callable, ok := calleval.(Callable)
	if !ok {
		fmt.Fprintln(os.Stderr, "Can only call functions and classes.")
		os.Exit(70)
	}

	var arguments []interface{}
	for _, arg := range c.Arguments {
		arguments = append(arguments, arg.Eval(env))
	}

	if len(arguments) != callable.Arity() {
		fmt.Fprintf(os.Stderr, "Expected %d arguments but got %d.\n", callable.Arity(), len(arguments))
		os.Exit(70)
	}

	return callable.Call(env, arguments)
}

func (f *ForStmt) Eval(env *Environment) interface{} {
	if f.Initializer != nil {
		f.Initializer.Eval(env)
	}

	for isTruthy(f.Condition.Eval(env)) {
		f.Body.Eval(env)

		if f.Increment != nil {
			f.Increment.Eval(env)
		}
	}

	return nil
}

func (w *WhileStmt) Eval(env *Environment) interface{} {
	for isTruthy(w.Condition.Eval(env)) {
		localEnv := NewEnvironmentWithParent(env)
		for _, stmt := range w.Body {
			stmt.Eval(localEnv)
		}
	}

	return nil
}

func (i *IfStmt) Eval(env *Environment) interface{} {
	condition := i.Condition.Eval(env)
	if isTruthy(condition) {
		localEnv := NewEnvironmentWithParent(env)
		var result interface{}
		for _, stmt := range i.Body {
			result = stmt.Eval(localEnv)
		}
		return result
	} else {
		localEnv := NewEnvironmentWithParent(env)
		var result interface{}
		for _, stmt := range i.Else {
			result = stmt.Eval(localEnv)
		}
		return result
	}
}

func (b *BlockStmt) Eval(env *Environment) interface{} {
	localEnv := NewEnvironmentWithParent(env)
	for _, stmt := range b.Statements {
		stmt.Eval(localEnv)
	}
	return nil
}

func (a *AssignStmt) Eval(env *Environment) interface{} {
	value := a.Value.Eval(env)
	if err := env.Assign(a.Name, value); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(70)
	}
	return value
}

func (v *VarStmt) Eval(env *Environment) interface{} {
	var value interface{}
	if v.Initializer != nil {
		value = v.Initializer.Eval(env)
	}
	if v.VarUsed {
		env.Define(v.Name, value)
	} else {
		if err := env.Assign(v.Name, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(70)
		}
	}
	return "nil"
}

func (i *Identifier) Eval(env *Environment) interface{} {
	value, err := env.Get(i.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Undefined variable '%s'.\n[line %d]\n", i.Name, i.Line)
		os.Exit(70)
	}
	if value == nil {
		value = "nil"
	}
	return value
}

func (p *PrintStatement) Eval(env *Environment) interface{} {
	value := p.Expression.Eval(env)
	fmt.Println(value)
	return nil
}

func (e *ExpressionStatement) Eval(env *Environment) interface{} {
	return e.Expression.Eval(env)
}

func (l *Literal) Eval(env *Environment) interface{} {
	if l.Value == nil {
		return "nil"
	}
	if l.Type == "number" {
		if num, ok := l.Value.(string); ok && isNumber(num) {
			value, err := ConvertStringToFloat(num, 0)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
			return value
		}
	}
	return l.Value
}

// Eval method for Grouping evaluates the inner expression
func (g *Grouping) Eval(env *Environment) interface{} {
	return g.Expression.Eval(env) // Evaluate the expression inside parentheses
}

// Eval method for Unary handles unary operators like ! and -
func (u *Unary) Eval(env *Environment) interface{} {
	rightVal := u.Right.Eval(env) // Evaluate the right-hand expression

	switch u.Operator.Type {
	case "BANG": // Logical NOT
		return !isTruthy(rightVal)

	case "MINUS": // Negation
		switch num := rightVal.(type) {
		case float64:
			return -num // Negate the float64 number
		case int:
			return -float64(num) // Negate the integer by converting to float64
		case string:
			// Use the ConvertStringToFloat function to handle string conversion and error reporting
			value, err := ConvertStringToFloat(num, u.Line)
			if err != nil {
				raiseRuntimeError(u.Line)
			}
			return -value // Negate the converted float64 number
		default:
			raiseRuntimeError(u.Line)
		}
	}

	return nil
}

// Eval method for Binary expressions (for future operators)
func (b *Binary) Eval(env *Environment) interface{} {
	leftVal := b.Left.Eval(env)
	rightVal := b.Right.Eval(env)

	switch b.Operator.Lexeme {
	case PLUS: // Handle addition
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)
		// Handle addition of numbers
		if leftIsNum && rightIsNum {
			return leftNum + rightNum
		}

		// Handle string concatenation
		if leftStr, ok := leftVal.(string); ok {
			if rightStr, ok := rightVal.(string); ok {
				return leftStr + rightStr // Concatenate two strings
			}
		}

		// Raise an error for incompatible types
		raiseRuntimeError(b.Line)
	case MINUS: // Handle subtraction
		return handleBinaryNumberOperation(leftVal, rightVal, "-", b.Line)
	case STAR: // Handle multiplication
		return handleBinaryNumberOperation(leftVal, rightVal, "*", b.Line)
	case SLASH: // Handle division
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			// Check for division by zero
			if rightNum == 0 {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Cannot divide by zero\n", b.Line)
				os.Exit(1)
			}
			return leftNum / rightNum
		}

		// Raise an error for incompatible types
		raiseRuntimeError(b.Line)
	case GT:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			return leftNum > rightNum
		}

		raiseRuntimeError(b.Line)
	case LT:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			return leftNum < rightNum
		}

		raiseRuntimeError(b.Line)
	case GREATER_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			return leftNum >= rightNum
		}

		raiseRuntimeError(b.Line)
	case LESS_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			return leftNum <= rightNum
		}

		raiseRuntimeError(b.Line)
	case BANG_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			return leftNum != rightNum
		} else if !leftIsNum && !rightIsNum {
			return leftVal != rightVal
		}

		return false
	case EQUAL_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if leftIsNum && rightIsNum {
			return leftNum == rightNum
		} else if !leftIsNum && !rightIsNum {
			return leftVal == rightVal
		}

		return false
	case "or":
		if isTruthy(b.Left.Eval(env)) {
			return leftVal
		}
		return rightVal
	case "and":
		if !isTruthy(b.Left.Eval(env)) {
			return leftVal
		}
		return rightVal
	}

	return nil
}

// Helper function to handle number operations (+, -, *, /) for binary expressions
func handleBinaryNumberOperation(leftVal, rightVal interface{}, operator string, line int) interface{} {
	leftNum, leftIsNum := toNumber(leftVal)
	rightNum, rightIsNum := toNumber(rightVal)

	if leftIsNum && rightIsNum {
		switch operator {
		case "-":
			return leftNum - rightNum
		case "*":
			return leftNum * rightNum
		}
	}

	// Raise an error for incompatible types
	raiseRuntimeError(line)
	return nil
}

// Helper function to check truthiness (used in logical NOT)
func isTruthy(value interface{}) bool {
	if value == nil || value == "nil" {
		return false
	}
	if boolean, ok := value.(bool); ok {
		return boolean
	}
	return true
}

// Helper function to check if a string is a number
func isNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

// ConvertStringToFloat tries to convert a string to a float64.
// If it fails, it returns an error with the associated line number.
func ConvertStringToFloat(input string, line int) (float64, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("[line %d] Error: Invalid number format '%s'", line, input)
	}
	return value, nil
}

// Helper function to convert interface{} to float64 for number operations
// Returns a float64 and a boolean indicating if the conversion was successful.
func toNumber(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// Helper function to raise a type error for binary operations
func raiseRuntimeError(line int) {
	fmt.Fprintf(os.Stderr, "Operands must be a number.\n[line %d]", line)
	os.Exit(70)
}
