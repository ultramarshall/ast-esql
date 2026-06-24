package main

import (
	"fmt"
	"esql-ast-tool/pkg/parser"
)

func main() {
	input := `CASE WHEN score > 80 THEN 1 ELSE 0 END`
	l := parser.NewLexer(input)
	tokens := l.Tokenize()
	
	fmt.Println("=== TOKENS ===")
	for i, tok := range tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}
