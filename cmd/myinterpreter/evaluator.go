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
	// If it's a number string, convert it
	if l.Type == "number" {
		if num, ok := l.Value.(string); ok && isNumber(num) {
			value, err := ConvertStringToFloat(num, 0) // 0 for line number as it's a literal
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
			return value
		}
	}

	return l.Value // Return the literal value (true, false, or nil)
}

// Eval method for Grouping evaluates the inner expression
func (g *Grouping) Eval() interface{} {
	return g.Expression.Eval() // Evaluate the expression inside parentheses
}

// Eval method for Unary handles unary operators like ! and -
func (u *Unary) Eval() interface{} {
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
	leftVal := b.Left.Eval()
	rightVal := b.Right.Eval()
	
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
		raiseBinaryTypeError(b.Line, leftVal, rightVal, "+")
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
		raiseBinaryTypeError(b.Line, leftVal, rightVal, "/")
	case GT:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if(leftIsNum && rightIsNum) {
			return leftNum > rightNum
		}

		raiseBinaryTypeError(b.Line, leftVal, rightVal, GT)
	case LT:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if(leftIsNum && rightIsNum) {
			return leftNum < rightNum
		}

		raiseBinaryTypeError(b.Line, leftVal, rightVal, LT)
	case GREATER_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if(leftIsNum && rightIsNum) {
			return leftNum >= rightNum
		}

		raiseBinaryTypeError(b.Line, leftVal, rightVal, GREATER_EQUAL)
	case LESS_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if(leftIsNum && rightIsNum) {
			return leftNum <= rightNum
		}

		raiseBinaryTypeError(b.Line, leftVal, rightVal, LESS_EQUAL)
	case BANG_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if(leftIsNum && rightIsNum) {
			return leftNum != rightNum
		} else if !leftIsNum && !rightIsNum {
			return leftVal != rightVal
		}

		raiseBinaryTypeError(b.Line, leftVal, rightVal, BANG_EQUAL)
	case EQUAL_EQUAL:
		leftNum, leftIsNum := toNumber(leftVal)
		rightNum, rightIsNum := toNumber(rightVal)

		if(leftIsNum && rightIsNum) {
			return leftNum == rightNum
		} else if !leftIsNum && !rightIsNum {
			return leftVal == rightVal
		}

		raiseBinaryTypeError(b.Line, leftVal, rightVal, EQUAL_EQUAL)
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
	raiseBinaryTypeError(line, leftVal, rightVal, operator)
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
	case float64:
		return v, true
	default:
		return 0, false
	}
}


// Helper function to raise a type error for binary operations
func raiseBinaryTypeError(line int, leftVal, rightVal interface{}, operator string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: Cannot apply '%s' to %T and %T\n", line, operator, leftVal, rightVal)
	os.Exit(1)
}