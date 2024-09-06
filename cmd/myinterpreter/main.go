// main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <command> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	rawFileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	logEnabled := false
	if command == "tokenize" {
		logEnabled = true
	}

	switch command {
	case "tokenize":
		scanner := NewLexer(string(rawFileContent), logEnabled)
		scanner.ScanTokens()
	case "parse":
		scanner := NewLexer(string(rawFileContent), logEnabled)
		scanner.ScanTokens() // Tokenize first
		parser := NewParser(scanner)
		ast := parser.Parse()
		fmt.Println(ast.String()) // Print the AST
	case "evaluate":
		scanner := NewLexer(string(rawFileContent), logEnabled)
		scanner.ScanTokens() // Tokenize first
		parser := NewParser(scanner)
		ast := parser.Parse()
		result := ast.Eval() // Evaluate the AST
		// fmt.Println(ast.String())
		fmt.Println(result)  // Print the evaluated result
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// DebugStruct function to print the values of Lexer and Parser
// func DebugStruct(lexer *Lexer, parser *Parser) {
// 	fmt.Println("===== Lexer State =====")
// 	fmt.Printf("Source Code: %s\n", lexer.source)
// 	fmt.Printf("Current Character: %c\n", lexer.ch)
// 	fmt.Printf("Current Position: %d\n", lexer.position)
// 	fmt.Printf("Next Position: %d\n", lexer.nextPosition)
// 	fmt.Printf("Line: %d\n", lexer.line)
// 	fmt.Println("Tokens:")
// 	for _, token := range lexer.tokens {
// 		fmt.Printf("Type: %s, Lexeme: %s, Literal: %s, Line: %d\n",
// 			token.Type, token.Lexeme, token.Literal, token.Line)
// 	}

	
// }