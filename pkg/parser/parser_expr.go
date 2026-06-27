package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

func (p *Parser) parseExpression() ASTNode {
	debugPrint("  [parseExpression] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// STOP - jangan parse jika token bukan bagian dari expression
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
		debugPrint("    [parseComparison] returning IS NULL/NOT NULL node\n")
		return nullNode
	}

	// Handle BETWEEN
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
		betweenNode.AddChild(node)
		betweenNode.AddChild(lower)
		betweenNode.AddChild(upper)

		debugPrint("    [parseComparison] returning BETWEEN node\n")
		return betweenNode
	}

	// Handle NOT ... (including NOT BETWEEN)
	if p.curToken.Type == token.NOT {
		debugPrint("    [parseComparison] found NOT, checking next token\n")
		tok := p.curToken

		// Save position
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

			// Create BETWEEN node
			betweenNode := NewASTNode(BetweenNode, "BETWEEN", tok.Line, tok.Column)
			betweenNode.AddChild(node)
			betweenNode.AddChild(lower)
			betweenNode.AddChild(upper)

			// Wrap with NOT
			notNode := NewASTNode(UnaryOpNode, "NOT", tok.Line, tok.Column)
			notNode.AddChild(betweenNode)

			debugPrint("    [parseComparison] returning NOT BETWEEN node\n")
			return notNode
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
				return unaryNode
			}
			return node
		}
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
		debugPrint("    [parseUnary] END (unary minus): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	// NOT is handled in parseComparison for NOT BETWEEN
	// But if we encounter NOT here, treat as unary NOT
	if p.curToken.Type == token.NOT {
		tok := p.curToken
		p.nextToken()
		operand := p.parseComparison()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		debugPrint("    [parseUnary] END (unary not): token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return unaryNode
	}

	result := p.parsePrimary()
	debugPrint("    [parseUnary] END (primary): token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)
	return result
}
