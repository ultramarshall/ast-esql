package parser

import (
	"esql-ast-tool/internal/token"
)

func (p *Parser) parseCall() ASTNode {
	node := NewASTNode(CallNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse procedure name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse arguments
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			arg := p.parseExpression()
			if arg.Type != "" {
				node.AddChild(arg)
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}
