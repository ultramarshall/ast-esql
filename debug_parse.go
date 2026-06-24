// debug_parse.go
package main

import (
	"esql-ast-tool/pkg/parser"
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("examples/test_cast.esql")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	input := string(data)
	p := parser.NewParser(input)

	// Print semua tokens dulu
	fmt.Println("=== TOKENS ===")
	tokens := p.GetTokens()
	for i, tok := range tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}

	fmt.Println("\n=== START PARSING ===")

	// Override parser dengan debug
	program := parseProgramWithDebug(p)

	if len(p.Errors()) > 0 {
		fmt.Println("\n=== ERRORS ===")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err)
		}
	}

	fmt.Printf("\n=== STATEMENTS: %d ===\n", len(program.Statements))
}

func parseProgramWithDebug(p *parser.Parser) parser.Program {
	program := parser.NewProgram()
	maxLoops := 1000
	loopCount := 0

	for p.CurToken().Type != "EOF" {
		loopCount++
		if loopCount > maxLoops {
			fmt.Printf("ERROR: Infinite loop detected at token: %s, literal: '%s', line: %d\n",
				p.CurToken().Type, p.CurToken().Literal, p.CurToken().Line)
			break
		}

		// Debug setiap iterasi
		fmt.Printf("[Loop %d] Token: %s, Literal: '%s', Line: %d\n",
			loopCount, p.CurToken().Type, p.CurToken().Literal, p.CurToken().Line)

		// Skip END tokens
		if p.CurToken().Type == "END" {
			nextToken := p.PeekToken()
			if nextToken.Type == "IF" {
				p.NextToken()
				p.NextToken()
				if p.CurToken().Type == ";" {
					p.NextToken()
				}
				continue
			}
			if nextToken.Type == "MODULE" {
				p.NextToken()
				p.NextToken()
				if p.CurToken().Type == ";" {
					p.NextToken()
				}
				continue
			}
			p.NextToken()
			continue
		}

		stmt := p.ParseStatement()
		if stmt.Type != "" {
			program.AddStatement(stmt)
			fmt.Printf("  Added statement: %s\n", stmt.Type)
		}

		// Konsumsi semicolon jika ada
		if p.CurToken().Type == ";" {
			p.NextToken()
		} else if p.CurToken().Type != "EOF" && p.CurToken().Type != "END" {
			// Safety: jangan advance
			fmt.Printf("  WARNING: Unexpected token: %s, literal: '%s'\n",
				p.CurToken().Type, p.CurToken().Literal)
			p.NextToken()
		}
	}

	return program
}
