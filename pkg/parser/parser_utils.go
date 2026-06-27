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
