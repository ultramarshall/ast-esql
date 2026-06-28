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

	// Handle IS NULL / IS NOT NULL
	if p.curToken.Type == token.ISNULL || p.curToken.Type == token.NOTNULL {
		debugPrint("    [parseComparison] found IS NULL/NOT NULL: %s\n", p.curToken.Literal)
		tok := p.curToken
		p.nextToken()

		var nullNode ASTNode
		if tok.Type == token.ISNULL {
			nullNode = NewASTNode(IsNullNode, "IS NULL", tok.Line, tok.Column)
		} else {
			nullNode = NewASTNode(IsNotNullNode, "IS NOT NULL", tok.Line, tok.Column)
		}
		nullNode.AddChild(node)
		nullNode.Span = combineSpan(node, nullNode)
		debugPrint("    [parseComparison] returning IS NULL/NOT NULL node\n")
		return nullNode
	}

	// Handle NOT ... (including NOT BETWEEN)
	if p.curToken.Type == token.NOT {
		debugPrint("    [parseComparison] found NOT, checking next token\n")
		tok := p.curToken
		pos := p.position
		p.nextToken() // consume NOT

		// Check if next token is BETWEEN
		if p.curToken.Type == token.BETWEEN {
			debugPrint("    [parseComparison] found NOT BETWEEN\n")
			p.nextToken() // consume BETWEEN

			lower := p.parseAdditive()
			if lower.Type == "" {
				p.errors = append(p.errors,
					fmt.Sprintf("expected lower bound in NOT BETWEEN expression at line %d", tok.Line))
				return node
			}

			if p.curToken.Type != token.AND {
				p.errors = append(p.errors,
					fmt.Sprintf("expected AND in NOT BETWEEN expression, got %s at line %d",
						p.curToken.Type, p.curToken.Line))
				return node
			}
			p.nextToken()

			upper := p.parseAdditive()
			if upper.Type == "" {
				p.errors = append(p.errors,
					fmt.Sprintf("expected upper bound in NOT BETWEEN expression at line %d", tok.Line))
				return node
			}

			// Buat BetweenNode dengan flag Not = true
			betweenNode := NewASTNode(BetweenNode, "BETWEEN", tok.Line, tok.Column)
			betweenNode.Not = true
			betweenNode.AddChild(node)
			betweenNode.AddChild(lower)
			betweenNode.AddChild(upper)
			betweenNode.Span = combineSpans(node, lower, upper)

			debugPrint("    [parseComparison] returning NOT BETWEEN node\n")
			return betweenNode
		} else {
			// Not NOT BETWEEN, rewind and handle as unary NOT
			debugPrint("    [parseComparison] NOT followed by %s, treating as unary NOT\n", p.curToken.Type)
			p.position = pos
			p.curToken = p.tokens[p.position]

			// Parse as unary NOT
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
			return node
		}
	}

	// Handle BETWEEN (regular, tanpa NOT)
	if p.curToken.Type == token.BETWEEN {
		debugPrint("    [parseComparison] found BETWEEN\n")
		tok := p.curToken
		p.nextToken()

		lower := p.parseAdditive()
		if lower.Type == "" {
			p.errors = append(p.errors,
				fmt.Sprintf("expected lower bound in BETWEEN expression at line %d", tok.Line))
			return node
		}

		if p.curToken.Type != token.AND {
			p.errors = append(p.errors,
				fmt.Sprintf("expected AND in BETWEEN expression, got %s at line %d",
					p.curToken.Type, p.curToken.Line))
			return node
		}
		p.nextToken()

		upper := p.parseAdditive()
		if upper.Type == "" {
			p.errors = append(p.errors,
				fmt.Sprintf("expected upper bound in BETWEEN expression at line %d", tok.Line))
			return node
		}

		betweenNode := NewASTNode(BetweenNode, tok.Literal, tok.Line, tok.Column)
		betweenNode.Not = false
		betweenNode.AddChild(node)
		betweenNode.AddChild(lower)
		betweenNode.AddChild(upper)
		betweenNode.Span = combineSpans(node, lower, upper)

		debugPrint("    [parseComparison] returning BETWEEN node\n")
		return betweenNode
	}

	// Handle regular comparison operators
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
			compNode.Span = combineSpan(node, right)
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
