# Dokumentasi Struktur & Kode Project

Dokumen ini dihasilkan secara otomatis untuk memetakan seluruh struktur folder dan isi kode di dalam project ini.

## Struktur Project (Tree)

```text
.
├── cmd
│   └── esql-ast
│       └── main.go
├── DOC.md
├── esql-ast
├── examples
│   ├── sample.esql
│   ├── test_case.esql
│   ├── test_case_nested_if.esql
│   ├── test_case_searched_only.esql
│   ├── test_case_simple.esql
│   ├── test_case_simple_only.esql
│   ├── test_cast.esql
│   └── test_nested_cast.esql
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
│   │   ├── parser_expr.go
│   │   ├── parser.go
│   │   ├── parser_primary.go
│   │   ├── parser_stmt_call.go
│   │   ├── parser_stmt_control.go
│   │   ├── parser_stmt_create.go
│   │   ├── parser_stmt_declare.go
│   │   ├── parser_stmt.go
│   │   ├── parser_stmt_if.go
│   │   ├── parser_stmt_loop.go
│   │   ├── parser_stmt_set.go
│   │   └── parser_utils.go
│   └── printer
│       └── printer.go
└── tests
    ├── baseline
    │   ├── case_full.analyze.txt
    │   ├── case_full.generate.txt
    │   ├── case_full.pretty.txt
    │   ├── case_nested_if.analyze.txt
    │   ├── case_nested_if.generate.txt
    │   ├── case_nested_if.pretty.txt
    │   ├── case_searched.analyze.txt
    │   ├── case_searched.generate.txt
    │   ├── case_searched.pretty.txt
    │   ├── case_simple.analyze.txt
    │   ├── case_simple.generate.txt
    │   ├── case_simple.pretty.txt
    │   ├── cast.analyze.txt
    │   ├── cast.generate.txt
    │   ├── cast.json.txt
    │   ├── cast.pretty.txt
    │   ├── nested_cast.generate.txt
    │   ├── nested_cast.pretty.txt
    │   ├── sample.analyze.txt
    │   ├── sample.generate.txt
    │   ├── sample.json.txt
    │   └── sample.pretty.txt
    ├── diff
    │   ├── case_full.analyze.txt.diff
    │   ├── case_nested_if.analyze.txt.diff
    │   ├── case_searched.analyze.txt.diff
    │   └── cast.analyze.txt.diff
    └── output
        ├── case_full.analyze.txt
        ├── case_full.generate.txt
        ├── case_full.pretty.txt
        ├── case_nested_if.analyze.txt
        ├── case_nested_if.generate.txt
        ├── case_nested_if.pretty.txt
        ├── case_searched.analyze.txt
        ├── case_searched.generate.txt
        ├── case_searched.pretty.txt
        ├── case_simple.analyze.txt
        ├── case_simple.generate.txt
        ├── case_simple.pretty.txt
        ├── cast.analyze.txt
        ├── cast.generate.txt
        ├── cast.json.txt
        ├── cast.pretty.txt
        ├── nested_cast.generate.txt
        ├── nested_cast.pretty.txt
        ├── sample.analyze.txt
        ├── sample.generate.txt
        ├── sample.json.txt
        └── sample.pretty.txt

15 directories, 80 files
```

## Isi Kode Berdasarkan File

### File: `go.mod`

```text
module esql-ast-tool

go 1.22

```

---

### File: `examples/test_case.esql`

```text
CREATE COMPUTE MODULE TestCase
    DECLARE score INTEGER;
    DECLARE grade STRING;
    DECLARE status STRING;
    DECLARE result INTEGER;
    
    SET score = 85;
    
    -- Simple CASE
    SET grade = CASE score
        WHEN 90 THEN 'A'
        WHEN 80 THEN 'B'
        WHEN 70 THEN 'C'
        ELSE 'D'
    END;
    
    -- Searched CASE
    SET status = CASE 
        WHEN score >= 90 THEN 'Excellent'
        WHEN score >= 80 THEN 'Good'
        WHEN score >= 70 THEN 'Average'
        ELSE 'Poor'
    END;
    
    -- CASE in IF condition
    IF CASE WHEN score > 80 THEN 1 ELSE 0 END = 1 THEN
        SET result = 100;
    END IF;
    
    -- Nested CASE
    SET result = CASE 
        WHEN CASE WHEN score > 80 THEN 1 ELSE 0 END = 1 THEN 100
        ELSE 0
    END;
END MODULE;

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

### File: `examples/test_case_simple.esql`

```text
CREATE COMPUTE MODULE TestCaseSimple
    DECLARE score INTEGER;
    DECLARE grade STRING;
    SET score = 85;
    SET grade = CASE score
        WHEN 90 THEN 'A'
        WHEN 80 THEN 'B'
        ELSE 'C'
    END;
END MODULE;

```

---

### File: `examples/test_case_searched_only.esql`

```text
CREATE COMPUTE MODULE TestCaseSearched
    DECLARE score INTEGER;
    DECLARE status STRING;
    SET score = 85;
    SET status = CASE 
        WHEN score >= 90 THEN 'Excellent'
        WHEN score >= 80 THEN 'Good'
        WHEN score >= 70 THEN 'Average'
        ELSE 'Poor'
    END;
END MODULE;

```

---

### File: `examples/test_cast.esql`

```text
CREATE COMPUTE MODULE TestCast
    DECLARE intVar INTEGER;
    DECLARE strVar STRING;
    DECLARE result INTEGER;
    
    SET strVar = '123';
    SET intVar = CAST(strVar AS INTEGER);
    SET result = CAST('456' AS INTEGER);
    SET strVar = CAST(789 AS STRING);
    
    IF CAST(strVar AS INTEGER) > 100 THEN
        SET result = 1;
    END IF;
END MODULE;

```

---

### File: `examples/test_nested_cast.esql`

```text
CREATE COMPUTE MODULE TestNestedCast
    DECLARE strVar STRING;
    DECLARE intVar INTEGER;
    DECLARE result INTEGER;
    
    SET strVar = '123';
    SET intVar = CAST(CAST(strVar AS STRING) AS INTEGER);
    SET result = CAST(CAST('456' AS INTEGER) AS INTEGER);
    SET strVar = CAST(CAST(789 AS STRING) AS STRING);
    
    IF CAST(CAST(strVar AS STRING) AS INTEGER) > 100 THEN
        SET result = 1;
    END IF;
END MODULE;

```

---

### File: `examples/test_case_simple_only.esql`

```text
CREATE COMPUTE MODULE TestCaseSimple
    DECLARE score INTEGER;
    DECLARE grade STRING;
    SET score = 85;
    SET grade = CASE score
        WHEN 90 THEN 'A'
        WHEN 80 THEN 'B'
        ELSE 'C'
    END;
END MODULE;

```

---

### File: `examples/test_case_nested_if.esql`

```text
CREATE COMPUTE MODULE TestCaseNestedIf
    DECLARE score INTEGER;
    DECLARE result INTEGER;
    SET score = 85;
    IF CASE WHEN score > 80 THEN 1 ELSE 0 END = 1 THEN
        SET result = 100;
    END IF;
END MODULE;

```

---

### File: `pkg/parser/parser_stmt_set.go`

```go
package parser

import (
	"fmt"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parseSet() ASTNode {
	debugPrint("  [parseSet] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(SetNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse target
	var target ASTNode
	switch p.curToken.Type {
	case token.IDENTIFIER:
		target = p.parseIdentifier()
	case token.ENVIRONMENT:
		target = p.parseFieldReferenceFromKeyword("Environment")
	case token.FIELD:
		target = p.parseField()
	default:
		target = p.parseExpression()
	}

	if target.Type != "" {
		targetWrapper := NewASTNode(BlockNode, "target", target.Line, target.Column)
		targetWrapper.AddChild(target)
		node.AddChild(targetWrapper)
	}

	// Parse '='
	if p.curToken.Type == token.ASSIGN || p.curToken.Type == token.EQ {
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
		p.nextToken()
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	debugPrint("  [parseSet] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
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
	CastNode           NodeType = "Cast"
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

### File: `pkg/parser/parser_stmt_create.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

// ============================================
// CREATE COMPUTE MODULE
// ============================================

func (p *Parser) parseCreate() ASTNode {
	node := NewASTNode(CreateNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.COMPUTE {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
			moduleNode := p.parseComputeModule()
			node.AddChild(moduleNode)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseComputeModule() ASTNode {
	moduleNode := NewASTNode(ModuleNode, "COMPUTE MODULE", p.curToken.Line, p.curToken.Column)
	moduleNode.Value = "COMPUTE"

	// Parse module name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		moduleNode.AddChild(nameNode)
		p.nextToken()
	}

	// Parse body
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			moduleNode.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}

	// Consume END MODULE
	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
		}
	}

	return moduleNode
}

// ============================================
// MODULE Statement
// ============================================

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

	// Parse body
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// FUNCTION Statement
// ============================================

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
		p.parseFunctionParameters(&node)
	}

	// Parse return type
	if p.curToken.Type == token.RETURNS {
		p.parseFunctionReturnType(&node)
	}

	// Parse function body
	if p.curToken.Type == token.BEGIN {
		p.parseFunctionBody(&node)
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFunctionParameters(node *ASTNode) {
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

func (p *Parser) parseFunctionReturnType(node *ASTNode) {
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

func (p *Parser) parseFunctionBody(node *ASTNode) {
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

// ============================================
// PROCEDURE Statement
// ============================================

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
		p.parseProcedureParameters(&node)
	}

	// Parse procedure body
	if p.curToken.Type == token.BEGIN {
		p.parseProcedureBody(&node)
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseProcedureParameters(node *ASTNode) {
	p.nextToken()
	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		if p.curToken.Type == token.IDENTIFIER {
			paramNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			if paramNode.Value == nil {
				paramNode.Value = paramNode.Token
			}
			node.AddChild(paramNode)
			p.nextToken()
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

func (p *Parser) parseProcedureBody(node *ASTNode) {
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

```

---

### File: `pkg/parser/parser_primary.go`

```go
package parser

import (
	"fmt"
	"strconv"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parsePrimary() ASTNode {
	debugPrint("    [parsePrimary] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	switch p.curToken.Type {
	case token.IDENTIFIER:
		result := p.parseIdentifier()
		debugPrint("    [parsePrimary] after IDENTIFIER: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.NUMBER:
		result := p.parseNumber()
		debugPrint("    [parsePrimary] after NUMBER: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.STRING:
		result := p.parseString()
		debugPrint("    [parsePrimary] after STRING: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.LPAREN:
		result := p.parseGroupedExpression()
		debugPrint("    [parsePrimary] after LPAREN: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.CASE:
		result := p.parseCase()
		debugPrint("    [parsePrimary] after CASE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.CAST:
		result := p.parseCast()
		debugPrint("    [parsePrimary] after CAST: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.DOT:
		debugPrint("    [parsePrimary] WARNING: DOT without identifier, skipping...\n")
		p.nextToken()
		return ASTNode{}
	default:
		debugPrint("    [parsePrimary] UNKNOWN: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return ASTNode{}
	}
}

func (p *Parser) parseIdentifier() ASTNode {
	debugPrint("      [parseIdentifier] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type == token.LPAREN {
		return p.parseFunctionCall(node)
	}

	if p.curToken.Type == token.DOT {
		return p.parseFieldReference(node)
	}

	debugPrint("      [parseIdentifier] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

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
	fieldNode := NewASTNode(FieldReferenceNode, "field", base.Line, base.Column)
	fieldNode.AddChild(base)

	if base.Value != nil {
		fieldNode.Value = base.Value
	} else {
		fieldNode.Value = base.Token
	}

	for p.curToken.Type == token.DOT {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			identNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			identNode.Value = p.curToken.Literal

			newFieldNode := NewASTNode(FieldReferenceNode, "field", fieldNode.Line, fieldNode.Column)
			newFieldNode.AddChild(fieldNode)
			newFieldNode.AddChild(identNode)

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

func (p *Parser) parseCast() ASTNode {
	debugPrint("    [parseCast] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(CastNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = "CAST"
	p.nextToken()

	debugPrint("    [parseCast] after CAST: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.LPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '(' after CAST, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] after '(': token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	expr := p.parseExpression()
	if expr.Type != "" {
		node.AddChild(expr)
	}

	debugPrint("    [parseCast] after expression: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.AS {
		p.errors = append(p.errors,
			fmt.Sprintf("expected 'AS' in CAST expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] after AS: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors,
			fmt.Sprintf("expected type name after AS in CAST, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}

	typeNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	typeNode.Value = p.curToken.Literal
	node.AddChild(typeNode)
	p.nextToken()

	debugPrint("    [parseCast] after type: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.RPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ')' in CAST expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseCase() ASTNode {
	debugPrint("[parseCase] START: token=%s, literal='%s', line=%d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(CaseNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("[parseCase] after CASE: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	var isSimpleCase bool
	var caseExpr ASTNode

	if p.curToken.Type != token.WHEN {
		isSimpleCase = true
		debugPrint("[parseCase] Simple CASE, parsing expression\n")
		caseExpr = p.parseExpression()
		if caseExpr.Type != "" {
			node.AddChild(caseExpr)
		}
		debugPrint("[parseCase] after expression: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	} else {
		debugPrint("[parseCase] Searched CASE\n")
	}

	whenCount := 0
	for p.curToken.Type == token.WHEN {
		whenCount++
		debugPrint("[parseCase] Parsing WHEN #%d at line %d\n", whenCount, p.curToken.Line)
		whenNode := p.parseWhen(isSimpleCase)
		if whenNode.Type != "" {
			node.AddChild(whenNode)
		}
		debugPrint("[parseCase] after WHEN #%d: token=%s, literal='%s'\n",
			whenCount, p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.ELSE {
		debugPrint("[parseCase] Parsing ELSE\n")
		p.nextToken()
		elseExpr := p.parseExpression()
		if elseExpr.Type != "" {
			elseNode := NewASTNode(BlockNode, "else", elseExpr.Line, elseExpr.Column)
			elseNode.AddChild(elseExpr)
			node.AddChild(elseNode)
		}
		debugPrint("[parseCase] after ELSE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.END {
		debugPrint("[parseCase] Found END, consuming it\n")
		p.nextToken()
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected END in CASE expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
	}

	debugPrint("[parseCase] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.THEN || p.curToken.Type == token.ELSE || p.curToken.Type == token.END {
		debugPrint("[parseCase] Next token is %s, returning early\n", p.curToken.Type)
		return node
	}

	return node
}

func (p *Parser) parseWhen(isSimpleCase bool) ASTNode {
	debugPrint("  [parseWhen] START: token=%s, literal='%s', isSimple=%v\n",
		p.curToken.Type, p.curToken.Literal, isSimpleCase)

	node := NewASTNode(WhenNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("  [parseWhen] after WHEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	var condition ASTNode

	if isSimpleCase {
		debugPrint("  [parseWhen] Simple CASE: parsing value\n")
		condition = p.parseExpression()
		if condition.Type != "" {
			node.AddChild(condition)
		}
		debugPrint("  [parseWhen] after value: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	} else {
		debugPrint("  [parseWhen] Searched CASE: parsing condition\n")
		condition = p.parseExpression()
		if condition.Type != "" {
			node.AddChild(condition)
		}
		debugPrint("  [parseWhen] after condition: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.THEN {
		debugPrint("  [parseWhen] Found THEN\n")
		p.nextToken()
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected THEN in CASE WHEN, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}

	debugPrint("  [parseWhen] after THEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	result := p.parseExpression()
	if result.Type != "" {
		node.AddChild(result)
	}

	debugPrint("  [parseWhen] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

```

---

### File: `pkg/parser/parser_expr.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseExpression() ASTNode {
	debugPrint("  [parseExpression] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	left := p.parseLogicalOr()

	if p.curToken.Type == token.EQ || p.curToken.Type == token.ASSIGN ||
		p.curToken.Type == token.NOT_EQ ||
		p.curToken.Type == token.LT || p.curToken.Type == token.GT ||
		p.curToken.Type == token.LTE || p.curToken.Type == token.GTE {
		debugPrint("  [parseExpression] FOUND OPERATOR: %s\n", p.curToken.Literal)
		tok := p.curToken
		p.nextToken()
		right := p.parseAdditive()
		if right.Type != "" {
			compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
			compNode.AddChild(left)
			compNode.AddChild(right)
			debugPrint("  [parseExpression] returning comparison node\n")
			return compNode
		}
	}

	return left
}

func (p *Parser) parseLogicalOr() ASTNode {
	debugPrint("    [parseLogicalOr] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
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
	debugPrint("    [parseLogicalAnd] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
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
	debugPrint("    [parseComparison] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := p.parseAdditive()

	debugPrint("    [parseComparison] after parseAdditive: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.EQ || p.curToken.Type == token.NOT_EQ ||
		p.curToken.Type == token.LT || p.curToken.Type == token.GT ||
		p.curToken.Type == token.LTE || p.curToken.Type == token.GTE {
		debugPrint("    [parseComparison] found operator: %s\n", p.curToken.Literal)
		tok := p.curToken
		p.nextToken()
		right := p.parseAdditive()
		if right.Type != "" {
			compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
			compNode.AddChild(node)
			compNode.AddChild(right)
			debugPrint("    [parseComparison] returning comparison node\n")
			return compNode
		}
	}

	debugPrint("    [parseComparison] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseAdditive() ASTNode {
	debugPrint("    [parseAdditive] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

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

	debugPrint("    [parseAdditive] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseMultiplicative() ASTNode {
	debugPrint("    [parseMultiplicative] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

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

	debugPrint("    [parseMultiplicative] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseUnary() ASTNode {
	debugPrint("    [parseUnary] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.MINUS || p.curToken.Type == token.NOT {
		tok := p.curToken
		p.nextToken()
		operand := p.parsePrimary()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		debugPrint("    [parseUnary] END (unary): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	result := p.parsePrimary()
	debugPrint("    [parseUnary] END (primary): token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
	return result
}

```

---

### File: `pkg/parser/parser.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

var DebugMode = false

func debugPrint(format string, args ...interface{}) {
	if DebugMode {
		fmt.Printf(format, args...)
	}
}

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

func (p *Parser) GetCurToken() token.Token {
	return p.curToken
}

func (p *Parser) GetPeekToken() token.Token {
	return p.peekToken
}

func (p *Parser) GetNextToken() {
	p.nextToken()
}

func (p *Parser) GetPosition() int {
	return p.position
}

func (p *Parser) DebugTokens() {
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

func (p *Parser) GetTokens() []token.Token {
	return p.tokens
}

func (p *Parser) DebugPrintTokens() {
	fmt.Println("=== TOKENS ===")
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// ParseProgram - entry point
func (p *Parser) ParseProgram() Program {
	program := NewProgram()

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			program.AddStatement(stmt)
		}

		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.EOF && p.curToken.Type != token.END {
			p.nextToken()
		}
	}

	return program
}

```

---

### File: `pkg/parser/parser_stmt_loop.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

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

```

---

### File: `pkg/parser/parser_stmt_if.go`

```go
package parser

import (
	"fmt"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parseIf() ASTNode {
	debugPrint("[parseIf] START: token=%s, literal='%s', line=%d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(IfNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("[parseIf] after IF: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Parse condition
	cond := p.parseExpression()
	if cond.Type != "" {
		debugPrint("[parseIf] condition parsed: type=%s\n", cond.Type)
		node.AddChild(cond)
	}

	debugPrint("[parseIf] after condition: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Expect THEN
	if p.curToken.Type == token.THEN {
		debugPrint("[parseIf] Found THEN, consuming it\n")
		p.nextToken()
	} else {
		debugPrint("[parseIf] ERROR: expected THEN, got %s\n", p.curToken.Type)
		p.errors = append(p.errors,
			fmt.Sprintf("expected THEN after IF condition, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		p.nextToken()
		return node
	}

	debugPrint("[parseIf] after THEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

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
		debugPrint("[parseIf] Found ELSE\n")
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
		debugPrint("[parseIf] Found END\n")
		p.nextToken()
		if p.curToken.Type == token.IF {
			debugPrint("[parseIf] Found IF after END\n")
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	debugPrint("[parseIf] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

```

---

### File: `pkg/parser/parser_stmt_call.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

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

```

---

### File: `pkg/parser/parser_stmt_declare.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

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

```

---

### File: `pkg/parser/parser_stmt.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseStatement() ASTNode {
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
		p.nextToken()
		return ASTNode{}
	default:
		expr := p.parseExpression()
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return expr
	}
}

```

---

### File: `pkg/parser/parser_stmt_control.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

// ============================================
// RETURN
// ============================================

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

// ============================================
// THROW
// ============================================

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

// ============================================
// PROPAGATE
// ============================================

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

// ============================================
// MOVE
// ============================================

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

// ============================================
// CONTINUE
// ============================================

func (p *Parser) parseContinue() ASTNode {
	node := NewASTNode(ContinueNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	if p.curToken.Type == token.IDENTIFIER {
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return node
}

// ============================================
// BREAK
// ============================================

func (p *Parser) parseBreak() ASTNode {
	node := NewASTNode(BreakNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// ============================================
// LABEL
// ============================================

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

### File: `pkg/parser/parser_utils.go`

```go
package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseFieldReferenceFromKeyword(keyword string) ASTNode {
	debugPrint("  [parseFieldReferenceFromKeyword] keyword=%s\n", keyword)

	baseNode := NewASTNode(IdentifierNode, keyword, p.curToken.Line, p.curToken.Column)
	baseNode.Value = keyword
	p.nextToken()

	return p.parseFieldReference(baseNode)
}

func (p *Parser) parseField() ASTNode {
	debugPrint("  [parseField] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(FieldReferenceNode, "FIELD", p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER || p.curToken.Type == token.ENVIRONMENT {
		base := p.parseIdentifier()
		node.AddChild(base)

		for p.curToken.Type == token.DOT {
			node = p.parseFieldReference(node)
		}
	}

	debugPrint("  [parseField] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseExpressionStatement() ASTNode {
	expr := p.parseExpression()
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return expr
}

func (p *Parser) parseComparisonFromNode(left ASTNode) ASTNode {
	node := left

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

	case parser.CastNode:
		return g.generateCast(node, level)

	case parser.CaseNode:
		return g.generateCase(node, level)

	case parser.WhenNode:
		return g.generateWhen(node, false)

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

	// Generate condition (child 0)
	if len(node.Children) > 0 {
		cond := node.Children[0]
		sb.WriteString(g.generateNode(cond, 0))
	}

	sb.WriteString(" THEN\n")

	// Generate then block (child 1)
	if len(node.Children) > 1 {
		thenBlock := node.Children[1]
		for _, stmt := range thenBlock.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}

	// Generate else block if exists (child 2)
	if len(node.Children) > 2 {
		elseBlock := node.Children[2]
		sb.WriteString(indentStr + "ELSE\n")
		for _, stmt := range elseBlock.Children {
			sb.WriteString(g.generateNode(stmt, level+1))
		}
	}

	sb.WriteString(indentStr + "END IF;\n")
	return sb.String()
}

func (g *Generator) generateBlock(node parser.ASTNode, level int) string {
	var sb strings.Builder

	// Untuk BlockNode dengan token "else" di CASE expression
	if node.Token == "else" && len(node.Children) > 0 {
		sb.WriteString("ELSE " + g.generateNode(node.Children[0], 0))
		return sb.String()
	}

	// Untuk BlockNode biasa (then, target, value, condition)
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

func (g *Generator) generateCast(node parser.ASTNode, level int) string {
	var sb strings.Builder

	sb.WriteString("CAST(")

	// Generate expression
	if len(node.Children) > 0 {
		sb.WriteString(g.generateNode(node.Children[0], 0))
	}

	sb.WriteString(" AS ")

	// Generate type
	if len(node.Children) > 1 {
		sb.WriteString(g.generateNode(node.Children[1], 0))
	}

	sb.WriteString(")")

	return sb.String()
}

func (g *Generator) generateCase(node parser.ASTNode, level int) string {
	var sb strings.Builder

	sb.WriteString("CASE")

	// Cek apakah ini simple CASE (child 0 adalah expression, bukan WHEN)
	if len(node.Children) > 0 && node.Children[0].Type != parser.WhenNode {
		// Simple CASE: CASE expression
		sb.WriteString(" " + g.generateNode(node.Children[0], 0))

		// Generate WHEN clauses (mulai dari index 1)
		for i := 1; i < len(node.Children); i++ {
			child := node.Children[i]
			if child.Type == parser.WhenNode {
				sb.WriteString(" " + g.generateWhen(child, true))
			} else if child.Type == parser.BlockNode && child.Token == "else" {
				sb.WriteString(" " + g.generateNode(child, 0))
			}
		}
	} else {
		// Searched CASE: CASE WHEN condition THEN result ...
		for _, child := range node.Children {
			if child.Type == parser.WhenNode {
				sb.WriteString(" " + g.generateWhen(child, false))
			} else if child.Type == parser.BlockNode && child.Token == "else" {
				sb.WriteString(" " + g.generateNode(child, 0))
			}
		}
	}

	sb.WriteString(" END")
	return sb.String()
}

func (g *Generator) generateWhen(node parser.ASTNode, isSimpleCase bool) string {
	var sb strings.Builder

	sb.WriteString("WHEN ")

	if len(node.Children) >= 2 {
		if isSimpleCase {
			// Simple CASE: WHEN value THEN result
			sb.WriteString(g.generateNode(node.Children[0], 0))
			sb.WriteString(" THEN ")
			sb.WriteString(g.generateNode(node.Children[1], 0))
		} else {
			// Searched CASE: WHEN condition THEN result
			sb.WriteString(g.generateNode(node.Children[0], 0))
			sb.WriteString(" THEN ")
			sb.WriteString(g.generateNode(node.Children[1], 0))
		}
	}

	return sb.String()
}

```

---

### File: `pkg/analyzer/analyzer.go`

```go
package analyzer

import (
	"esql-ast-tool/pkg/parser"
	"sort"
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

	// Sort variable names for consistent output
	var varNames []string
	for name := range a.variables {
		varNames = append(varNames, name)
	}
	sort.Strings(varNames)

	// Rebuild variables with sorted order
	sortedVariables := make(map[string]VariableInfo)
	for _, name := range varNames {
		sortedVariables[name] = a.variables[name]
	}

	// Sort function names
	var funcNames []string
	for name := range a.functions {
		funcNames = append(funcNames, name)
	}
	sort.Strings(funcNames)

	sortedFunctions := make(map[string]FunctionInfo)
	for _, name := range funcNames {
		sortedFunctions[name] = a.functions[name]
	}

	// Sort procedure names
	var procNames []string
	for name := range a.procedures {
		procNames = append(procNames, name)
	}
	sort.Strings(procNames)

	sortedProcedures := make(map[string]ProcedureInfo)
	for _, name := range procNames {
		sortedProcedures[name] = a.procedures[name]
	}

	// Sort used variables
	usedVars := make([]string, 0, len(a.usedVariables))
	for v := range a.usedVariables {
		usedVars = append(usedVars, v)
	}
	sort.Strings(usedVars)

	// Sort defined variables
	definedVars := make([]string, 0, len(a.definedVariables))
	for v := range a.definedVariables {
		definedVars = append(definedVars, v)
	}
	sort.Strings(definedVars)

	return AnalysisResult{
		Variables:        sortedVariables,
		Functions:        sortedFunctions,
		Procedures:       sortedProcedures,
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
	case parser.CastNode:
		// Cast doesn't define variables, just analyze children
		for _, child := range node.Children {
			a.analyzeNode(child)
		}

	case parser.CaseNode:
		// Analyze semua children (WHEN clauses dan ELSE)
		for _, child := range node.Children {
			a.analyzeNode(child)
		}
	case parser.WhenNode:
		// Analyze condition dan result
		for _, child := range node.Children {
			a.analyzeNode(child)
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
		// Tampilkan condition jika ada
		if len(node.Children) > 0 {
			sb.WriteString(fmt.Sprintf("%s%s [line: %d, col: %d]\n", indent, displayName, node.Line, node.Column))
			// Print condition
			cond := node.Children[0]
			sb.WriteString(p.printNode(cond, level+1))
			// Print then block
			if len(node.Children) > 1 {
				sb.WriteString(p.printNode(node.Children[1], level+1))
			}
			// Print else block if exists
			if len(node.Children) > 2 {
				sb.WriteString(p.printNode(node.Children[2], level+1))
			}
			return sb.String()
		}

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

	case parser.CastNode:
		displayName = "Cast"
	case parser.CaseNode:
		displayName = "Case"
	case parser.WhenNode:
		displayName = "When"

	}

	sb.WriteString(fmt.Sprintf("%s%s [line: %d, col: %d]\n", indent, displayName, node.Line, node.Column))

	for _, child := range node.Children {
		sb.WriteString(p.printNode(child, level+1))
	}

	return sb.String()
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
	"sort"

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
		debug    = flag.Bool("debug", false, "Enable debug output")
	)
	flag.Parse()

	// Enable debug mode if flag is set
	if *debug {
		parser.DebugMode = true
	}

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
			// Sort variable names for consistent output
			var varNames []string
			for name := range analysisResult.Variables {
				varNames = append(varNames, name)
			}
			sort.Strings(varNames)
			for _, name := range varNames {
				info := analysisResult.Variables[name]
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
	CAST        = "CAST"
	CASE        = "CASE"
	WHEN        = "WHEN"
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
		"CREATE":      CREATE,
		"DECLARE":     DECLARE,
		"SET":         SET,
		"IF":          IF,
		"ELSE":        ELSE,
		"ELSEIF":      ELSEIF,
		"WHILE":       WHILE,
		"FOR":         FOR,
		"RETURN":      RETURN,
		"THROW":       THROW,
		"PROPAGATE":   PROPAGATE,
		"MOVE":        MOVE,
		"CONTINUE":    CONTINUE,
		"BREAK":       BREAK,
		"LABEL":       LABEL,
		"MODULE":      MODULE,
		"FUNCTION":    FUNCTION,
		"PROCEDURE":   PROCEDURE,
		"CALL":        CALL,
		"BEGIN":       BEGIN,
		"END":         END,
		"THEN":        THEN,
		"DO":          DO,
		"AS":          AS,
		"RETURNS":     RETURNS,
		"DEFAULT":     DEFAULT,
		"COMPUTE":     COMPUTE,
		"FIELD":       FIELD,
		"ENVIRONMENT": ENVIRONMENT,
		"DATABASE":    DATABASE,
		"PASSTHRU":    PASSTHRU,
		"AND":         AND,
		"OR":          OR,
		"NOT":         NOT,
		"TO":          TO,
		"CAST":        CAST,
		"CASE":        CASE,
		"WHEN":        WHEN,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

```

---

