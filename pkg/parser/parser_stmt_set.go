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
		// Gunakan Span untuk posisi
		targetWrapper := NewASTNode(BlockNode, "target", target.Span.Start.Line, target.Span.Start.Column)
		targetWrapper.AddChild(target)
		targetWrapper.Span = target.Span
		node.AddChild(targetWrapper)
	}

	// Parse '='
	if p.curToken.Type == token.ASSIGN || p.curToken.Type == token.EQ {
		p.nextToken()
		value := p.parseExpression()
		if value.Type != "" {
			// Gunakan Span untuk posisi
			valueWrapper := NewASTNode(BlockNode, "value", value.Span.Start.Line, value.Span.Start.Column)
			valueWrapper.AddChild(value)
			valueWrapper.Span = value.Span
			node.AddChild(valueWrapper)
			// Update node span to include value
			node.Span.End = value.Span.End
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
