package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

// Helper: combine span dari left dan right
func combineSpan(left, right ASTNode) Span {
	return Span{
		Start: left.Span.Start,
		End:   right.Span.End,
	}
}

// Helper: combine span dari multiple nodes
func combineSpans(nodes ...ASTNode) Span {
	if len(nodes) == 0 {
		return Span{}
	}
	span := nodes[0].Span
	for _, n := range nodes[1:] {
		if n.Span.End.Line > span.End.Line ||
			(n.Span.End.Line == span.End.Line && n.Span.End.Column > span.End.Column) {
			span.End = n.Span.End
		}
	}
	return span
}

func (p *Parser) parseExpression() ASTNode {
	debugPrint("  [parseExpression] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.END || p.curToken.Type == token.THEN ||
		p.curToken.Type == token.ELSE || p.curToken.Type == token.WHEN ||
		p.curToken.Type == token.EOF {
		debugPrint("  [parseExpression] STOP: token is %s, not an expression\n", p.curToken.Type)
		return ASTNode{}
	}

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
			compNode.Span = combineSpan(left, right)
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
			binOp.Span = combineSpan(node, right)
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
			binOp.Span = combineSpan(node, right)
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

	// Cek operator-operator yang mungkin muncul setelah node
	return p.parseComparisonSuffix(node)
}

// parseComparisonSuffix menangani operator setelah node kiri
func (p *Parser) parseComparisonSuffix(left ASTNode) ASTNode {
	switch p.curToken.Type {
	case token.ISNULL, token.NOTNULL:
		return p.parseIsNull(left)

	case token.NOT:
		return p.parseNotOperator(left)

	case token.BETWEEN:
		return p.parseBetween(left, false)

	case token.LIKE:
		return p.parseLike(left, false)

	case token.EQ, token.NOT_EQ, token.LT, token.GT, token.LTE, token.GTE:
		return p.parseComparisonOperator(left)

	default:
		debugPrint("    [parseComparison] END: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return left
	}
}

// parseIsNull menangani IS NULL / IS NOT NULL
func (p *Parser) parseIsNull(left ASTNode) ASTNode {
	debugPrint("    [parseIsNull] found IS NULL/NOT NULL: %s\n", p.curToken.Literal)
	tok := p.curToken
	p.nextToken()

	var nullNode ASTNode
	if tok.Type == token.ISNULL {
		nullNode = NewASTNode(IsNullNode, "IS NULL", tok.Line, tok.Column)
	} else {
		nullNode = NewASTNode(IsNotNullNode, "IS NOT NULL", tok.Line, tok.Column)
	}
	nullNode.AddChild(left)
	nullNode.Span = combineSpan(left, nullNode)
	debugPrint("    [parseIsNull] returning IS NULL/NOT NULL node\n")
	return nullNode
}

// parseNotOperator menangani NOT (termasuk NOT BETWEEN, NOT LIKE)
func (p *Parser) parseNotOperator(left ASTNode) ASTNode {
	debugPrint("    [parseNotOperator] found NOT, checking next token\n")
	tok := p.curToken
	pos := p.position
	p.nextToken() // consume NOT

	switch p.curToken.Type {
	case token.BETWEEN:
		debugPrint("    [parseNotOperator] found NOT BETWEEN\n")
		return p.parseBetween(left, true)

	case token.LIKE:
		debugPrint("    [parseNotOperator] found NOT LIKE\n")
		return p.parseLike(left, true)

	default:
		// Not NOT BETWEEN/LIKE, treat as unary NOT
		debugPrint("    [parseNotOperator] NOT followed by %s, treating as unary NOT\n", p.curToken.Type)
		p.position = pos
		p.curToken = p.tokens[p.position]

		tok = p.curToken
		p.nextToken()
		right := p.parseComparison()
		if right.Type != "" {
			unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
			unaryNode.AddChild(right)
			unaryNode.Span = combineSpan(
				NewASTNode(IdentifierNode, tok.Literal, tok.Line, tok.Column),
				right,
			)
			return unaryNode
		}
		return left
	}
}

// parseBetween menangani BETWEEN / NOT BETWEEN
func (p *Parser) parseBetween(left ASTNode, isNot bool) ASTNode {
	debugPrint("    [parseBetween] parsing BETWEEN (not=%v)\n", isNot)
	tok := p.curToken
	if isNot {
		p.nextToken() // consume BETWEEN (sudah di-consume di parseNotOperator)
	} else {
		p.nextToken() // consume BETWEEN
	}

	lower := p.parseAdditive()
	if lower.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected lower bound in BETWEEN expression at line %d", tok.Line))
		return left
	}

	if p.curToken.Type != token.AND {
		p.errors = append(p.errors,
			fmt.Sprintf("expected AND in BETWEEN expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return left
	}
	p.nextToken()

	upper := p.parseAdditive()
	if upper.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected upper bound in BETWEEN expression at line %d", tok.Line))
		return left
	}

	betweenNode := NewASTNode(BetweenNode, tok.Literal, tok.Line, tok.Column)
	betweenNode.Not = isNot
	betweenNode.AddChild(left)
	betweenNode.AddChild(lower)
	betweenNode.AddChild(upper)
	betweenNode.Span = combineSpans(left, lower, upper)

	debugPrint("    [parseBetween] returning BETWEEN node (not=%v)\n", isNot)
	return betweenNode
}

// parseLike menangani LIKE / NOT LIKE
func (p *Parser) parseLike(left ASTNode, isNot bool) ASTNode {
	debugPrint("    [parseLike] parsing LIKE (not=%v)\n", isNot)
	tok := p.curToken
	if isNot {
		p.nextToken() // consume LIKE (sudah di-consume di parseNotOperator)
	} else {
		p.nextToken() // consume LIKE
	}

	pattern := p.parseAdditive()
	if pattern.Type == "" {
		p.errors = append(p.errors,
			fmt.Sprintf("expected pattern in LIKE expression at line %d", tok.Line))
		return left
	}

	likeNode := NewASTNode(LikeNode, tok.Literal, tok.Line, tok.Column)
	likeNode.Not = isNot
	likeNode.AddChild(left)
	likeNode.AddChild(pattern)
	likeNode.Span = combineSpans(left, pattern)

	debugPrint("    [parseLike] returning LIKE node (not=%v)\n", isNot)
	return likeNode
}

// parseComparisonOperator menangani operator comparison biasa (=, <, >, dll)
func (p *Parser) parseComparisonOperator(left ASTNode) ASTNode {
	debugPrint("    [parseComparisonOperator] found operator: %s\n", p.curToken.Literal)
	tok := p.curToken
	p.nextToken()
	right := p.parseAdditive()
	if right.Type != "" {
		compNode := NewASTNode(ComparisonNode, tok.Literal, tok.Line, tok.Column)
		compNode.AddChild(left)
		compNode.AddChild(right)
		compNode.Span = combineSpan(left, right)
		debugPrint("    [parseComparisonOperator] returning comparison node\n")
		return compNode
	}
	return left
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
		binOp.Span = combineSpan(node, right)
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
		binOp.Span = combineSpan(node, right)
		node = binOp
	}

	debugPrint("    [parseMultiplicative] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseUnary() ASTNode {
	debugPrint("    [parseUnary] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.MINUS {
		tok := p.curToken
		p.nextToken()
		operand := p.parsePrimary()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		unaryNode.Span = combineSpan(
			NewASTNode(IdentifierNode, tok.Literal, tok.Line, tok.Column),
			operand,
		)
		debugPrint("    [parseUnary] END (unary minus): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	// Note: NOT is now handled in parseComparison, but keep this for safety
	if p.curToken.Type == token.NOT {
		tok := p.curToken
		p.nextToken()
		operand := p.parseComparison()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		unaryNode.Span = combineSpan(
			NewASTNode(IdentifierNode, tok.Literal, tok.Line, tok.Column),
			operand,
		)
		debugPrint("    [parseUnary] END (unary not): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	result := p.parsePrimary()
	debugPrint("    [parseUnary] END (primary): token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
	return result
}
