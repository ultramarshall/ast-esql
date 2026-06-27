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
