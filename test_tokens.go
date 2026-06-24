package main

import (
	"fmt"
	"esql-ast-tool/pkg/parser"
)

func main() {
	input := `CREATE COMPUTE MODULE TestModule
    DECLARE myVar INTEGER;
    SET myVar = 42;
    IF myVar > 0 THEN
        SET Environment.Variables.Status = 'OK';
    END IF;
END MODULE;`
	l := parser.NewLexer(input)
	tokens := l.Tokenize()
	for _, t := range tokens {
		fmt.Printf("%s %s line:%d col:%d\n", t.Type, t.Literal, t.Line, t.Column)
	}
}
