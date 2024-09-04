package main

import (
	"fmt"
	"os"
	"unicode"
)

// Token structure to represent each token in the source code
type Token struct {
	Type    string
	Lexeme  string
	Literal string
	Line    int
}

// Lexer structure to maintain the state of lexical analysis
type Lexer struct {
	source      string
	line        int
	errors      []string
	position    int
	nextPosition int
	ch          byte
	tokens      []Token
}

// NewLexer initializes a new lexer for the given source code
func NewLexer(source string) *Lexer {
	l := &Lexer{
		source: source,
		line:   1,
		errors: []string{},
		tokens: []Token{},
	}
	l.readChar()
	return l
}

// readChar reads the next character and updates position in the source
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.source) {
		l.ch = 0 // End of file
	} else {
		l.ch = l.source[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

// ScanTokens processes the source and generates tokens
func (l *Lexer) ScanTokens() {
	for l.ch != 0 {
		switch {
		case l.isDigit():
			l.handleNumberLiteral()
		case l.ch == '(':
			l.addToken("LEFT_PAREN", "(", "")
			fmt.Println("LEFT_PAREN ( null")
		case l.ch == ')':
			l.addToken("RIGHT_PAREN", ")", "")
			fmt.Println("RIGHT_PAREN ) null")
		case l.ch == '{':
			l.addToken("LEFT_BRACE", "{", "")
			fmt.Println("LEFT_BRACE { null")
		case l.ch == '}':
			l.addToken("RIGHT_BRACE", "}", "")
			fmt.Println("RIGHT_BRACE } null")
		case l.ch == '*':
			l.addToken("STAR", "*", "")
			fmt.Println("STAR * null")
		case l.ch == '.':
			l.addToken("DOT", ".", "")
			fmt.Println("DOT . null")
		case l.ch == ',':
			l.addToken("COMMA", ",", "")
			fmt.Println("COMMA , null")
		case l.ch == '+':
			l.addToken("PLUS", "+", "")
			fmt.Println("PLUS + null")
		case l.ch == '-':
			l.addToken("MINUS", "-", "")
			fmt.Println("MINUS - null")
		case l.ch == ';':
			l.addToken("SEMICOLON", ";", "")
			fmt.Println("SEMICOLON ; null")
		case l.ch == '=':
			if l.peekChar() == '=' {
				l.addToken("EQUAL_EQUAL", "==", "")
				fmt.Println("EQUAL_EQUAL == null")
				l.readChar()
			} else {
				l.addToken("EQUAL", "=", "")
				fmt.Println("EQUAL = null")
			}
		case l.ch == '!':
			if l.peekChar() == '=' {
				l.addToken("BANG_EQUAL", "!=", "")
				fmt.Println("BANG_EQUAL != null")
				l.readChar()
			} else {
				l.addToken("BANG", "!", "")
				fmt.Println("BANG ! null")
			}
		case l.ch == '<':
			if l.peekChar() == '=' {
				l.addToken("LESS_EQUAL", "<=", "")
				fmt.Println("LESS_EQUAL <= null")
				l.readChar()
			} else {
				l.addToken("LESS", "<", "")
				fmt.Println("LESS < null")
			}
		case l.ch == '>':
			if l.peekChar() == '=' {
				l.addToken("GREATER_EQUAL", ">=", "")
				fmt.Println("GREATER_EQUAL >= null")
				l.readChar()
			} else {
				l.addToken("GREATER", ">", "")
				fmt.Println("GREATER > null")
			}
		case l.ch == '/':
			if l.peekChar() == '/' {
				l.skipComment()
			} else {
				l.addToken("SLASH", "/", "")
				fmt.Println("SLASH / null")
			}
		case l.ch == '"':
			l.handleStringLiteral()
		default:
			if l.isAlpha() {
				l.handleIdentifier()
				continue
			}
			if l.isWhitespace() {
				if l.ch == '\n' {
					l.line++
				}
				l.readChar()
				continue
			}
			l.reportError(l.ch)
		}
		l.readChar()
	}
	fmt.Println("EOF  null")

	if len(l.errors) > 0 {
		os.Exit(65)
	}
}

// addToken creates a new token and appends it to the token list
func (l *Lexer) addToken(tokenType, lexeme, literal string) {
	token := Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    l.line,
	}
	l.tokens = append(l.tokens, token)
}

// handleIdentifier processes identifiers or reserved keywords
func (l *Lexer) handleIdentifier() {
	startPosition := l.position
	for l.isAlpha() || l.isDigit() {
		l.readChar()
	}

	identifier := l.source[startPosition:l.position]
	keyword, ok := RESERVED_WORDS[identifier]
	if !ok {
		l.addToken("IDENTIFIER", identifier, "")
		fmt.Printf("IDENTIFIER %s null\n", identifier)
	} else {
		l.addToken(keyword, identifier, "")
		fmt.Printf("%s %s null\n", keyword, identifier)
	}
}

// handleStringLiteral processes string literals
func (l *Lexer) handleStringLiteral() {
	startPosition := l.position
	for {
		l.readChar()
		if l.ch == '"' {
			literal := l.source[startPosition+1 : l.position]
			l.addToken("STRING", literal, literal)
			fmt.Printf("STRING \"%s\" %s\n", literal, literal)
			return
		} else if l.ch == 0 || l.ch == '\n' {
			l.reportErrorUnterminatedString()
			return
		}
	}
}

// handleNumberLiteral processes numeric literals
func (l *Lexer) handleNumberLiteral() {
	startPosition := l.position
	for l.isDigit() {
		l.readChar()
	}
	if l.ch == '.' && l.isDigitAtNextPosition() {
		l.readChar()
		for l.isDigit() {
			l.readChar()
		}
	}

	literal := l.source[startPosition:l.position]
	l.addToken("NUMBER", literal, formatAsFloat(literal))
	fmt.Printf("NUMBER %s %s\n", literal, formatAsFloat(literal))
	l.position--
	l.nextPosition--
}

// skipComment skips comments starting with '//'
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	if l.ch == '\n' {
		l.line++
	}
}

// Helper functions
func (l *Lexer) isDigit() bool {
	return l.ch >= '0' && l.ch <= '9'
}

func (l *Lexer) isAlpha() bool {
	return unicode.IsLetter(rune(l.ch)) || l.ch == '_'
}

func (l *Lexer) isWhitespace() bool {
	return l.ch == ' ' || l.ch == '\r' || l.ch == '\t' || l.ch == '\n'
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.source) {
		return 0
	}
	return l.source[l.nextPosition]
}

func (l *Lexer) isDigitAtNextPosition() bool {
	if l.nextPosition >= len(l.source) {
		return false
	}
	return unicode.IsDigit(rune(l.source[l.nextPosition]))
}

// Error handling
func (l *Lexer) reportErrorUnterminatedString() {
	errorMessage := fmt.Sprintf("[line %d] Error: Unterminated string.", l.line)
	fmt.Fprintln(os.Stderr, errorMessage)
	l.errors = append(l.errors, errorMessage)
}

func (s *Lexer) reportError(content byte) {
	errorMessage := fmt.Sprintf("[line %d] Error: Unexpected character: %c", s.line, content)
	fmt.Fprintln(os.Stderr, errorMessage)
	s.errors = append(s.errors, errorMessage)
}

// Helper function to format a number literal as a floating-point number
func formatAsFloat(literal string) string {
	if !containsDecimalPoint(literal) {
		return literal + ".0" // Add ".0" for integers
	}
	// If there's a decimal point, trim unnecessary trailing zeros but keep at least one digit after the decimal
	return trimTrailingZeros(literal)
}

// Helper function to check if a literal contains a decimal point
func containsDecimalPoint(literal string) bool {
	for _, ch := range literal {
		if ch == '.' {
			return true
		}
	}
	return false
}

// Helper function to trim trailing zeros from the fractional part
func trimTrailingZeros(literal string) string {
	decimalPos := -1
	for i := range literal {
		if literal[i] == '.' {
			decimalPos = i
			break
		}
	}

	if decimalPos == -1 {
		return literal // No decimal point, return the literal as is
	}

	// Start from the end of the string and remove trailing zeros
	endPos := len(literal)
	for endPos > decimalPos+1 && literal[endPos-1] == '0' {
		endPos--
	}

	// Ensure there's at least one digit after the decimal
	if endPos == decimalPos+1 {
		return literal[:endPos] + "0"
	}

	return literal[:endPos]
}
