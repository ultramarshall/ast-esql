# Dokumentasi Struktur & Kode Project

Dokumen ini dihasilkan secara otomatis untuk memetakan seluruh struktur folder dan isi kode di dalam project ini.

## Struktur Project (Tree)

```text
.
├── cmd
│   └── esql-ast
│       └── main.go
├── debug.go
├── DOC.md
├── esql-ast
├── examples
│   └── sample.esql
├── generate_doc.sh
├── go.mod
├── internal
│   └── token
│       └── token.go
├── Makefile
├── pkg
│   ├── analyzer
│   │   └── analyzer.go
│   ├── generator
│   │   └── generator.go
│   ├── parser
│   │   ├── ast.go
│   │   ├── lexer.go
│   │   └── parser.go
│   └── printer
│       └── printer.go
└── test_tokens.go

11 directories, 16 files
```

## Isi Kode Berdasarkan File

### File: `go.mod`

```text
module esql-ast-tool

go 1.22

```

---

### File: `examples/sample.esql`

```text
CREATE COMPUTE MODULE TestModule
    DECLARE myVar INTEGER;
    SET myVar = 42;
    IF myVar > 0 THEN
        SET Environment.Variables.Status = 'OK';
    END IF;
END MODULE;

```

---

### File: `debug.go`

```go
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

```

---

### File: `pkg/parser/ast.go`

```go
package parser

import (
	"encoding/json"
)

type NodeType string

const (
	// Statements
	ProgramNode     NodeType = "Program"
	ModuleNode      NodeType = "Module"
	FunctionNode    NodeType = "Function"
	ProcedureNode   NodeType = "Procedure"
	DeclareNode     NodeType = "Declare"
	SetNode         NodeType = "Set"
	IfNode          NodeType = "If"
	ElseNode        NodeType = "Else"
	ElseIfNode      NodeType = "ElseIf"
	WhileNode       NodeType = "While"
	ForNode         NodeType = "For"
	CaseNode        NodeType = "Case"
	WhenNode        NodeType = "When"
	OtherwiseNode   NodeType = "Otherwise"
	ReturnNode      NodeType = "Return"
	ThrowNode       NodeType = "Throw"
	CreateNode      NodeType = "Create"
	EnvironmentNode NodeType = "Environment"
	DatabaseNode    NodeType = "Database"
	PassthruNode    NodeType = "Passthru"
	MoveNode        NodeType = "Move"
	PropagateNode   NodeType = "Propagate"
	ContinueNode    NodeType = "Continue"
	BreakNode       NodeType = "Break"
	LabelNode       NodeType = "Label"
	BlockNode       NodeType = "Block"
	CallNode        NodeType = "Call"

	// Expressions
	BinaryOpNode       NodeType = "BinaryOp"
	UnaryOpNode        NodeType = "UnaryOp"
	ComparisonNode     NodeType = "Comparison"
	FunctionCallNode   NodeType = "FunctionCall"
	FieldReferenceNode NodeType = "FieldReference"
	ArrayIndexNode     NodeType = "ArrayIndex"
	LiteralNode        NodeType = "Literal"
	IdentifierNode     NodeType = "Identifier"
)

type ASTNode struct {
	Type     NodeType    `json:"type"`
	Value    interface{} `json:"value,omitempty"`
	Children []ASTNode   `json:"children,omitempty"`
	Line     int         `json:"line"`
	Column   int         `json:"column"`
	Token    string      `json:"-"`
}

func NewASTNode(nodeType NodeType, token string, line, column int) ASTNode {
	return ASTNode{
		Type:     nodeType,
		Children: []ASTNode{},
		Line:     line,
		Column:   column,
		Token:    token,
	}
}

func (n *ASTNode) AddChild(child ASTNode) {
	n.Children = append(n.Children, child)
}

func (n ASTNode) ToJSON() ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}

type Program struct {
	Statements []ASTNode `json:"statements"`
}

func NewProgram() Program {
	return Program{
		Statements: []ASTNode{},
	}
}

func (p *Program) AddStatement(stmt ASTNode) {
	p.Statements = append(p.Statements, stmt)
}

func (p Program) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

```

---

### File: `pkg/parser/parser.go`

```go
package parser

import (
	"fmt"
	"strconv"

	"esql-ast-tool/internal/token"
)

type Parser struct {
	l         *Lexer
	tokens    []token.Token
	position  int
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(input string) *Parser {
	l := NewLexer(input)
	tokens := l.Tokenize()
	p := &Parser{
		l:        l,
		tokens:   tokens,
		position: 0,
		errors:   []string{},
	}

	if len(tokens) > 0 {
		p.curToken = tokens[0]
		if len(tokens) > 1 {
			p.peekToken = tokens[1]
		}
	}

	return p
}

func (p *Parser) nextToken() {
	p.position++
	if p.position < len(p.tokens) {
		p.curToken = p.tokens[p.position]
		if p.position+1 < len(p.tokens) {
			p.peekToken = p.tokens[p.position+1]
		} else {
			p.peekToken = token.Token{Type: token.EOF}
		}
	} else {
		p.curToken = token.Token{Type: token.EOF}
		p.peekToken = token.Token{Type: token.EOF}
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) DebugTokens() {
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// GetTokens returns all tokens for debugging
func (p *Parser) GetTokens() []token.Token {
	return p.tokens
}

// DebugPrintTokens prints all tokens
func (p *Parser) DebugPrintTokens() {
	fmt.Println("=== TOKENS ===")
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// Parse program
func (p *Parser) ParseProgram() Program {
	program := NewProgram()

	for p.curToken.Type != token.EOF {
		// Debug
		// fmt.Printf("DEBUG ParseProgram: token=%s, literal='%s'\n", p.curToken.Type, p.curToken.Literal)

		// Skip END tokens - tapi hati-hati jangan skip token yang penting
		if p.curToken.Type == token.END {
			// Cek apakah ini END IF atau END MODULE
			nextToken := p.peekToken
			if nextToken.Type == token.IF {
				p.nextToken() // consume END
				p.nextToken() // consume IF
				if p.curToken.Type == token.SEMICOLON {
					p.nextToken()
				}
				continue
			}
			if nextToken.Type == token.MODULE {
				p.nextToken() // consume END
				p.nextToken() // consume MODULE
				if p.curToken.Type == token.SEMICOLON {
					p.nextToken()
				}
				continue
			}
			// END lainnya, skip
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt.Type != "" {
			program.AddStatement(stmt)
		}

		// Konsumsi semicolon jika ada
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.EOF && p.curToken.Type != token.END {
			// Safety advance - tapi hati-hati
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseStatement() ASTNode {
	// Debug - matikan setelah selesai debugging
	// fmt.Printf("DEBUG parseStatement: token=%s, literal='%s', line=%d\n",
	//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	switch p.curToken.Type {
	case token.CREATE:
		return p.parseCreate()
	case token.DECLARE:
		return p.parseDeclare()
	case token.SET:
		return p.parseSet()
	case token.IF:
		return p.parseIf()
	case token.WHILE:
		return p.parseWhile()
	case token.FOR:
		return p.parseFor()
	case token.RETURN:
		return p.parseReturn()
	case token.THROW:
		return p.parseThrow()
	case token.PROPAGATE:
		return p.parsePropagate()
	case token.MOVE:
		return p.parseMove()
	case token.CONTINUE:
		return p.parseContinue()
	case token.BREAK:
		return p.parseBreak()
	case token.LABEL:
		return p.parseLabel()
	case token.MODULE:
		return p.parseModuleStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.PROCEDURE:
		return p.parseProcedureStatement()
	case token.CALL:
		return p.parseCall()
	case token.END:
		// Skip END tokens
		p.nextToken()
		return ASTNode{}
	default:
		// Parse sebagai expression
		expr := p.parseExpression()
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return expr
	}
}

// Parse CALL statement
func (p *Parser) parseCall() ASTNode {
	node := NewASTNode(CallNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse procedure name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse arguments
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			arg := p.parseExpression()
			if arg.Type != "" {
				node.AddChild(arg)
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// Parse MODULE statement
func (p *Parser) parseModuleStatement() ASTNode {
	node := NewASTNode(ModuleNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse module name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse body statements - jangan skip apapun
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		// Debug
		// fmt.Printf("DEBUG parseModuleStatement: token=%s, literal='%s', line=%d\n",
		//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}

		// Jangan panggil p.nextToken() di sini karena parseStatement sudah handle
		// Hanya konsumsi semicolon jika ada dan parser belum lanjut
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			// Hanya advance jika benar-benar perlu
			// Tapi hati-hati, ini bisa menyebabkan skip token
			// p.nextToken()
		}
	}

	// Konsumsi END MODULE
	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
		}
	}

	// Optional semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// Parse FUNCTION statement
func (p *Parser) parseFunctionStatement() ASTNode {
	node := NewASTNode(FunctionNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse function name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		if nameNode.Value == nil {
			nameNode.Value = nameNode.Token
		}
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse parameters
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			if p.curToken.Type == token.IDENTIFIER {
				paramNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				if paramNode.Value == nil {
					paramNode.Value = paramNode.Token
				}
				node.AddChild(paramNode)
				p.nextToken()
				// Skip type
				if p.curToken.Type == token.IDENTIFIER {
					p.nextToken()
				}
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	// Parse return type
	if p.curToken.Type == token.RETURNS {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			returnType := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			if returnType.Value == nil {
				returnType.Value = returnType.Token
			}
			node.AddChild(returnType)
			p.nextToken()
		}
	}

	// Parse function body
	if p.curToken.Type == token.BEGIN {
		p.nextToken()
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				node.AddChild(stmt)
			}
			if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.END {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// Parse PROCEDURE statement
func (p *Parser) parseProcedureStatement() ASTNode {
	node := NewASTNode(ProcedureNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse procedure name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		if nameNode.Value == nil {
			nameNode.Value = nameNode.Token
		}
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse parameters
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			if p.curToken.Type == token.IDENTIFIER {
				paramNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				if paramNode.Value == nil {
					paramNode.Value = paramNode.Token
				}
				node.AddChild(paramNode)
				p.nextToken()
				// Skip type
				if p.curToken.Type == token.IDENTIFIER {
					p.nextToken()
				}
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	// Parse procedure body
	if p.curToken.Type == token.BEGIN {
		p.nextToken()
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				node.AddChild(stmt)
			}
			if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.END {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseDeclare() ASTNode {
	node := NewASTNode(DeclareNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse variable name
	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors, "expected identifier after DECLARE")
		return node
	}

	nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	nameNode.Value = p.curToken.Literal
	node.AddChild(nameNode)
	p.nextToken()

	// Parse type (INTEGER, STRING, etc.)
	if p.curToken.Type == token.IDENTIFIER {
		typeNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		typeNode.Value = p.curToken.Literal
		node.AddChild(typeNode)
		p.nextToken()
	}

	// Optional DEFAULT
	if p.curToken.Type == token.DEFAULT {
		p.nextToken()
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseSet() ASTNode {
	// Debug
	// fmt.Printf("DEBUG parseSet: token=%s, literal='%s', line=%d\n",
	//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(SetNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse target - ini bisa identifier atau field reference
	target := p.parseExpression()
	if target.Type != "" {
		targetWrapper := NewASTNode(BlockNode, "target", target.Line, target.Column)
		targetWrapper.AddChild(target)
		node.AddChild(targetWrapper)
	}

	// Debug
	// fmt.Printf("DEBUG parseSet: after target, token=%s, literal='%s'\n",
	//     p.curToken.Type, p.curToken.Literal)

	// Cek ASSIGN
	if p.curToken.Type == token.ASSIGN {
		// fmt.Println("DEBUG parseSet: found ASSIGN")
		p.nextToken()
		value := p.parseExpression()
		if value.Type != "" {
			valueWrapper := NewASTNode(BlockNode, "value", value.Line, value.Column)
			valueWrapper.AddChild(value)
			node.AddChild(valueWrapper)
		}
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '=' in SET statement, got %s '%s' at line %d",
				p.curToken.Type, p.curToken.Literal, p.curToken.Line))
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseIf() ASTNode {
	node := NewASTNode(IfNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse condition
	cond := p.parseExpression()
	if cond.Type != "" {
		condWrapper := NewASTNode(BlockNode, "condition", cond.Line, cond.Column)
		condWrapper.AddChild(cond)
		node.AddChild(condWrapper)
	}

	// Expect THEN
	if p.curToken.Type == token.THEN {
		p.nextToken()
	}

	// Parse THEN block
	thenBlock := NewASTNode(BlockNode, "then", p.curToken.Line, p.curToken.Column)

	for p.curToken.Type != token.END && p.curToken.Type != token.ELSE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			thenBlock.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}
	node.AddChild(thenBlock)

	// Parse ELSE if ada
	if p.curToken.Type == token.ELSE {
		p.nextToken()
		elseBlock := NewASTNode(BlockNode, "else", p.curToken.Line, p.curToken.Column)
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				elseBlock.AddChild(stmt)
			}
			if p.curToken.Type == token.SEMICOLON {
				p.nextToken()
			}
		}
		node.AddChild(elseBlock)
	}

	// Konsumsi END IF
	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.IF {
			p.nextToken()
		}
	}

	// Optional semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseWhile() ASTNode {
	node := NewASTNode(WhileNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	condition := p.parseExpression()
	if condition.Type != "" {
		node.AddChild(condition)
	}

	if p.curToken.Type == token.DO {
		p.nextToken()
	}

	bodyNode := NewASTNode(BlockNode, "body", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			bodyNode.AddChild(stmt)
		}
		if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			p.nextToken()
		}
	}
	node.AddChild(bodyNode)

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.WHILE {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFor() ASTNode {
	node := NewASTNode(ForNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER {
		varNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		varNode.Value = varNode.Token
		node.AddChild(varNode)
		p.nextToken()
	}

	if p.curToken.Type == token.AS {
		p.nextToken()
	}

	for p.curToken.Type != token.DO && p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.DO {
		p.nextToken()
	}

	bodyNode := NewASTNode(BlockNode, "body", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			bodyNode.AddChild(stmt)
		}
		if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			p.nextToken()
		}
	}
	node.AddChild(bodyNode)

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.FOR {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseReturn() ASTNode {
	node := NewASTNode(ReturnNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseThrow() ASTNode {
	node := NewASTNode(ThrowNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parsePropagate() ASTNode {
	node := NewASTNode(PropagateNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseBreak() ASTNode {
	node := NewASTNode(BreakNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseMove() ASTNode {
	node := NewASTNode(MoveNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	target := p.parseExpression()
	if target.Type != "" {
		node.AddChild(target)
	}

	if p.curToken.Type == token.TO {
		p.nextToken()
		source := p.parseExpression()
		if source.Type != "" {
			node.AddChild(source)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseContinue() ASTNode {
	node := NewASTNode(ContinueNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	if p.curToken.Type == token.IDENTIFIER { // fix for continue label perhaps
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseLabel() ASTNode {
	node := NewASTNode(LabelNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER {
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseCreate() ASTNode {
	node := NewASTNode(CreateNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.COMPUTE {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()

			// Create Module node
			moduleNode := NewASTNode(ModuleNode, "COMPUTE MODULE", p.curToken.Line, p.curToken.Column)
			moduleNode.Value = "COMPUTE"

			// Parse module name
			if p.curToken.Type == token.IDENTIFIER {
				nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				nameNode.Value = p.curToken.Literal
				moduleNode.AddChild(nameNode)
				p.nextToken()
			}

			// Parse body statements
			for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
				// Debug
				// fmt.Printf("DEBUG parseCreate: token=%s, literal='%s', line=%d\n",
				//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

				stmt := p.parseStatement()
				if stmt.Type != "" {
					moduleNode.AddChild(stmt)
				}

				// Hanya konsumsi semicolon jika perlu
				if p.curToken.Type == token.SEMICOLON {
					p.nextToken()
				}
			}

			// Konsumsi END MODULE
			if p.curToken.Type == token.END {
				p.nextToken()
				if p.curToken.Type == token.MODULE {
					p.nextToken()
				}
			}

			node.AddChild(moduleNode)
		}
	}

	// Optional semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseExpressionStatement() ASTNode {
	expr := p.parseExpression()
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return expr
}

func (p *Parser) parseExpression() ASTNode {
	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() ASTNode {
	node := p.parseLogicalAnd()

	for p.curToken.Type == token.OR {
		tok := p.curToken
		p.nextToken()
		right := p.parseLogicalAnd()
		if right.Type != "" {
			binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
			binOp.AddChild(node)
			binOp.AddChild(right)
			node = binOp
		}
	}

	return node
}

func (p *Parser) parseLogicalAnd() ASTNode {
	node := p.parseComparison()

	for p.curToken.Type == token.AND {
		tok := p.curToken
		p.nextToken()
		right := p.parseComparison()
		if right.Type != "" {
			binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
			binOp.AddChild(node)
			binOp.AddChild(right)
			node = binOp
		}
	}

	return node
}

func (p *Parser) parseComparison() ASTNode {
	node := p.parseAdditive()

	if p.curToken.Type == token.EQ || p.curToken.Type == token.NOT_EQ ||
		p.curToken.Type == token.LT || p.curToken.Type == token.GT ||
		p.curToken.Type == token.LTE || p.curToken.Type == token.GTE {
		tok := p.curToken
		p.nextToken()
		right := p.parseAdditive()
		if right.Type != "" {
			compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
			compNode.AddChild(node)
			compNode.AddChild(right)
			return compNode
		}
	}

	return node
}

func (p *Parser) parseAdditive() ASTNode {
	node := p.parseMultiplicative()

	for p.curToken.Type == token.PLUS || p.curToken.Type == token.MINUS {
		tok := p.curToken
		p.nextToken()
		right := p.parseMultiplicative()
		binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
		binOp.AddChild(node)
		binOp.AddChild(right)
		node = binOp
	}

	return node
}

func (p *Parser) parseMultiplicative() ASTNode {
	node := p.parseUnary()

	for p.curToken.Type == token.ASTERISK || p.curToken.Type == token.SLASH ||
		p.curToken.Type == token.MODULO {
		tok := p.curToken
		p.nextToken()
		right := p.parseUnary()
		binOp := NewASTNode(BinaryOpNode, tok.Literal, tok.Line, tok.Column)
		binOp.AddChild(node)
		binOp.AddChild(right)
		node = binOp
	}

	return node
}

func (p *Parser) parseUnary() ASTNode {
	if p.curToken.Type == token.MINUS || p.curToken.Type == token.NOT {
		tok := p.curToken
		p.nextToken()
		operand := p.parsePrimary()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		return unaryNode
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ASTNode {
	switch p.curToken.Type {
	case token.IDENTIFIER:
		return p.parseIdentifier()
	case token.NUMBER:
		return p.parseNumber()
	case token.STRING:
		return p.parseString()
	case token.LPAREN:
		return p.parseGroupedExpression()
	default:
		// Jika token tidak dikenal, jangan buat node error
		// Biarkan parser melanjutkan
		return ASTNode{}
	}
}

func (p *Parser) parseIdentifier() ASTNode {
	node := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()

	// Function call
	if p.curToken.Type == token.LPAREN {
		return p.parseFunctionCall(node)
	}

	// Field reference - Environment.Variables.Status
	if p.curToken.Type == token.DOT {
		return p.parseFieldReference(node)
	}

	return node
}

func (p *Parser) parseFunctionCall(name ASTNode) ASTNode {
	funcName := name.Token
	if name.Value != nil {
		if str, ok := name.Value.(string); ok {
			funcName = str
		}
	}
	node := NewASTNode(FunctionCallNode, funcName, name.Line, name.Column)
	node.Value = funcName
	p.nextToken()

	if p.curToken.Type != token.RPAREN {
		arg := p.parseExpression()
		if arg.Type != "" {
			node.AddChild(arg)
		}

		for p.curToken.Type == token.COMMA {
			p.nextToken()
			arg = p.parseExpression()
			if arg.Type != "" {
				node.AddChild(arg)
			}
		}
	}

	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFieldReference(base ASTNode) ASTNode {
	// Base adalah root (Environment)
	fieldNode := NewASTNode(FieldReferenceNode, "field", base.Line, base.Column)
	fieldNode.AddChild(base)

	// Set initial value
	if base.Value != nil {
		fieldNode.Value = base.Value
	} else {
		fieldNode.Value = base.Token
	}

	// Lanjutkan parsing untuk rest of chain
	for p.curToken.Type == token.DOT {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			identNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			identNode.Value = p.curToken.Literal

			// Buat FieldReference baru yang membungkus yang lama
			newFieldNode := NewASTNode(FieldReferenceNode, "field", fieldNode.Line, fieldNode.Column)
			newFieldNode.AddChild(fieldNode)
			newFieldNode.AddChild(identNode)

			// Update value dengan path lengkap
			if fieldNode.Value != nil {
				newFieldNode.Value = fieldNode.Value.(string) + "." + p.curToken.Literal
			} else {
				newFieldNode.Value = p.curToken.Literal
			}

			fieldNode = newFieldNode
			p.nextToken()
		}
	}

	return fieldNode
}

func (p *Parser) parseNumber() ASTNode {
	val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
	node := NewASTNode(LiteralNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = val
	p.nextToken()
	return node
}

func (p *Parser) parseString() ASTNode {
	node := NewASTNode(LiteralNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()
	return node
}

func (p *Parser) parseGroupedExpression() ASTNode {
	p.nextToken()
	expr := p.parseExpression()
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}
	return expr
}

```

---

### File: `pkg/parser/lexer.go`

```go
package parser

import (
	"strings"

	"esql-ast-tool/internal/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	if l.ch == '-' && l.peekChar() == '-' {
		// Single line comment
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		l.skipWhitespace()
	} else if l.ch == '/' && l.peekChar() == '*' {
		// Multi-line comment
		l.readChar()
		l.readChar()
		for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}
		if l.ch != 0 {
			l.readChar()
			l.readChar()
		}
		l.skipWhitespace()
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	l.skipComments()

	var tok token.Token
	tok.Line = l.line
	tok.Column = l.column
	tok.Pos = l.position

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ASSIGN
			tok.Literal = string(l.ch)
		}
	case '+':
		tok.Type = token.PLUS
		tok.Literal = string(l.ch)
	case '-':
		tok.Type = token.MINUS
		tok.Literal = string(l.ch)
	case '*':
		tok.Type = token.ASTERISK
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = token.SLASH
		tok.Literal = string(l.ch)
	case '%':
		tok.Type = token.MODULO
		tok.Literal = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.NOT_EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.LTE
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.LT
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.GTE
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.GT
			tok.Literal = string(l.ch)
		}
	case ',':
		tok.Type = token.COMMA
		tok.Literal = string(l.ch)
	case ';':
		tok.Type = token.SEMICOLON
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = token.LPAREN
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = token.RPAREN
		tok.Literal = string(l.ch)
	case '{':
		tok.Type = token.LBRACE
		tok.Literal = string(l.ch)
	case '}':
		tok.Type = token.RBRACE
		tok.Literal = string(l.ch)
	case '[':
		tok.Type = token.LBRACKET
		tok.Literal = string(l.ch)
	case ']':
		tok.Type = token.RBRACKET
		tok.Literal = string(l.ch)
	case '.':
		tok.Type = token.DOT
		tok.Literal = string(l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(strings.ToUpper(tok.Literal))
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else if l.ch == '\'' {
			tok.Type = token.STRING
			tok.Literal = l.readString()
			return tok
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	str := l.input[position:l.position]
	l.readChar()
	return str
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) Tokenize() []token.Token {
	var tokens []token.Token
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

```

---

### File: `pkg/generator/generator.go`

```go
package generator

import (
	"fmt"
	"strings"

	"esql-ast-tool/pkg/parser"
)

type Generator struct {
	indent string
}

func NewGenerator() *Generator {
	return &Generator{indent: "    "}
}

func (g *Generator) Generate(program parser.Program) string {
	var sb strings.Builder
	for i, stmt := range program.Statements {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(g.generateNode(stmt, 0))
	}
	return sb.String()
}

func (g *Generator) generateNode(node parser.ASTNode, level int) string {
	if node.Type == "" {
		return ""
	}

	_ = strings.Repeat(g.indent, level) // avoid unused variable

	switch node.Type {
	case parser.CreateNode:
		return g.generateCreate(node, level)

	case parser.ModuleNode:
		return g.generateModule(node, level)

	case parser.DeclareNode:
		return g.generateDeclare(node, level)

	case parser.SetNode:
		return g.generateSet(node, level)

	case parser.IfNode:
		return g.generateIf(node, level)

	case parser.BlockNode:
		return g.generateBlock(node, level)

	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok && name != "" {
			return name
		}
		return node.Token

	case parser.FieldReferenceNode:
		return g.generateFieldReference(node)

	case parser.LiteralNode:
		if str, ok := node.Value.(string); ok {
			return "'" + str + "'"
		}
		if num, ok := node.Value.(float64); ok {
			return fmt.Sprintf("%v", num)
		}
		return node.Token

	case parser.ComparisonNode, parser.BinaryOpNode:
		return g.generateBinaryOp(node, level)

	default:
		var sb strings.Builder
		if node.Token != "" {
			sb.WriteString(node.Token)
		}
		for _, child := range node.Children {
			sb.WriteString(" " + g.generateNode(child, level))
		}
		return sb.String()
	}
}

func (g *Generator) generateCreate(node parser.ASTNode, level int) string {
	var sb strings.Builder
	sb.WriteString("CREATE COMPUTE MODULE")
	for _, child := range node.Children {
		if child.Type == parser.ModuleNode {
			sb.WriteString(" " + g.generateModule(child, level))
			break
		}
	}
	return sb.String()
}

func (g *Generator) generateModule(node parser.ASTNode, level int) string {
	var sb strings.Builder
	moduleName := "UnnamedModule"
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			if name, ok := child.Value.(string); ok && name != "" {
				moduleName = name
				break
			}
		}
	}

	sb.WriteString(moduleName + "\n")
	for _, child := range node.Children {
		if child.Type != parser.IdentifierNode {
			sb.WriteString(g.generateNode(child, level+1))
		}
	}
	sb.WriteString("END MODULE;\n")
	return sb.String()
}

func (g *Generator) generateDeclare(node parser.ASTNode, level int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(g.indent, level) + "DECLARE ")

	name := ""
	typ := ""
	for _, child := range node.Children {
		if child.Type == parser.IdentifierNode {
			if name == "" {
				if v, ok := child.Value.(string); ok {
					name = v
				} else {
					name = child.Token
				}
			} else {
				typ = child.Token
			}
		}
	}
	sb.WriteString(name + " " + typ + ";\n")
	return sb.String()
}

func (g *Generator) generateSet(node parser.ASTNode, level int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat(g.indent, level) + "SET ")

	var target, value string
	for _, child := range node.Children {
		if child.Type == parser.BlockNode {
			if child.Token == "target" && len(child.Children) > 0 {
				target = g.generateNode(child.Children[0], 0)
			} else if child.Token == "value" && len(child.Children) > 0 {
				value = g.generateNode(child.Children[0], 0)
			}
		}
	}

	sb.WriteString(target + " = " + value + ";\n")
	return sb.String()
}

func (g *Generator) generateIf(node parser.ASTNode, level int) string {
	var sb strings.Builder
	indentStr := strings.Repeat(g.indent, level)

	sb.WriteString(indentStr + "IF ")

	for _, child := range node.Children {
		if child.Type == parser.BlockNode && child.Token == "condition" && len(child.Children) > 0 {
			sb.WriteString(g.generateNode(child.Children[0], 0))
			break
		}
	}

	sb.WriteString(" THEN\n")

	for _, child := range node.Children {
		if child.Type == parser.BlockNode && child.Token == "then" {
			for _, stmt := range child.Children {
				sb.WriteString(g.generateNode(stmt, level+1))
			}
		}
	}

	sb.WriteString(indentStr + "END IF;\n")
	return sb.String()
}

func (g *Generator) generateBlock(node parser.ASTNode, level int) string {
	var sb strings.Builder
	for _, child := range node.Children {
		sb.WriteString(g.generateNode(child, level))
	}
	return sb.String()
}

func (g *Generator) generateFieldReference(node parser.ASTNode) string {
	if node.Value != nil {
		if path, ok := node.Value.(string); ok {
			return path
		}
	}
	var parts []string
	for _, child := range node.Children {
		parts = append(parts, g.generateNode(child, 0))
	}
	return strings.Join(parts, ".")
}

func (g *Generator) generateBinaryOp(node parser.ASTNode, level int) string {
	if len(node.Children) != 2 {
		return node.Token
	}
	left := g.generateNode(node.Children[0], 0)
	right := g.generateNode(node.Children[1], 0)
	return left + " " + node.Token + " " + right
}

```

---

### File: `pkg/analyzer/analyzer.go`

```go
package analyzer

import (
	"esql-ast-tool/pkg/parser"
)

type AnalysisResult struct {
	Variables        map[string]VariableInfo  `json:"variables"`
	Functions        map[string]FunctionInfo  `json:"functions"`
	Procedures       map[string]ProcedureInfo `json:"procedures"`
	UsedVariables    []string                 `json:"usedVariables"`
	DefinedVariables []string                 `json:"definedVariables"`
	Issues           []string                 `json:"issues"`
}

type VariableInfo struct {
	Type string `json:"type"`
	Line int    `json:"line"`
}

type FunctionInfo struct {
	Parameters []string `json:"parameters"`
	ReturnType string   `json:"returnType"`
	Line       int      `json:"line"`
}

type ProcedureInfo struct {
	Parameters []string `json:"parameters"`
	Line       int      `json:"line"`
}

type Analyzer struct {
	variables        map[string]VariableInfo
	functions        map[string]FunctionInfo
	procedures       map[string]ProcedureInfo
	usedVariables    map[string]bool
	definedVariables map[string]bool
	issues           []string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		variables:        make(map[string]VariableInfo),
		functions:        make(map[string]FunctionInfo),
		procedures:       make(map[string]ProcedureInfo),
		usedVariables:    make(map[string]bool),
		definedVariables: make(map[string]bool),
		issues:           []string{},
	}
}

func (a *Analyzer) Analyze(program parser.Program) AnalysisResult {
	for _, stmt := range program.Statements {
		a.analyzeNode(stmt)
	}

	usedVars := make([]string, 0, len(a.usedVariables))
	for v := range a.usedVariables {
		usedVars = append(usedVars, v)
	}

	definedVars := make([]string, 0, len(a.definedVariables))
	for v := range a.definedVariables {
		definedVars = append(definedVars, v)
	}

	return AnalysisResult{
		Variables:        a.variables,
		Functions:        a.functions,
		Procedures:       a.procedures,
		UsedVariables:    usedVars,
		DefinedVariables: definedVars,
		Issues:           a.issues,
	}
}

func (a *Analyzer) analyzeNode(node parser.ASTNode) {
	if node.Type == "" {
		return
	}

	switch node.Type {
	case parser.DeclareNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				varType := "UNKNOWN"
				if len(node.Children) > 1 && node.Children[1].Type == parser.IdentifierNode {
					if v, ok := node.Children[1].Value.(string); ok {
						varType = v
					}
				}
				a.variables[name] = VariableInfo{
					Type: varType,
					Line: node.Line,
				}
				a.definedVariables[name] = true
			}
		}

	case parser.SetNode:
		if len(node.Children) > 0 {
			a.analyzeNode(node.Children[0])
		}
		if len(node.Children) > 1 {
			a.analyzeNode(node.Children[1])
		}

	case parser.FunctionNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "UNKNOWN",
					Line:       node.Line,
				}
			}
		}

	case parser.ProcedureNode:
		if len(node.Children) > 0 && node.Children[0].Type == parser.IdentifierNode {
			if name, ok := node.Children[0].Value.(string); ok {
				a.procedures[name] = ProcedureInfo{
					Parameters: []string{},
					Line:       node.Line,
				}
			}
		}

	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok {
			a.usedVariables[name] = true
		}

	case parser.FunctionCallNode:
		if name, ok := node.Value.(string); ok {
			if _, exists := a.functions[name]; !exists {
				a.functions[name] = FunctionInfo{
					Parameters: []string{},
					ReturnType: "BUILTIN",
					Line:       node.Line,
				}
			}
		}
	}

	for _, child := range node.Children {
		a.analyzeNode(child)
	}
}

```

---

### File: `pkg/printer/printer.go`

```go
package printer

import (
	"fmt"
	"strings"

	"esql-ast-tool/pkg/parser"
)

type Printer struct {
	indent string
}

func NewPrinter() *Printer {
	return &Printer{indent: "  "}
}

func (p *Printer) PrintProgram(program parser.Program) string {
	var sb strings.Builder
	sb.WriteString("Program [line: 1, col: 1]\n")

	for _, stmt := range program.Statements {
		sb.WriteString(p.printNode(stmt, 1))
	}

	return sb.String()
}

func (p *Printer) printNode(node parser.ASTNode, level int) string {
	if node.Type == "" {
		return ""
	}

	indent := strings.Repeat(p.indent, level)
	var sb strings.Builder

	// Determine display name
	displayName := string(node.Type)
	switch node.Type {
	case parser.CreateNode:
		displayName = "CreateModule"
	case parser.ModuleNode:
		if node.Value != nil {
			if val, ok := node.Value.(string); ok && val == "COMPUTE" {
				displayName = "ComputeModule"
			} else {
				displayName = "Module"
			}
		} else {
			displayName = "Module"
		}
	case parser.DeclareNode:
		displayName = "Declare"
	case parser.SetNode:
		displayName = "Set"
	case parser.IfNode:
		displayName = "If"
	case parser.BlockNode:
		switch node.Token {
		case "condition":
			displayName = "Condition"
		case "then":
			displayName = "Then"
		case "else":
			displayName = "Else"
		case "target":
			displayName = "Target"
		case "value":
			displayName = "Value"
		default:
			displayName = "Block"
		}
	case parser.ComparisonNode:
		if node.Token != "" {
			displayName = fmt.Sprintf("Comparison (%s)", node.Token)
		} else {
			displayName = "Comparison"
		}
	case parser.FieldReferenceNode:
		if node.Value != nil {
			if path, ok := node.Value.(string); ok {
				displayName = fmt.Sprintf("FieldReference (%s)", path)
			} else {
				displayName = "FieldReference"
			}
		} else {
			displayName = "FieldReference"
		}
	case parser.IdentifierNode:
		if name, ok := node.Value.(string); ok && name != "" && name != "error" {
			displayName = fmt.Sprintf("Identifier: %s", name)
		} else {
			displayName = fmt.Sprintf("Identifier: %s", node.Token)
		}
	case parser.LiteralNode:
		if valStr, ok := node.Value.(string); ok {
			displayName = fmt.Sprintf("Literal: '%s'", valStr)
		} else if num, ok := node.Value.(float64); ok {
			displayName = fmt.Sprintf("Literal: %v", num)
		} else {
			displayName = "Literal"
		}
	}

	sb.WriteString(fmt.Sprintf("%s%s [line: %d, col: %d]\n", indent, displayName, node.Line, node.Column))

	// Print children
	for _, child := range node.Children {
		sb.WriteString(p.printNode(child, level+1))
	}

	return sb.String()
}

```

---

### File: `test_tokens.go`

```go
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

```

---

### File: `cmd/esql-ast/main.go`

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"esql-ast-tool/pkg/analyzer"
	"esql-ast-tool/pkg/generator"
	"esql-ast-tool/pkg/parser"
	"esql-ast-tool/pkg/printer"
)

func main() {
	var (
		file     = flag.String("f", "", "ESQL file to parse")
		code     = flag.String("c", "", "ESQL code string to parse")
		jsonOut  = flag.Bool("json", false, "Output AST as JSON")
		pretty   = flag.Bool("pretty", false, "Pretty print AST")
		analyze  = flag.Bool("analyze", false, "Perform analysis")
		validate = flag.Bool("validate", false, "Validate AST")
		generate = flag.Bool("generate", false, "Generate ESQL code from AST")
		output   = flag.String("o", "", "Output file")
	)
	flag.Parse()

	if *file == "" && *code == "" {
		fmt.Println("Usage: esql-ast -f <file> or -c <code>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var input string
	if *file != "" {
		data, err := os.ReadFile(*file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(data)
	} else {
		input = *code
	}

	p := parser.NewParser(input)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	var result string

	if *jsonOut {
		jsonData, err := program.ToJSON()
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		result = string(jsonData)
	} else if *pretty {
		pr := printer.NewPrinter()
		result = pr.PrintProgram(program)
	} else if *generate {
		gen := generator.NewGenerator()
		result = gen.Generate(program)
	} else {
		result = fmt.Sprintf("Program has %d statements\n", len(program.Statements))
	}

	if *analyze {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		result += "\n=== Analysis Results ===\n"
		result += fmt.Sprintf("Variables defined: %d\n", len(analysisResult.DefinedVariables))
		result += fmt.Sprintf("Variables used: %d\n", len(analysisResult.UsedVariables))

		if len(analysisResult.Variables) > 0 {
			result += "\nVariables:\n"
			for name, info := range analysisResult.Variables {
				result += fmt.Sprintf("  %s: %s (line %d)\n", name, info.Type, info.Line)
			}
		}

		if len(analysisResult.Functions) > 0 {
			result += "\nFunctions:\n"
			for name, info := range analysisResult.Functions {
				result += fmt.Sprintf("  %s (line %d)\n", name, info.Line)
			}
		}
	}

	if *validate {
		an := analyzer.NewAnalyzer()
		analysisResult := an.Analyze(program)

		result += "\n=== Validation Results ===\n"
		if len(analysisResult.Issues) > 0 {
			for _, issue := range analysisResult.Issues {
				result += fmt.Sprintf("  %s\n", issue)
			}
		} else {
			result += "  No issues found\n"
		}
	}

	if *output != "" {
		err := os.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing output: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(result)
	}
}

```

---

### File: `internal/token/token.go`

```go
package token

const (
	// Keywords and identifiers
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"
	EOF        = "EOF"
	ILLEGAL    = "ILLEGAL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	DOT       = "."

	// Keywords
	CREATE      = "CREATE"
	DECLARE     = "DECLARE"
	SET         = "SET"
	IF          = "IF"
	ELSE        = "ELSE"
	ELSEIF      = "ELSEIF"
	WHILE       = "WHILE"
	FOR         = "FOR"
	RETURN      = "RETURN"
	THROW       = "THROW"
	PROPAGATE   = "PROPAGATE"
	MOVE        = "MOVE"
	CONTINUE    = "CONTINUE"
	BREAK       = "BREAK"
	LABEL       = "LABEL"
	MODULE      = "MODULE"
	FUNCTION    = "FUNCTION"
	PROCEDURE   = "PROCEDURE"
	CALL        = "CALL"
	BEGIN       = "BEGIN"
	END         = "END"
	THEN        = "THEN"
	DO          = "DO"
	AS          = "AS"
	RETURNS     = "RETURNS"
	DEFAULT     = "DEFAULT"
	COMPUTE     = "COMPUTE"
	FIELD       = "FIELD"
	ENVIRONMENT = "ENVIRONMENT"
	DATABASE    = "DATABASE"
	PASSTHRU    = "PASSTHRU"
	AND         = "AND"
	OR          = "OR"
	NOT         = "NOT"
	TO          = "TO"
)

type Token struct {
	Type    string
	Literal string
	Line    int
	Column  int
	Pos     int
}

func LookupIdent(ident string) string {
	keywords := map[string]string{
		"CREATE":    CREATE,
		"DECLARE":   DECLARE,
		"SET":       SET,
		"IF":        IF,
		"ELSE":      ELSE,
		"ELSEIF":    ELSEIF,
		"WHILE":     WHILE,
		"FOR":       FOR,
		"RETURN":    RETURN,
		"THROW":     THROW,
		"PROPAGATE": PROPAGATE,
		"MOVE":      MOVE,
		"CONTINUE":  CONTINUE,
		"BREAK":     BREAK,
		"LABEL":     LABEL,
		"MODULE":    MODULE,
		"FUNCTION":  FUNCTION,
		"PROCEDURE": PROCEDURE,
		"CALL":      CALL,
		"BEGIN":     BEGIN,
		"END":       END,
		"THEN":      THEN,
		"DO":        DO,
		"AS":        AS,
		"RETURNS":   RETURNS,
		"DEFAULT":   DEFAULT,
		"COMPUTE":   COMPUTE,
		"FIELD":     FIELD,
		// "ENVIRONMENT": ENVIRONMENT,
		"DATABASE": DATABASE,
		"PASSTHRU": PASSTHRU,
		"AND":      AND,
		"OR":       OR,
		"NOT":      NOT,
		"TO":       TO,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

```

---

