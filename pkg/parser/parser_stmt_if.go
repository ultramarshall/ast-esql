package parser

import (
	"fmt"

	"esql-ast-tool/internal/token"
)

func (p *Parser) parseIf() ASTNode {
	debugPrint("[parseIf] START: token=%s, literal='%s', line=%d\n",
		p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(IfNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken() // consume IF

	debugPrint("[parseIf] after IF: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Parse condition - parseExpression will handle NOT
	cond := p.parseExpression()
	if cond.Type != "" {
		debugPrint("[parseIf] condition parsed: type=%s\n", cond.Type)
		node.AddChild(cond)
	} else {
		debugPrint("[parseIf] WARNING: empty condition\n")
	}

	debugPrint("[parseIf] after condition: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Expect THEN
	if p.curToken.Type == token.THEN {
		debugPrint("[parseIf] Found THEN, consuming it\n")
		p.nextToken() // consume THEN
	} else {
		debugPrint("[parseIf] ERROR: expected THEN, got %s\n", p.curToken.Type)
		p.errors = append(p.errors,
			fmt.Sprintf("expected THEN after IF condition, got %s at line %d",
				p.curToken.Type, p.curToken.Line))
		// Advance to prevent infinite loop
		if p.curToken.Type != token.EOF {
			p.nextToken()
		}
		return node
	}

	debugPrint("[parseIf] after THEN: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	// Parse THEN block
	thenBlock := NewASTNode(BlockNode, "then", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.ELSE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			thenBlock.AddChild(stmt)
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
	}
	node.AddChild(thenBlock)

	// Parse ELSE if ada
	if p.curToken.Type == token.ELSE {
		debugPrint("[parseIf] Found ELSE\n")
		p.nextToken()
		elseBlock := NewASTNode(BlockNode, "else", p.curToken.Line, p.curToken.Column)
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				elseBlock.AddChild(stmt)
			}
			if p.curToken.Type == token.SEMICOLON {
				p.nextToken()
			}
		}
		node.AddChild(elseBlock)
	}

	// Konsumsi END IF
	if p.curToken.Type == token.END {
		debugPrint("[parseIf] Found END\n")
		p.nextToken()
		if p.curToken.Type == token.IF {
			debugPrint("[parseIf] Found IF after END\n")
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	debugPrint("[parseIf] END: token=%s, literal='%s'\n",
		p.curToken.Type, p.curToken.Literal)

	return node
}
