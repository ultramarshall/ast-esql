package main

import (
	"esql-ast-tool/pkg/parser"
	"esql-ast-tool/pkg/printer"
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("examples/sample.esql")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	input := string(data)
	p := parser.NewParser(input)

	fmt.Println("=== PARSING WITH DEBUG ===")
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("\n=== ERRORS ===")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err)
		}
	}

	fmt.Println("\n=== AST ===")
	pr := printer.NewPrinter()
	fmt.Println(pr.PrintProgram(program))
}
