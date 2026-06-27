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
