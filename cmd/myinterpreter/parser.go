// parser.go
package main

import (
	"fmt"
	"os"
)

type Parser struct {
	lexer  *Lexer
	tokens []Token
	pos    int
}

// NewParser initializes a new parser with the lexer input
func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer:  lexer,
		tokens: lexer.tokens,
		pos:    0,
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
	}

	// If it's an unexpected token, we report an error
	p.error("Expected literal")
	return nil
}

// match checks if the current token matches the expected type
func (p *Parser) match(expectedType string) bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	if p.tokens[p.pos].Type == expectedType {
		p.pos++
		return true
	}
	return false
}

// error reports a parsing error
func (p *Parser) error(msg string) {
	token := p.tokens[p.pos]
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, msg)
	os.Exit(1)
}
