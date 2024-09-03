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
	EQUAL = '='
	BANG = '!'
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
	for i := 0; i < len(s.source); i++ {
		switch content := s.source[i]; content{
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
		case EQUAL:
			if i+1 < len(s.source) && s.source[i+1] == EQUAL {
				fmt.Println("EQUAL_EQUAL == null")
				i++ // skip the next character as it's part of ==
			} else {
				fmt.Println("EQUAL = null")
			}
		case BANG:
			if i+1 < len(s.source) && s.source[i+1] == EQUAL {
				fmt.Println("BANG_EQUAL != null")
				i++ // skip the next character as it's part of ==
			} else {
				fmt.Println("BANG ! null")
			}
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

func (s *Scanner) reportError(content byte) {
	errorMessage := fmt.Sprintf("[line %d] Error: Unexpected character: %c", s.line, content)
	fmt.Fprintln(os.Stderr, errorMessage)
	s.errors = append(s.errors, errorMessage)
}

func isWhitespace(r byte) bool {
	return r == ' ' || r == '\r' || r == '\t' || r == '\n'
}
