package main

import (
	"fmt"
	"os"
)

type Parser struct {
	lexer   *Lexer
	pos     int
	mode    string
	inClass bool
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
	} else if p.match("LEFT_BRACE") {
		return p.blockStatement()
	} else if p.match("IF") {
		return p.ifStatement()
	} else if p.match("WHILE") {
		return p.whileStatement()
	} else if p.match("FOR") {
		return p.forStatement()
	} else if p.match("FUN") {
		return p.functionDeclaration()
	} else if p.match("RETURN") {
		return p.returnStatement()
	} else if p.match("CLASS") { // <-- New branch for class declarations
		return p.classDeclaration()
	}
	return p.expressionStatement()
}

func (p *Parser) classDeclaration() Stmt {
	p.consume("IDENTIFIER", "Expect class name.")
	className := p.previous().Lexeme

	var superclass *Identifier
	if p.match("LESS") {
		p.consume("IDENTIFIER", "Expect superclass name.")
		superclass = &Identifier{
			Name: p.previous().Lexeme,
			Line: p.previous().Line,
		}
		if className == superclass.Name {
			p.customError("A class can't inherit from itself.", className, superclass.Line)
			os.Exit(65)
		}
	}

	p.consume("LEFT_BRACE", "Expect '{' before class body.")

	// Save the current state and mark that we're in a class.
	enclosingClass := p.inClass
	p.inClass = true
	var methods []*FunctionStmt
	for !p.check("RIGHT_BRACE") && !p.isAtEnd() {
		methods = append(methods, p.parseMethod())
	}
	p.consume("RIGHT_BRACE", "Expect '}' after class body.")
	// Restore the previous inClass state.
	p.inClass = enclosingClass

	return &ClassStmt{
		Name:       className,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (p *Parser) parseMethod() *FunctionStmt {
	// The method name is an identifier.
	p.consume("IDENTIFIER", "Expect method name.")
	name := p.previous().Lexeme

	// Parse the parameter list.
	p.consume("LEFT_PAREN", "Expect '(' after method name.")
	var parameters []string
	if !p.check("RIGHT_PAREN") {
		for {
			p.consume("IDENTIFIER", "Expect parameter name.")
			parameters = append(parameters, p.previous().Lexeme)
			if !p.match("COMMA") {
				break
			}
		}
	}
	p.consume("RIGHT_PAREN", "Expect ')' after parameters.")

	// Parse the method body as a block.
	p.consume("LEFT_BRACE", "Expect '{' before method body.")
	body := p.blockStatement().(*BlockStmt)
	return &FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr = nil
	if !p.check("SEMICOLON") {
		value = p.parseAssignment()
	}
	p.consume("SEMICOLON", "Expect ';' after return value.")
	return &ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) functionDeclaration() Stmt {
	p.consume("IDENTIFIER", "Expect function name.")
	name := p.previous().Lexeme

	p.consume("LEFT_PAREN", "Expect '(' after function name.")
	var parameters []string
	if !p.check("RIGHT_PAREN") {
		// Parse at least one parameter.
		for {
			p.consume("IDENTIFIER", "Expect parameter name.")
			parameters = append(parameters, p.previous().Lexeme)
			if !p.match("COMMA") {
				break
			}
		}
	}
	p.consume("RIGHT_PAREN", "Expect ')' after parameters.")

	p.consume("LEFT_BRACE", "Expect '{' before function body.")
	body := p.blockStatement().(*BlockStmt)

	return &FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume("LEFT_PAREN", "Expect '(' after 'for'.")

	if p.check("LEFT_BRACE") {
		p.customError("Expect expression.", "{", p.previous().Line+1)
		p.customError("Expect ';' after expression.", ")", p.previous().Line+1)
		os.Exit(65)
	}

	var initializer Stmt
	if p.match("SEMICOLON") {
		initializer = nil
	} else if p.match("VAR") {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	if p.check("LEFT_BRACE") {
		p.customError("Expect expression.", "{", p.previous().Line+1)
		p.customError("Expect ';' after expression.", ")", p.previous().Line+1)
		os.Exit(65)
	}

	var condition Expr
	if !p.check("SEMICOLON") {
		condition = p.parseAssignment()
	} else {
		condition = &Literal{Value: true, Type: "boolean"}
	}
	p.consume("SEMICOLON", "Expect ';' after loop condition.")

	if p.check("LEFT_BRACE") {
		p.customError("Expect expression.", "{", p.previous().Line+1)
		os.Exit(65)
	}
	var increment Expr
	if !p.check("RIGHT_PAREN") {
		increment = p.parseAssignment()
	}
	p.consume("RIGHT_PAREN", "Expect ')' after for clauses.")

	if p.check("VAR") {
		p.customError("Expect expression.", "var", p.previous().Line)
		os.Exit(65)
	}
	body := p.parseStatement()

	return &ForStmt{
		Initializer: initializer,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume("LEFT_PAREN", "Expect '(' after if statement")
	condition := p.parseAssignment()
	p.consume("RIGHT_PAREN", "Expect ')' after if condition expression")

	var body []Stmt
	if p.check("VAR") {
		p.customError("Expect expression.", "var", p.previous().Line)
		os.Exit(65)
	}

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
		Body:      body,
	}
}

func (p *Parser) ifStatement() Stmt {
	// Parse condition.
	p.consume("LEFT_PAREN", "Expect '(' after if statement")
	condition := p.parseAssignment()
	p.consume("RIGHT_PAREN", "Expect ')' after if condition expression")

	// Parse "if" branch.
	var ifBranch []Stmt

	if p.check("VAR") {
		p.customError("Expect expression.", "var", p.previous().Line)
		os.Exit(65)
	}
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
		if p.check("VAR") {
			p.customError("Expect expression.", "var", p.previous().Line)
			os.Exit(65)
		}
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

// varDeclaration parses a variable declaration.
func (p *Parser) varDeclaration() Stmt {
	p.consume("IDENTIFIER", "Expect variable name.")
	identifier := p.previous()
	var initializer Expr
	if p.match("EQUAL") {
		initializer = p.parseAssignment()
	} else {
		initializer = &Literal{Value: nil, Type: "nil"}
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
	if p.match("EQUAL") {
		equals := p.previous()
		value := p.parseAssignment()

		// Handle assignment to a property (e.g. object.property = value)
		if getExpr, ok := expr.(*Get); ok {
			return &Set{
				Object: getExpr.Object,
				Name:   getExpr.Name,
				Value:  value,
			}
		} else if identifier, ok := expr.(*Identifier); ok {
			// Assignment to a plain variable.
			return &AssignStmt{
				Name:  identifier.Name,
				Value: value,
				Line:  equals.Line,
			}
		}

		p.error("Invalid assignment target.")
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
	return p.parseCall()
}

func (p *Parser) parseCall() Expr {
	expr := p.parsePrimary()
	for {
		if p.match("LEFT_PAREN") {
			expr = p.finishCall(expr)
		} else if p.match("DOT") {
			p.consume("IDENTIFIER", "Expect property name after '.'.")
			name := p.previous()
			expr = &Get{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check("RIGHT_PAREN") {
		arguments = append(arguments, p.parseAssignment())
		for p.match("COMMA") {
			arguments = append(arguments, p.parseAssignment())
		}
	}
	p.consume("RIGHT_PAREN", "Expect ')' after arguments.")
	return &CallExpr{Callee: callee, Arguments: arguments}
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
	case p.match("THIS"):
		if !p.inClass {
			// p.error("Can't use 'this' outside of a class.")
			fmt.Fprintf(os.Stderr, "[line %d] Error at 'this': Can't use 'this' outside of a class.", p.lexer.tokens[p.pos].Line)
			os.Exit(65)
		}
		return &This{Keyword: p.previous()}
	case p.match("SUPER"):
		keyword := p.previous()
		p.consume("DOT", "Expect '.' after 'super'.")
		p.consume("IDENTIFIER", "Expect superclass method name.")
		method := p.previous()
		return &Super{Keyword: keyword, Method: method}
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

func (p *Parser) customError(msg, lexem string, line int) {
	fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", line, lexem, msg)
}
