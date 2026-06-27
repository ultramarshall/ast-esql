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
