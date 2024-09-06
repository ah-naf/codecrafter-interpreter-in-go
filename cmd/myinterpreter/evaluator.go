// evaluator.go
package main

import (
	"fmt"
	"os"
	"strconv"
)

// ExprEvaluator is an interface for expressions that can be evaluated
type ExprEvaluator interface {
	Eval() interface{} // Method to evaluate the expression
}

// Eval method for Literal evaluates and returns the value of the literal
func (l *Literal) Eval() interface{} {
	if l.Value == nil {
		return "nil"
	}
	return l.Value // Return the literal value (true, false, or nil)
}

// Eval method for Grouping evaluates the inner expression
func (g *Grouping) Eval() interface{} {
	return g.Expression.Eval() // Evaluate the expression inside parentheses
}

// Eval method for Unary handles unary operators like ! and -
func (u *Unary) Eval() interface{} {
	fmt.Printf("%#v\n", u.String())
	rightVal := u.Right.Eval() // Evaluate the right-hand expression

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
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
			return -value // Negate the converted float64 number
		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Operand must be a number\n", u.Line)
			return nil
		}
	}

	return nil
}

// Eval method for Binary expressions (for future operators)
func (b *Binary) Eval() interface{} {
	// leftVal := b.Left.Eval()
	// rightVal := b.Right.Eval()

	// switch b.Operator.Type {
	// // Add binary operator cases like +, -, *, /, etc., if needed
	// }

	return nil
}

// Helper function to check truthiness (used in logical NOT)
func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if boolean, ok := value.(bool); ok {
		return boolean
	}
	return true
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