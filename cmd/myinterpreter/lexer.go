package main

import (
	"fmt"
	"os"
)

type Lexer struct {
	source string
	line   int
	errors []string
	position int
	nextPosition int
	ch byte
}

func NewLexer(source string) *Lexer {
	l := &Lexer{
		source: source,
		line:   1,
		errors: []string{},
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.source) {
        l.ch = 0
    } else {
        l.ch = l.source[l.nextPosition]
    }
    l.position = l.nextPosition
    l.nextPosition++
}

func (l *Lexer) ScanTokens() {
	for l.ch != 0 {
		switch l.ch {
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
			if l.peekChar() == EQUAL {
				fmt.Println("EQUAL_EQUAL == null")
				l.readChar()
			} else {
				fmt.Println("EQUAL = null")
			}
		case BANG:
			if l.peekChar() == EQUAL {
				fmt.Println("BANG_EQUAL != null")
				l.readChar()
			} else {
				fmt.Println("BANG ! null")
			}
		case LT:
			if l.peekChar() == EQUAL {
				fmt.Println("LESS_EQUAL <= null")
				l.readChar()
			} else {
				fmt.Println("LESS < null")
			}
		case GT:
			if l.peekChar() == EQUAL {
				fmt.Println("GREATER_EQUAL >= null")
				l.readChar()
			} else {
				fmt.Println("GREATER > null")
			}
		case SLASH:
			if l.peekChar() == SLASH {
				// It's a comment, skip until end of line
				l.skipComment()
			} else {
				fmt.Println("SLASH / null")
			}
		case '"':
			l.handleStringLiteral()
		default:
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

func (l *Lexer) handleStringLiteral() {
	startPosition := l.position
	for {
		l.readChar()
		if l.ch == '"' {
			// Found the closing quote
			literal := l.source[startPosition+1 : l.position]
			fmt.Printf("STRING \"%s\" %s\n", literal, literal)
			return
		} else if l.ch == 0 {
			// End of file reached, unterminated string
			l.reportErrorUnterminatedString()
			return
		} else if l.ch == '\n' {
			// Unterminated string on the current line
			l.reportErrorUnterminatedString()
			l.line++
			return
		}
	}
}

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

func (l *Lexer) isWhitespace() bool {
	return l.ch == ' ' || l.ch == '\r' || l.ch == '\t' || l.ch == '\n'
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	
	if l.ch == '\n' {
		l.line++
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.source) {
		return 0
	} else {
		return l.source[l.nextPosition]
	}
}