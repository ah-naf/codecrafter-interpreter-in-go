package main

const (
	LEFT_PAREN  = '('
	RIGHT_PAREN = ')'
	LEFT_BRACE  = '{'
	RIGHT_BRACE = '}'
	STAR        = '*'
	DOT         = '.'
	COMMA       = ','
	PLUS        = '+'
	MINUS       = '-'
	SEMICOLON   = ';'
	EQUAL       = '='
	BANG        = '!'
	LT          = '<'
	GT          = '>'
	SLASH       = '/'
)

var RESERVED_WORDS = map[string]string {
	"and": "AND",
	"class": "CLASS",
	"else": "ELSE",
	"false": "FALSE",
    "for": "FOR",
	"fun": "FUN",
    "if": "IF",
    "nil": "NIL",
    "or": "OR",
	"print": "PRINT",
	"return": "RETURN",
	"super": "SUPER",
	"this": "THIS",
    "true": "TRUE",
    "var": "VAR",
	"while": "WHILE",
}