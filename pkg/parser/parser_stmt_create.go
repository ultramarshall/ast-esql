package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

// ============================================
// CREATE COMPUTE MODULE
// ============================================

func (p *Parser) parseCreate() ASTNode {
	p.nextToken() // consume CREATE

	if p.curToken.Type == token.COMPUTE {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
			// Langsung return module node, tanpa dibungkus CreateNode
			return p.parseComputeModule()
		}
	}

	// Fallback
	return ASTNode{}
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

	// Parse body - LANGSUNG parse statement, JANGAN lewat parseCreate lagi
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
		endLine := p.curToken.Line
		endCol := p.curToken.Column + 3
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			endCol = p.curToken.Column + len(p.curToken.Literal)
			p.nextToken()
		}
		moduleNode.Span.End = Position{Line: endLine, Column: endCol}
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
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected function name, got %s", p.curToken.Type))
		return node
	}

	// Parse parameters (harus ada LPAREN)
	if p.curToken.Type == token.LPAREN {
		p.parseFunctionParameters(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' after function name, got %s", p.curToken.Type))
		return node
	}

	// Parse return type
	if p.curToken.Type == token.RETURNS {
		p.parseFunctionReturnType(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected RETURNS, got %s", p.curToken.Type))
		return node
	}

	// Parse function body
	if p.curToken.Type == token.BEGIN {
		p.parseFunctionBody(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected BEGIN, got %s", p.curToken.Type))
		return node
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
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

func (p *Parser) parseFunctionParameters(node *ASTNode) {
	p.nextToken() // consume '('

	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		// Simpan posisi awal parameter
		paramStartLine := p.curToken.Line
		paramStartCol := p.curToken.Column

		// 1. Parse mode
		var mode string
		var modeStartLine, modeStartCol int
		if p.curToken.Type == token.IN || p.curToken.Type == token.OUT || p.curToken.Type == token.INOUT {
			mode = p.curToken.Literal
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
			p.nextToken()
		} else {
			mode = "IN"
			// Jika tidak ada mode, posisi awal adalah nama parameter
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
		}

		// 2. Parse parameter name
		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter name, got %s", p.curToken.Type))
			break
		}
		nameStartLine := p.curToken.Line
		nameStartCol := p.curToken.Column
		paramName := p.curToken.Literal
		p.nextToken()

		// 3. Parse parameter type
		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter type, got %s", p.curToken.Type))
			break
		}
		typeStartLine := p.curToken.Line
		typeStartCol := p.curToken.Column
		paramType := p.curToken.Literal
		p.nextToken()

		// 4. Create Parameter node
		paramNode := NewASTNode(ParameterNode, "parameter", paramStartLine, paramStartCol)

		// Mode node
		modeNode := NewASTNode(LiteralNode, mode, modeStartLine, modeStartCol)
		modeNode.Value = mode
		modeNode.Span = Span{
			Start: Position{Line: modeStartLine, Column: modeStartCol},
			End:   Position{Line: modeStartLine, Column: modeStartCol + len(mode)},
		}
		paramNode.AddChild(modeNode)

		// Name node
		nameNode := NewASTNode(IdentifierNode, paramName, nameStartLine, nameStartCol)
		nameNode.Value = paramName
		nameNode.Span = Span{
			Start: Position{Line: nameStartLine, Column: nameStartCol},
			End:   Position{Line: nameStartLine, Column: nameStartCol + len(paramName)},
		}
		paramNode.AddChild(nameNode)

		// Type node
		typeNode := NewASTNode(IdentifierNode, paramType, typeStartLine, typeStartCol)
		typeNode.Value = paramType
		typeNode.Span = Span{
			Start: Position{Line: typeStartLine, Column: typeStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		paramNode.AddChild(typeNode)

		// Span parameter: dari awal parameter sampai akhir type
		paramNode.Span = Span{
			Start: Position{Line: paramStartLine, Column: paramStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		node.AddChild(paramNode)

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}
}

func (p *Parser) parseFunctionReturnType(node *ASTNode) {
	p.nextToken() // consume RETURNS
	if p.curToken.Type == token.IDENTIFIER {
		returnType := NewASTNode(ReturnTypeNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		returnType.Value = p.curToken.Literal
		node.AddChild(returnType)
		p.nextToken()
	}
}

func (p *Parser) parseProcedureBody(node *ASTNode) {
	p.nextToken() // consume BEGIN
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}
		// parseStatement sudah memajukan token ke setelah statement,
		// tidak perlu p.nextToken() lagi di sini.
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
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected procedure name, got %s", p.curToken.Type))
		return node
	}

	// Parse parameters
	if p.curToken.Type == token.LPAREN {
		p.parseProcedureParameters(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' after procedure name, got %s", p.curToken.Type))
		return node
	}

	// Parse procedure body
	if p.curToken.Type == token.BEGIN {
		p.parseProcedureBody(&node)
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected BEGIN, got %s", p.curToken.Type))
		return node
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseProcedureParameters(node *ASTNode) {
	p.nextToken() // consume '('

	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		paramStartLine := p.curToken.Line
		paramStartCol := p.curToken.Column

		var mode string
		var modeStartLine, modeStartCol int
		if p.curToken.Type == token.IN || p.curToken.Type == token.OUT || p.curToken.Type == token.INOUT {
			mode = p.curToken.Literal
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
			p.nextToken()
		} else {
			mode = "IN"
			modeStartLine = p.curToken.Line
			modeStartCol = p.curToken.Column
		}

		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter name, got %s", p.curToken.Type))
			break
		}
		nameStartLine := p.curToken.Line
		nameStartCol := p.curToken.Column
		paramName := p.curToken.Literal
		p.nextToken()

		if p.curToken.Type != token.IDENTIFIER {
			p.errors = append(p.errors, fmt.Sprintf("expected parameter type, got %s", p.curToken.Type))
			break
		}
		typeStartLine := p.curToken.Line
		typeStartCol := p.curToken.Column
		paramType := p.curToken.Literal
		p.nextToken()

		paramNode := NewASTNode(ParameterNode, "parameter", paramStartLine, paramStartCol)

		modeNode := NewASTNode(LiteralNode, mode, modeStartLine, modeStartCol)
		modeNode.Value = mode
		modeNode.Span = Span{
			Start: Position{Line: modeStartLine, Column: modeStartCol},
			End:   Position{Line: modeStartLine, Column: modeStartCol + len(mode)},
		}
		paramNode.AddChild(modeNode)

		nameNode := NewASTNode(IdentifierNode, paramName, nameStartLine, nameStartCol)
		nameNode.Value = paramName
		nameNode.Span = Span{
			Start: Position{Line: nameStartLine, Column: nameStartCol},
			End:   Position{Line: nameStartLine, Column: nameStartCol + len(paramName)},
		}
		paramNode.AddChild(nameNode)

		typeNode := NewASTNode(IdentifierNode, paramType, typeStartLine, typeStartCol)
		typeNode.Value = paramType
		typeNode.Span = Span{
			Start: Position{Line: typeStartLine, Column: typeStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		paramNode.AddChild(typeNode)

		paramNode.Span = Span{
			Start: Position{Line: paramStartLine, Column: paramStartCol},
			End:   Position{Line: typeStartLine, Column: typeStartCol + len(paramType)},
		}
		node.AddChild(paramNode)

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}
}
