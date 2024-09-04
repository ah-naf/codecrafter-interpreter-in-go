package main

import (
	"fmt"
	"os"
)

type Parser struct {
	lexer      *Lexer
	currentTok string
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.advance() // Initialize first token
	return p
}

func (p *Parser) advance() {
	if p.lexer.ch != 0 {
		p.lexer.readChar()
	}
	// Assign the next token (for simplicity, assume next token is a keyword or literal)
	p.currentTok = string(p.lexer.ch)
}

func (p *Parser) Parse() Expr {
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() Expr {
	switch p.currentTok {
	case "true":
		return &Literal{Value: true}
	case "false":
		return &Literal{Value: false}
	case "nil":
		return &Literal{Value: nil}
	default:
		p.reportError()
	}
	return nil
}

func (p *Parser) reportError() {
	fmt.Fprintln(os.Stderr, "Error: Unexpected input")
	os.Exit(1)
}
