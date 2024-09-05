package main

import (
	"fmt"
	"os"
)

type Parser struct {
	lexer *Lexer
	pos   int
}

// NewParser initializes a new parser with the lexer input
func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
		pos:   0,
	}
}

// Parse starts parsing and returns the resulting AST
func (p *Parser) Parse() Expr {
	return p.parseEquality()
}

func(p *Parser) parseEquality() Expr {
	expr := p.parseComparison()

	for p.match("EQUAL_EQUAL", "BANG_EQUAL") {
		operator := p.previous()
		right := p.parseComparison()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// parseComparison handles >, <, >=, <= operators
func (p *Parser) parseComparison() Expr {
	expr := p.parseAdditionSubstraction() // Start by parsing addition and subtraction

	for p.match("GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL") { // Look for comparison operators
		operator := p.previous()
		right := p.parseAdditionSubstraction() // Parse the right-hand operand
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// parseAdditionSubstraction handles + and - operators
func (p *Parser) parseAdditionSubstraction() Expr {
	expr := p.parseMultiplication()

	for p.match("PLUS", "MINUS") {
		operator := p.previous()
		right := p.parseMultiplication()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// parseMultiplication handles * and / operators with their precedence
func (p *Parser) parseMultiplication() Expr {
	expr := p.parseUnary() // Start by parsing unary operators

	for p.match("STAR", "SLASH") { // Look for * or / operators
		operator := p.previous()
		right := p.parseUnary() // Parse the right-hand operand (which could be a unary expression)
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// parseUnary handles unary operators (e.g., -23, !true) or forwards to primary expressions
func (p *Parser) parseUnary() Expr {
	if p.match("BANG", "MINUS") { // Check for the unary operators
		operator := p.previous()
		right := p.parseUnary() // Recursively parse the right-hand operand
		return &Unary{Operator: operator, Right: right}
	}

	// If it's not a unary expression, parse a primary expression
	return p.parsePrimary()
}

// parsePrimary handles numbers, strings, booleans, and parentheses
func (p *Parser) parsePrimary() Expr {
	switch {
	case p.match("TRUE"):
		return &Literal{Value: true}
	case p.match("FALSE"):
		return &Literal{Value: false}
	case p.match("NIL"):
		return &Literal{Value: nil}
	case p.match("NUMBER"):
		return &Literal{Value: p.previous().Literal}
	case p.match("STRING"):
		return &Literal{Value: p.previous().Literal}
	case p.match("LEFT_PAREN"):
		expr := p.parseEquality() // Recursively parse the inner expression inside parentheses
		p.consume("RIGHT_PAREN", "Expect ')' after expression.")
		return &Grouping{Expression: expr} // Directly return the expression, not a group node
	default:
		p.error("Expected expression.")
		return nil
	}
}

// match checks if the current token matches one of the expected types
func (p *Parser) match(types ...string) bool {
	if p.isAtEnd() {
		return false
	}
	for _, t := range types {
		if p.lexer.tokens[p.pos].Type == t {
			p.pos++
			return true
		}
	}
	return false
}

// previous returns the last matched token
func (p *Parser) previous() Token {
	return p.lexer.tokens[p.pos-1]
}

// consume checks for a specific token and advances, or throws an error if it doesn't match
func (p *Parser) consume(expectedType, errorMessage string) {
	if !p.match(expectedType) {
		p.error(errorMessage)
	}
}

// isAtEnd checks if the parser has reached the end of the token list
func (p *Parser) isAtEnd() bool {
	return p.pos >= len(p.lexer.tokens)
}

// error reports a parsing error
func (p *Parser) error(msg string) {
	token := p.lexer.tokens[p.pos]
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, msg)
	os.Exit(65)
}
