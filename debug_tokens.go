// debug_tokens.go
package main

import (
	"esql-ast-tool/pkg/parser"
	"fmt"
)

func main() {
	input := `CAST('123' AS INTEGER)`
	l := parser.NewLexer(input)
	tokens := l.Tokenize()

	fmt.Println("=== TOKENS ===")
	for i, tok := range tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}
