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
