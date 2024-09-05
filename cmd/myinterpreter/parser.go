// parser.go
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
		lexer: lexer, // Use the lexer directly
		pos:   0,
	}
}

// Parse starts parsing and returns the resulting AST
func (p *Parser) Parse() Expr {
	if p.match("TRUE") {
		return &Literal{Value: true}
	} else if p.match("FALSE") {
		return &Literal{Value: false}
	} else if p.match("NIL") {
		return &Literal{Value: nil}
	} else if p.match("NUMBER") {
		// Convert the token literal (string) to a float64 value
		value := p.previous().Literal
		return &Literal{Value: value} // Return number literal as float64
	} else if p.match("STRING") {
		value := p.previous().Literal
		return &Literal{Value: value}
	} else if p.match("LEFT_PAREN") {
		expr := p.Parse() // Recursively parse the inner expression
		p.consume("RIGHT_PAREN", "Expect ')' after expression.")
		return &Grouping{Expression: expr} // Return the grouping expression
	} else if p.match("BANG") { // Logical NOT operator
		operator := p.previous()
		right := p.Parse() // Recursively parse the right-hand expression
		return &Unary{Operator: operator, Right: right}
	} else if p.match("MINUS") { // Negation operator
		operator := p.previous()
		right := p.Parse() // Recursively parse the right-hand expression
		return &Unary{Operator: operator, Right: right}
	}

	// If it's an unexpected token, we report an error
	p.error("Expected literal")
	return nil
}

// match checks if the current token matches the expected type
func (p *Parser) match(expectedType string) bool {
	if p.pos >= len(p.lexer.tokens) { // Directly access lexer.tokens
		return false
	}
	if p.lexer.tokens[p.pos].Type == expectedType {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) previous() Token {
	return p.lexer.tokens[p.pos-1]
}

// consume checks for a specific token and advances, or throws an error if it doesn't match
func (p *Parser) consume(expectedType, errorMessage string) {
	if !p.match(expectedType) {
		p.error(errorMessage)
	}
}

// error reports a parsing error
func (p *Parser) error(msg string) {
	token := p.lexer.tokens[p.pos] // Directly access lexer.tokens
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, msg)
	os.Exit(1)
}
