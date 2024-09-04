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

	switch command {
	case "tokenize":
		scanner := NewLexer(string(rawFileContent))
		scanner.ScanTokens()
	case "parse":
		scanner := NewLexer(string(rawFileContent))
		parser := NewParser(scanner)
		ast := parser.Parse()
		fmt.Println(ast.String())
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

