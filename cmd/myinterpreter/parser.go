package main

import (
	"fmt"
	"os"
)

type Parser struct {
	lexer *Lexer
	pos   int
	mode  string
}

// NewParser initializes a new parser with the lexer input
func NewParser(lexer *Lexer, mode string) *Parser {
	return &Parser{
		lexer: lexer,
		pos:   0,
		mode:  mode,
	}
}

// Parse starts parsing and returns the resulting AST
func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}

	for !p.isAtEnd() {
		statements = append(statements, p.parseStatement())
	}

	return statements
}

// parseStatement handles either print statements or expression statements
func (p *Parser) parseStatement() Stmt {
	if p.match("PRINT") {
		return p.printStatement()
	} else if p.match("VAR") {
		return p.varDeclaration()
	}
	return p.expressionStatement()
}

// varDeclaration parses a variable declaration
func (p *Parser) varDeclaration() Stmt {
	// Expect an identifier after 'var'
	p.consume("IDENTIFIER", "Expect variable name.")
	identifier := p.previous()
	var initializer Expr
	if p.match("EQUAL") { // If '=' follows, there should be an initializer expression
		initializer = p.parseEquality()
	}

	// Ensure there's a semicolon after the variable declaration
	p.consume("SEMICOLON", "Expect ';' after variable declaration.")
	return &VarStmt{Name: identifier.Lexeme, Initializer: initializer}
}

// printStatement parses a print statement
func (p *Parser) printStatement() Stmt {
	expr := p.parseEquality() // Parse the expression after "print"
	if p.mode == "run" {
		if !p.checkSemicolon() {
			fmt.Fprintf(os.Stderr, "[line %d]: Expect ';' after expression\n", p.previous().Line)
			os.Exit(65)
		}
		p.consume("SEMICOLON", "Expect ';' after expression.")
	}
	return &PrintStatement{Expression: expr} // Return a PrintStatement node
}

// expressionStatement parses an expression statement
func (p *Parser) expressionStatement() Stmt {
	expr := p.parseEquality() // Parse the expression
	if p.mode == "run" {
		if !p.checkSemicolon() {
			fmt.Fprintf(os.Stderr, "[line %d]: Expect ';' after expression", p.previous().Line)
			os.Exit(65)
		}
		p.consume("SEMICOLON", "Expect ';' after expression.")
	}
	return &ExpressionStatement{Expression: expr} // Return an expression statement
}

func(p *Parser) parseEquality() Stmt {
	expr := p.parseComparison()

	for p.match("EQUAL_EQUAL", "BANG_EQUAL") {
		operator := p.previous()
		right := p.parseComparison()
		expr = &Binary{Left: expr, Operator: operator, Right: right, Line: operator.Line}
	}

	return expr
}

// parseComparison handles >, <, >=, <= operators
func (p *Parser) parseComparison() Expr {
	expr := p.parseAdditionSubstraction() // Start by parsing addition and subtraction

	for p.match("GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL") { // Look for comparison operators
		operator := p.previous()
		right := p.parseAdditionSubstraction() // Parse the right-hand operand
		expr = &Binary{Left: expr, Operator: operator, Right: right, Line: operator.Line}
	}

	return expr
}

// parseAdditionSubstraction handles + and - operators
func (p *Parser) parseAdditionSubstraction() Expr {
	expr := p.parseMultiplication()

	for p.match("PLUS", "MINUS") {
		operator := p.previous()
		right := p.parseMultiplication()
		expr = &Binary{Left: expr, Operator: operator, Right: right, Line: operator.Line}
	}

	return expr
}

// parseMultiplication handles * and / operators with their precedence
func (p *Parser) parseMultiplication() Expr {
	expr := p.parseUnary() // Start by parsing unary operators

	for p.match("STAR", "SLASH") { // Look for * or / operators
		operator := p.previous()
		right := p.parseUnary() // Parse the right-hand operand (which could be a unary expression)
		expr = &Binary{Left: expr, Operator: operator, Right: right, Line: operator.Line}
	}

	return expr
}

// parseUnary handles unary operators (e.g., -23, !true) or forwards to primary expressions
func (p *Parser) parseUnary() Expr {
	if p.match("BANG", "MINUS") { // Check for the unary operators
		operator := p.previous()
		right := p.parseUnary() // Recursively parse the right-hand operand
		return &Unary{Operator: operator, Right: right, Line: operator.Line}
	}

	// If it's not a unary expression, parse a primary expression
	return p.parsePrimary()
}

// parsePrimary handles numbers, strings, booleans, and parentheses
func (p *Parser) parsePrimary() Expr {
	switch {
	case p.match("TRUE"):
		return &Literal{Value: true, Type: "boolean"}
	case p.match("FALSE"):
		return &Literal{Value: false, Type: "boolean"}
	case p.match("NIL"):
		return &Literal{Value: nil, Type: "nil"}
	case p.match("NUMBER"):
		return &Literal{Value: p.previous().Literal, Type: "number"}
	case p.match("STRING"):
		return &Literal{Value: p.previous().Literal, Type: "string"}
	case p.match("IDENTIFIER"):
		return &Identifier{Name: p.previous().Lexeme}
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

// checkSemicolon checks if the current token is a semicolon without advancing the position
func (p *Parser) checkSemicolon() bool {
	if p.isAtEnd() {
		return false
	}
	return p.lexer.tokens[p.pos].Type == "SEMICOLON"
}


func (p *Parser) error(msg string) {
	if p.pos < len(p.lexer.tokens) {
		token := p.lexer.tokens[p.pos]
		fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, msg)
	} else {
		// Handle the case where the token list is exhausted
		fmt.Fprintf(os.Stderr, "[line %d] Error at end: %s\n", p.lexer.line, msg)
	}
	os.Exit(65)
}