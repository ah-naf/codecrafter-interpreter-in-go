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

// NewParser initializes a new parser with the lexer input.
func NewParser(lexer *Lexer, mode string) *Parser {
	return &Parser{
		lexer: lexer,
		pos:   0,
		mode:  mode,
	}
}

// Parse starts parsing and returns the resulting AST.
func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.parseStatement())
	}
	return statements
}

// parseStatement handles print statements, declarations, assignments, block statements, if statements, etc.
func (p *Parser) parseStatement() Stmt {
	if p.match("PRINT") {
		return p.printStatement()
	} else if p.match("VAR") {
		return p.varDeclaration()
	} else if p.match("IDENTIFIER") {
		return p.varAssignment()
	} else if p.match("LEFT_BRACE") {
		return p.blockStatement()
	} else if p.match("IF") {
		return p.ifStatement()
	} else if p.match("WHILE") {
		return p.whileStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) whileStatement() Stmt {
	p.consume("LEFT_PAREN", "Expect '(' after if statement")
	condition := p.parseAssignment()
	p.consume("RIGHT_PAREN", "Expect ')' after if condition expression")

	var body []Stmt
	if p.check("LEFT_BRACE") {
		p.consume("LEFT_BRACE", "Expect '{' after while")
		for !p.isAtEnd() && !p.check("RIGHT_BRACE") {
			body = append(body, p.parseStatement())
		}
		p.consume("RIGHT_BRACE", "Expect '}' after block.")
	} else {
		body = []Stmt{p.parseStatement()}
	}

	return &WhileStmt{
		Condition: condition,
		Body: body,
	}
}

func (p *Parser) ifStatement() Stmt {
	// Parse condition.
	p.consume("LEFT_PAREN", "Expect '(' after if statement")
	condition := p.parseAssignment()
	p.consume("RIGHT_PAREN", "Expect ')' after if condition expression")

	// Parse "if" branch.
	var ifBranch []Stmt
	if p.check("LEFT_BRACE") {
		p.consume("LEFT_BRACE", "Expect '{' after if")
		for !p.isAtEnd() && !p.check("RIGHT_BRACE") {
			ifBranch = append(ifBranch, p.parseStatement())
		}
		p.consume("RIGHT_BRACE", "Expect '}' after block.")
	} else {
		stmt := p.parseStatement()
		ifBranch = []Stmt{stmt}
	}

	// Optional: Parse "else" branch.
	var elseBranch []Stmt
	if p.match("ELSE") {
		if p.check("LEFT_BRACE") {
			p.consume("LEFT_BRACE", "Expect '{' after else")
			for !p.isAtEnd() && !p.check("RIGHT_BRACE") {
				elseBranch = append(elseBranch, p.parseStatement())
			}
			p.consume("RIGHT_BRACE", "Expect '}' after block in else")
		} else {
			stmt := p.parseStatement()
			elseBranch = []Stmt{stmt}
		}
	}

	return &IfStmt{
		Condition: condition,
		Body:      ifBranch,
		Else:      elseBranch,
	}
}

// blockStatement parses a block of statements enclosed in braces.
func (p *Parser) blockStatement() Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() && !p.check("RIGHT_BRACE") {
		statements = append(statements, p.parseStatement())
	}
	p.consume("RIGHT_BRACE", "Expect '}' after block.")
	return &BlockStmt{Statements: statements}
}

// varAssignment parses a variable assignment.
func (p *Parser) varAssignment() Stmt {
	identifier := p.previous()
	var initializer Expr
	if p.match("EQUAL") {
		initializer = p.parseAssignment()
	}
	p.consume("SEMICOLON", "Expect ';' after variable declaration.")
	return &VarStmt{
		Name:        identifier.Lexeme,
		Initializer: initializer,
		VarUsed:     false,
		Line:        identifier.Line,
	}
}

// varDeclaration parses a variable declaration.
func (p *Parser) varDeclaration() Stmt {
	p.consume("IDENTIFIER", "Expect variable name.")
	identifier := p.previous()
	var initializer Expr
	if p.match("EQUAL") {
		initializer = p.parseAssignment()
	}
	p.consume("SEMICOLON", "Expect ';' after variable declaration.")
	return &VarStmt{
		Name:        identifier.Lexeme,
		Initializer: initializer,
		VarUsed:     true,
		Line:        identifier.Line,
	}
}

// printStatement parses a print statement.
func (p *Parser) printStatement() Stmt {
	expr := p.parseAssignment()
	if p.mode == "run" {
		if !p.checkSemicolon() {
			fmt.Fprintf(os.Stderr, "[line %d]: Expect ';' after expression\n", p.previous().Line)
			os.Exit(65)
		}
		p.consume("SEMICOLON", "Expect ';' after expression.")
	}
	return &PrintStatement{Expression: expr}
}

// expressionStatement parses an expression statement.
func (p *Parser) expressionStatement() Stmt {
	expr := p.parseAssignment()
	if p.mode == "run" {
		if !p.checkSemicolon() {
			fmt.Fprintf(os.Stderr, "[line %d]: Expect ';' after expression", p.previous().Line)
			os.Exit(65)
		}
		p.consume("SEMICOLON", "Expect ';' after expression.")
	}
	return &ExpressionStatement{Expression: expr}
}

// parseAssignment parses assignment expressions (lowest precedence).
func (p *Parser) parseAssignment() Stmt {
	expr := p.parseLogicalAND()
	for p.match("EQUAL") {
		equals := p.previous()
		value := p.parseAssignment()
		if identifier, ok := expr.(*Identifier); ok {
			return &AssignStmt{
				Name:  identifier.Name,
				Value: value,
				Line:  equals.Line,
			}
		}
		// Optionally, report an error for invalid assignment target.
	}
	return expr
}

// parseLogicalAND parses "and" expressions.
func (p *Parser) parseLogicalAND() Expr {
	expr := p.parseLogicalOR()
	for p.match("AND") {
		operator := p.previous()
		right := p.parseLogicalOR()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return expr
}

// parseLogicalOR parses "or" expressions.
func (p *Parser) parseLogicalOR() Expr {
	expr := p.parseEquality()
	for p.match("OR") {
		operator := p.previous()
		right := p.parseEquality()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return expr
}

// parseEquality parses equality expressions.
func (p *Parser) parseEquality() Stmt {
	expr := p.parseComparison()
	for p.match("EQUAL_EQUAL", "BANG_EQUAL") {
		operator := p.previous()
		right := p.parseComparison()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return expr
}

// parseComparison handles >, <, >=, <= operators.
func (p *Parser) parseComparison() Expr {
	expr := p.parseAdditionSubstraction()
	for p.match("GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL") {
		operator := p.previous()
		right := p.parseAdditionSubstraction()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return expr
}

// parseAdditionSubstraction handles + and - operators.
func (p *Parser) parseAdditionSubstraction() Expr {
	expr := p.parseMultiplication()
	for p.match("PLUS", "MINUS") {
		operator := p.previous()
		right := p.parseMultiplication()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return expr
}

// parseMultiplication handles * and / operators.
// (Note the change: we now call parseUnary() here to break the recursion cycle.)
func (p *Parser) parseMultiplication() Expr {
	expr := p.parseUnary()
	for p.match("STAR", "SLASH") {
		operator := p.previous()
		right := p.parseUnary()
		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return expr
}

// parseUnary handles unary operators or defers to primary expressions.
func (p *Parser) parseUnary() Expr {
	if p.match("BANG", "MINUS") {
		operator := p.previous()
		right := p.parseUnary()
		return &Unary{
			Operator: operator,
			Right:    right,
			Line:     operator.Line,
		}
	}
	return p.parsePrimary()
}

// parsePrimary handles numbers, strings, booleans, identifiers, and grouping.
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
		return &Identifier{
			Name: p.previous().Lexeme,
			Line: p.previous().Line,
		}
	case p.match("LEFT_PAREN"):
		expr := p.parseAssignment()
		p.consume("RIGHT_PAREN", "Expect ')' after expression.")
		return &Grouping{Expression: expr}
	default:
		p.error("Expected expression.")
		return nil
	}
}

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

func (p *Parser) previous() Token {
	return p.lexer.tokens[p.pos-1]
}

func (p *Parser) consume(expectedType, errorMessage string) {
	if !p.match(expectedType) {
		p.error(errorMessage)
	}
}

func (p *Parser) isAtEnd() bool {
	return p.pos >= len(p.lexer.tokens)
}

func (p *Parser) checkSemicolon() bool {
	if p.isAtEnd() {
		return false
	}
	return p.lexer.tokens[p.pos].Type == "SEMICOLON"
}

func (p *Parser) check(tokenType string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.lexer.tokens[p.pos].Type == tokenType
}

func (p *Parser) error(msg string) {
	if p.pos < len(p.lexer.tokens) {
		token := p.lexer.tokens[p.pos]
		fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, msg)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error at end: %s\n", p.lexer.line, msg)
	}
	os.Exit(65)
}
