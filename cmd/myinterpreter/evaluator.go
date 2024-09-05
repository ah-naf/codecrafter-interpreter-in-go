// evaluator.go
package main

// ExprEvaluator is an interface for expressions that can be evaluated
type ExprEvaluator interface {
	Eval() interface{} // Method to evaluate the expression
}

// Eval method for Literal evaluates and returns the value of the literal
func (l *Literal) Eval() interface{} {
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
		if num, ok := rightVal.(float64); ok {
			return -num
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
