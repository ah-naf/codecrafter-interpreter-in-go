package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN  = '('
	RIGHT_PAREN = ')'
	LEFT_BRACE = '{'
	RIGHT_BRACE = '}'
	STAR = '*'
	DOT = '.'
	COMMA = ','
	PLUS = '+'
	MINUS = '-'
	SEMICOLON = ';'
)

type Scanner struct {
	source string
	line   int
	errors []string
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
		errors: []string{},
	}
}

func (s *Scanner) ScanTokens() {
	for _, content := range s.source {
		switch content {
		case LEFT_PAREN:
			fmt.Println("LEFT_PAREN ( null")
		case RIGHT_PAREN:
			fmt.Println("RIGHT_PAREN ) null")
		case LEFT_BRACE:
			fmt.Println("LEFT_BRACE { null")
		case RIGHT_BRACE:
			fmt.Println("RIGHT_BRACE } null")
		case STAR:
			fmt.Println("STAR * null")
		case DOT:
			fmt.Println("DOT . null")
		case COMMA:
			fmt.Println("COMMA , null")
		case PLUS:
			fmt.Println("PLUS + null")
		case MINUS:
			fmt.Println("MINUS - null")
		case SEMICOLON:
			fmt.Println("SEMICOLON ; null")
		default:
			if isWhitespace(content) {
				if content == '\n' {
					s.line++
				}
				continue
			}
			s.reportError(content)
		}
	}
	fmt.Println("EOF  null")

	if len(s.errors) > 0 {
		os.Exit(65)
	}
}

func (s *Scanner) reportError(content rune) {
	errorMessage := fmt.Sprintf("[line %d] Error: Unexpected character: %c", s.line, content)
	fmt.Fprintln(os.Stderr, errorMessage)
	s.errors = append(s.errors, errorMessage)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\r' || r == '\t' || r == '\n'
}
