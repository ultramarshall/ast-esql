package parser

import (
	"fmt"
	"strconv"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parsePrimary() ASTNode {
	debugPrint("    [parsePrimary] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	switch p.curToken.Type {
	case token.IDENTIFIER:
		result := p.parseIdentifier()
		debugPrint("    [parsePrimary] after IDENTIFIER: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.NUMBER:
		result := p.parseNumber()
		debugPrint("    [parsePrimary] after NUMBER: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.STRING:
		result := p.parseString()
		debugPrint("    [parsePrimary] after STRING: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.LPAREN:
		result := p.parseGroupedExpression()
		debugPrint("    [parsePrimary] after LPAREN: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.CASE:
		result := p.parseCase()
		debugPrint("    [parsePrimary] after CASE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.CAST:
		result := p.parseCast()
		debugPrint("    [parsePrimary] after CAST: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return result
	case token.DOT:
		debugPrint("    [parsePrimary] WARNING: DOT without identifier, skipping...\n")
		p.nextToken()
		return ASTNode{}
	default:
		debugPrint("    [parsePrimary] UNKNOWN: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
		return ASTNode{}
	}
}

func (p *Parser) parseIdentifier() ASTNode {
	debugPrint("      [parseIdentifier] token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type == token.LPAREN {
		return p.parseFunctionCall(node)
	}

	if p.curToken.Type == token.DOT {
		return p.parseFieldReference(node)
	}

	debugPrint("      [parseIdentifier] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseFunctionCall(name ASTNode) ASTNode {
	funcName := name.Token
	if name.Value != nil {
		if str, ok := name.Value.(string); ok {
			funcName = str
		}
	}
	node := NewASTNode(FunctionCallNode, funcName, name.Line, name.Column)
	node.Value = funcName
	p.nextToken()

	if p.curToken.Type != token.RPAREN {
		arg := p.parseExpression()
		if arg.Type != "" {
			node.AddChild(arg)
		}

		for p.curToken.Type == token.COMMA {
			p.nextToken()
			arg = p.parseExpression()
			if arg.Type != "" {
				node.AddChild(arg)
			}
		}
	}

	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFieldReference(base ASTNode) ASTNode {
	fieldNode := NewASTNode(FieldReferenceNode, "field", base.Line, base.Column)
	fieldNode.AddChild(base)

	if base.Value != nil {
		fieldNode.Value = base.Value
	} else {
		fieldNode.Value = base.Token
	}

	for p.curToken.Type == token.DOT {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			identNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			identNode.Value = p.curToken.Literal

			newFieldNode := NewASTNode(FieldReferenceNode, "field", fieldNode.Line, fieldNode.Column)
			newFieldNode.AddChild(fieldNode)
			newFieldNode.AddChild(identNode)

			if fieldNode.Value != nil {
				newFieldNode.Value = fieldNode.Value.(string) + "." + p.curToken.Literal
			} else {
				newFieldNode.Value = p.curToken.Literal
			}

			fieldNode = newFieldNode
			p.nextToken()
		}
	}

	return fieldNode
}

func (p *Parser) parseNumber() ASTNode {
	val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
	node := NewASTNode(LiteralNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = val
	p.nextToken()
	return node
}

func (p *Parser) parseString() ASTNode {
	node := NewASTNode(LiteralNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()
	return node
}

func (p *Parser) parseGroupedExpression() ASTNode {
	p.nextToken()
	expr := p.parseExpression()
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}
	return expr
}

func (p *Parser) parseCast() ASTNode {
	debugPrint("    [parseCast] START: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	node := NewASTNode(CastNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = "CAST"
	p.nextToken()

	debugPrint("    [parseCast] after CAST: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.LPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected '(' after CAST, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] after '(': token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	expr := p.parseExpression()
	if expr.Type != "" {
		node.AddChild(expr)
	}

	debugPrint("    [parseCast] after expression: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.AS {
		p.errors = append(p.errors,
			fmt.Sprintf("expected 'AS' in CAST expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] after AS: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors,
			fmt.Sprintf("expected type name after AS in CAST, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}

	typeNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	typeNode.Value = p.curToken.Literal
	node.AddChild(typeNode)
	p.nextToken()

	debugPrint("    [parseCast] after type: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type != token.RPAREN {
		p.errors = append(p.errors,
			fmt.Sprintf("expected ')' in CAST expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}
	p.nextToken()

	debugPrint("    [parseCast] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}

func (p *Parser) parseCase() ASTNode {
	debugPrint("[parseCase] START: token=%s, literal='%s', line=%d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(CaseNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("[parseCase] after CASE: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	var isSimpleCase bool
	var caseExpr ASTNode

	if p.curToken.Type != token.WHEN {
		isSimpleCase = true
		debugPrint("[parseCase] Simple CASE, parsing expression\n")
		caseExpr = p.parseExpression()
		if caseExpr.Type != "" {
			node.AddChild(caseExpr)
		}
		debugPrint("[parseCase] after expression: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	} else {
		debugPrint("[parseCase] Searched CASE\n")
	}

	whenCount := 0
	for p.curToken.Type == token.WHEN {
		whenCount++
		debugPrint("[parseCase] Parsing WHEN #%d at line %d\n", whenCount, p.curToken.Line)
		whenNode := p.parseWhen(isSimpleCase)
		if whenNode.Type != "" {
			node.AddChild(whenNode)
		}
		debugPrint("[parseCase] after WHEN #%d: token=%s, literal='%s'\n",
			whenCount, p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.ELSE {
		debugPrint("[parseCase] Parsing ELSE\n")
		p.nextToken()
		elseExpr := p.parseExpression()
		if elseExpr.Type != "" {
			elseNode := NewASTNode(BlockNode, "else", elseExpr.Line, elseExpr.Column)
			elseNode.AddChild(elseExpr)
			node.AddChild(elseNode)
		}
		debugPrint("[parseCase] after ELSE: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.END {
		debugPrint("[parseCase] Found END, consuming it\n")
		p.nextToken()
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected END in CASE expression, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
	}

	debugPrint("[parseCase] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	if p.curToken.Type == token.THEN || p.curToken.Type == token.ELSE || p.curToken.Type == token.END {
		debugPrint("[parseCase] Next token is %s, returning early\n", p.curToken.Type)
		return node
	}

	return node
}

func (p *Parser) parseWhen(isSimpleCase bool) ASTNode {
	debugPrint("  [parseWhen] START: token=%s, literal='%s', isSimple=%v\n",
		p.curToken.Type, p.curToken.Literal, isSimpleCase)

	node := NewASTNode(WhenNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	debugPrint("  [parseWhen] after WHEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	var condition ASTNode

	if isSimpleCase {
		debugPrint("  [parseWhen] Simple CASE: parsing value\n")
		condition = p.parseExpression()
		if condition.Type != "" {
			node.AddChild(condition)
		}
		debugPrint("  [parseWhen] after value: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	} else {
		debugPrint("  [parseWhen] Searched CASE: parsing condition\n")
		condition = p.parseExpression()
		if condition.Type != "" {
			node.AddChild(condition)
		}
		debugPrint("  [parseWhen] after condition: token=%s, literal='%s'\n",
			p.curToken.Type, p.curToken.Literal)
	}

	if p.curToken.Type == token.THEN {
		debugPrint("  [parseWhen] Found THEN\n")
		p.nextToken()
	} else {
		p.errors = append(p.errors,
			fmt.Sprintf("expected THEN in CASE WHEN, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		return node
	}

	debugPrint("  [parseWhen] after THEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	result := p.parseExpression()
	if result.Type != "" {
		node.AddChild(result)
	}

	debugPrint("  [parseWhen] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}
