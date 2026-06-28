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
		targetWrapper := NewASTNode(BlockNode, "target", target.Span.Start.Line, target.Span.Start.Column)
		targetWrapper.AddChild(target)
		targetWrapper.Span = target.Span
		node.AddChild(targetWrapper)
	}

	// Parse '='
	if p.curToken.Type == token.ASSIGN || p.curToken.Type == token.EQ {
		p.nextToken()

		// ✅ Cek apakah ada expression setelah '='
		if p.curToken.Type == token.SEMICOLON || p.curToken.Type == token.EOF {
			p.errors = append(p.errors,
				fmt.Sprintf("expected expression after '=' in SET statement at line %d",
					p.curToken.Line))
			// Consume semicolon or EOF to avoid infinite loop
			if p.curToken.Type == token.SEMICOLON {
				p.nextToken()
			}
			return node
		}

		// ✅ Cek juga jika token adalah ';' setelah whitespace
		// (ini sudah di-handle di atas)

		value := p.parseAssignmentRHS()
		if value.Type != "" {
			valueWrapper := NewASTNode(BlockNode, "value", value.Span.Start.Line, value.Span.Start.Column)
			valueWrapper.AddChild(value)
			valueWrapper.Span = value.Span
			node.AddChild(valueWrapper)
			node.Span.End = value.Span.End
		} else {
			// ✅ Kalau parseAssignmentRHS return kosong, tambahkan error
			p.errors = append(p.errors,
				fmt.Sprintf("invalid expression after '=' in SET statement at line %d",
					p.curToken.Line))
			return node
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

func (p *Parser) parseAssignmentRHS() ASTNode {
	// Parse sebagai additive dulu
	left := p.parseAdditive()

	debugPrint("  [parseAssignmentRHS] after parseAdditive: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Jika setelah parseAdditive masih ada operator comparison selain '='
	// Misalnya: >, <, >=, <=, !=
	if p.curToken.Type == token.GT || p.curToken.Type == token.LT ||
		p.curToken.Type == token.GTE || p.curToken.Type == token.LTE ||
		p.curToken.Type == token.NOT_EQ {
		// Ini adalah comparison, parse sebagai comparison
		debugPrint("  [parseAssignmentRHS] found comparison operator, parsing as comparison\n")
		return p.parseComparisonSuffix(left)
	}

	// ✅ Jika token saat ini adalah ';', berarti selesai
	if p.curToken.Type == token.SEMICOLON {
		debugPrint("  [parseAssignmentRHS] found ';', returning left\n")
		return left
	}

	// ✅ Jika token saat ini adalah '=', abaikan (sudah di-consume di parseSet)
	if p.curToken.Type == token.EQ || p.curToken.Type == token.ASSIGN {
		debugPrint("  [parseAssignmentRHS] found '=', but in SET context, skipping...\n")
		// Consume token '=' agar tidak menyebabkan infinite loop
		p.nextToken()
		return left
	}

	debugPrint("  [parseAssignmentRHS] returning left (no more operators)\n")
	return left
}
