package parser

import (
	"fmt"
	"strconv"

	"esql-ast-tool/internal/token"
)

type Parser struct {
	l         *Lexer
	tokens    []token.Token
	position  int
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(input string) *Parser {
	l := NewLexer(input)
	tokens := l.Tokenize()
	p := &Parser{
		l:        l,
		tokens:   tokens,
		position: 0,
		errors:   []string{},
	}

	if len(tokens) > 0 {
		p.curToken = tokens[0]
		if len(tokens) > 1 {
			p.peekToken = tokens[1]
		}
	}

	return p
}

func (p *Parser) nextToken() {
	p.position++
	if p.position < len(p.tokens) {
		p.curToken = p.tokens[p.position]
		if p.position+1 < len(p.tokens) {
			p.peekToken = p.tokens[p.position+1]
		} else {
			p.peekToken = token.Token{Type: token.EOF}
		}
	} else {
		p.curToken = token.Token{Type: token.EOF}
		p.peekToken = token.Token{Type: token.EOF}
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) DebugTokens() {
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// GetTokens returns all tokens for debugging
func (p *Parser) GetTokens() []token.Token {
	return p.tokens
}

// DebugPrintTokens prints all tokens
func (p *Parser) DebugPrintTokens() {
	fmt.Println("=== TOKENS ===")
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// Parse program
func (p *Parser) ParseProgram() Program {
	program := NewProgram()

	for p.curToken.Type != token.EOF {
		// Debug
		// fmt.Printf("DEBUG ParseProgram: token=%s, literal='%s'\n", p.curToken.Type, p.curToken.Literal)

		// Skip END tokens - tapi hati-hati jangan skip token yang penting
		if p.curToken.Type == token.END {
			// Cek apakah ini END IF atau END MODULE
			nextToken := p.peekToken
			if nextToken.Type == token.IF {
				p.nextToken() // consume END
				p.nextToken() // consume IF
				if p.curToken.Type == token.SEMICOLON {
					p.nextToken()
				}
				continue
			}
			if nextToken.Type == token.MODULE {
				p.nextToken() // consume END
				p.nextToken() // consume MODULE
				if p.curToken.Type == token.SEMICOLON {
					p.nextToken()
				}
				continue
			}
			// END lainnya, skip
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt.Type != "" {
			program.AddStatement(stmt)
		}

		// Konsumsi semicolon jika ada
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.EOF && p.curToken.Type != token.END {
			// Safety advance - tapi hati-hati
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseStatement() ASTNode {
	// Debug - matikan setelah selesai debugging
	// fmt.Printf("DEBUG parseStatement: token=%s, literal='%s', line=%d\n",
	//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	switch p.curToken.Type {
	case token.CREATE:
		return p.parseCreate()
	case token.DECLARE:
		return p.parseDeclare()
	case token.SET:
		return p.parseSet()
	case token.IF:
		return p.parseIf()
	case token.WHILE:
		return p.parseWhile()
	case token.FOR:
		return p.parseFor()
	case token.RETURN:
		return p.parseReturn()
	case token.THROW:
		return p.parseThrow()
	case token.PROPAGATE:
		return p.parsePropagate()
	case token.MOVE:
		return p.parseMove()
	case token.CONTINUE:
		return p.parseContinue()
	case token.BREAK:
		return p.parseBreak()
	case token.LABEL:
		return p.parseLabel()
	case token.MODULE:
		return p.parseModuleStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.PROCEDURE:
		return p.parseProcedureStatement()
	case token.CALL:
		return p.parseCall()
	case token.END:
		// Skip END tokens
		p.nextToken()
		return ASTNode{}
	default:
		// Parse sebagai expression
		expr := p.parseExpression()
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return expr
	}
}

// Parse CALL statement
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

// Parse MODULE statement
func (p *Parser) parseModuleStatement() ASTNode {
	node := NewASTNode(ModuleNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse module name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		nameNode.Value = p.curToken.Literal
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse body statements - jangan skip apapun
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		// Debug
		// fmt.Printf("DEBUG parseModuleStatement: token=%s, literal='%s', line=%d\n",
		//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

		stmt := p.parseStatement()
		if stmt.Type != "" {
			node.AddChild(stmt)
		}

		// Jangan panggil p.nextToken() di sini karena parseStatement sudah handle
		// Hanya konsumsi semicolon jika ada dan parser belum lanjut
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			// Hanya advance jika benar-benar perlu
			// Tapi hati-hati, ini bisa menyebabkan skip token
			// p.nextToken()
		}
	}

	// Konsumsi END MODULE
	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()
		}
	}

	// Optional semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// Parse FUNCTION statement
func (p *Parser) parseFunctionStatement() ASTNode {
	node := NewASTNode(FunctionNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse function name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		if nameNode.Value == nil {
			nameNode.Value = nameNode.Token
		}
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse parameters
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			if p.curToken.Type == token.IDENTIFIER {
				paramNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				if paramNode.Value == nil {
					paramNode.Value = paramNode.Token
				}
				node.AddChild(paramNode)
				p.nextToken()
				// Skip type
				if p.curToken.Type == token.IDENTIFIER {
					p.nextToken()
				}
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	// Parse return type
	if p.curToken.Type == token.RETURNS {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			returnType := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			if returnType.Value == nil {
				returnType.Value = returnType.Token
			}
			node.AddChild(returnType)
			p.nextToken()
		}
	}

	// Parse function body
	if p.curToken.Type == token.BEGIN {
		p.nextToken()
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				node.AddChild(stmt)
			}
			if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.END {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

// Parse PROCEDURE statement
func (p *Parser) parseProcedureStatement() ASTNode {
	node := NewASTNode(ProcedureNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse procedure name
	if p.curToken.Type == token.IDENTIFIER {
		nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		if nameNode.Value == nil {
			nameNode.Value = nameNode.Token
		}
		node.AddChild(nameNode)
		p.nextToken()
	}

	// Parse parameters
	if p.curToken.Type == token.LPAREN {
		p.nextToken()
		for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
			if p.curToken.Type == token.IDENTIFIER {
				paramNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				if paramNode.Value == nil {
					paramNode.Value = paramNode.Token
				}
				node.AddChild(paramNode)
				p.nextToken()
				// Skip type
				if p.curToken.Type == token.IDENTIFIER {
					p.nextToken()
				}
			}
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.RPAREN {
			p.nextToken()
		}
	}

	// Parse procedure body
	if p.curToken.Type == token.BEGIN {
		p.nextToken()
		for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			stmt := p.parseStatement()
			if stmt.Type != "" {
				node.AddChild(stmt)
			}
			if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
				p.nextToken()
			}
		}
		if p.curToken.Type == token.END {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseDeclare() ASTNode {
	node := NewASTNode(DeclareNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse variable name
	if p.curToken.Type != token.IDENTIFIER {
		p.errors = append(p.errors, "expected identifier after DECLARE")
		return node
	}

	nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	nameNode.Value = p.curToken.Literal
	node.AddChild(nameNode)
	p.nextToken()

	// Parse type (INTEGER, STRING, etc.)
	if p.curToken.Type == token.IDENTIFIER {
		typeNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		typeNode.Value = p.curToken.Literal
		node.AddChild(typeNode)
		p.nextToken()
	}

	// Optional DEFAULT
	if p.curToken.Type == token.DEFAULT {
		p.nextToken()
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseSet() ASTNode {
	// Debug
	// fmt.Printf("DEBUG parseSet: token=%s, literal='%s', line=%d\n",
	//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

	node := NewASTNode(SetNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse target - ini bisa identifier atau field reference
	target := p.parseExpression()
	if target.Type != "" {
		targetWrapper := NewASTNode(BlockNode, "target", target.Line, target.Column)
		targetWrapper.AddChild(target)
		node.AddChild(targetWrapper)
	}

	// Debug
	// fmt.Printf("DEBUG parseSet: after target, token=%s, literal='%s'\n",
	//     p.curToken.Type, p.curToken.Literal)

	// Cek ASSIGN
	if p.curToken.Type == token.ASSIGN {
		// fmt.Println("DEBUG parseSet: found ASSIGN")
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
	}

	// Consume semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseIf() ASTNode {
	node := NewASTNode(IfNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	// Parse condition
	cond := p.parseExpression()
	if cond.Type != "" {
		condWrapper := NewASTNode(BlockNode, "condition", cond.Line, cond.Column)
		condWrapper.AddChild(cond)
		node.AddChild(condWrapper)
	}

	// Expect THEN
	if p.curToken.Type == token.THEN {
		p.nextToken()
	}

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
		p.nextToken()
		if p.curToken.Type == token.IF {
			p.nextToken()
		}
	}

	// Optional semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseWhile() ASTNode {
	node := NewASTNode(WhileNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	condition := p.parseExpression()
	if condition.Type != "" {
		node.AddChild(condition)
	}

	if p.curToken.Type == token.DO {
		p.nextToken()
	}

	bodyNode := NewASTNode(BlockNode, "body", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			bodyNode.AddChild(stmt)
		}
		if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			p.nextToken()
		}
	}
	node.AddChild(bodyNode)

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.WHILE {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseFor() ASTNode {
	node := NewASTNode(ForNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER {
		varNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		varNode.Value = varNode.Token
		node.AddChild(varNode)
		p.nextToken()
	}

	if p.curToken.Type == token.AS {
		p.nextToken()
	}

	for p.curToken.Type != token.DO && p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.DO {
		p.nextToken()
	}

	bodyNode := NewASTNode(BlockNode, "body", p.curToken.Line, p.curToken.Column)
	for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			bodyNode.AddChild(stmt)
		}
		if p.curToken.Type != token.END && p.curToken.Type != token.EOF {
			p.nextToken()
		}
	}
	node.AddChild(bodyNode)

	if p.curToken.Type == token.END {
		p.nextToken()
		if p.curToken.Type == token.FOR {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseReturn() ASTNode {
	node := NewASTNode(ReturnNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseThrow() ASTNode {
	node := NewASTNode(ThrowNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type != token.SEMICOLON {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parsePropagate() ASTNode {
	node := NewASTNode(PropagateNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		expr := p.parseExpression()
		if expr.Type != "" {
			node.AddChild(expr)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseBreak() ASTNode {
	node := NewASTNode(BreakNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseMove() ASTNode {
	node := NewASTNode(MoveNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	target := p.parseExpression()
	if target.Type != "" {
		node.AddChild(target)
	}

	if p.curToken.Type == token.TO {
		p.nextToken()
		source := p.parseExpression()
		if source.Type != "" {
			node.AddChild(source)
		}
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseContinue() ASTNode {
	node := NewASTNode(ContinueNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	if p.curToken.Type == token.IDENTIFIER { // fix for continue label perhaps
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return node
}

func (p *Parser) parseLabel() ASTNode {
	node := NewASTNode(LabelNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.IDENTIFIER {
		node.Value = p.curToken.Literal
		p.nextToken()
	}

	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseCreate() ASTNode {
	node := NewASTNode(CreateNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	p.nextToken()

	if p.curToken.Type == token.COMPUTE {
		p.nextToken()
		if p.curToken.Type == token.MODULE {
			p.nextToken()

			// Create Module node
			moduleNode := NewASTNode(ModuleNode, "COMPUTE MODULE", p.curToken.Line, p.curToken.Column)
			moduleNode.Value = "COMPUTE"

			// Parse module name
			if p.curToken.Type == token.IDENTIFIER {
				nameNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
				nameNode.Value = p.curToken.Literal
				moduleNode.AddChild(nameNode)
				p.nextToken()
			}

			// Parse body statements
			for p.curToken.Type != token.END && p.curToken.Type != token.EOF {
				// Debug
				// fmt.Printf("DEBUG parseCreate: token=%s, literal='%s', line=%d\n",
				//     p.curToken.Type, p.curToken.Literal, p.curToken.Line)

				stmt := p.parseStatement()
				if stmt.Type != "" {
					moduleNode.AddChild(stmt)
				}

				// Hanya konsumsi semicolon jika perlu
				if p.curToken.Type == token.SEMICOLON {
					p.nextToken()
				}
			}

			// Konsumsi END MODULE
			if p.curToken.Type == token.END {
				p.nextToken()
				if p.curToken.Type == token.MODULE {
					p.nextToken()
				}
			}

			node.AddChild(moduleNode)
		}
	}

	// Optional semicolon
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return node
}

func (p *Parser) parseExpressionStatement() ASTNode {
	expr := p.parseExpression()
	if p.curToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return expr
}

func (p *Parser) parseExpression() ASTNode {
	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() ASTNode {
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
	node := p.parseAdditive()

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

func (p *Parser) parseAdditive() ASTNode {
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

	return node
}

func (p *Parser) parseMultiplicative() ASTNode {
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

	return node
}

func (p *Parser) parseUnary() ASTNode {
	if p.curToken.Type == token.MINUS || p.curToken.Type == token.NOT {
		tok := p.curToken
		p.nextToken()
		operand := p.parsePrimary()
		unaryNode := NewASTNode(UnaryOpNode, tok.Literal, tok.Line, tok.Column)
		unaryNode.AddChild(operand)
		return unaryNode
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ASTNode {
	switch p.curToken.Type {
	case token.IDENTIFIER:
		return p.parseIdentifier()
	case token.NUMBER:
		return p.parseNumber()
	case token.STRING:
		return p.parseString()
	case token.LPAREN:
		return p.parseGroupedExpression()
	default:
		// Jika token tidak dikenal, jangan buat node error
		// Biarkan parser melanjutkan
		return ASTNode{}
	}
}

func (p *Parser) parseIdentifier() ASTNode {
	node := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	node.Value = p.curToken.Literal
	p.nextToken()

	// Function call
	if p.curToken.Type == token.LPAREN {
		return p.parseFunctionCall(node)
	}

	// Field reference - Environment.Variables.Status
	if p.curToken.Type == token.DOT {
		return p.parseFieldReference(node)
	}

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
	// Base adalah root (Environment)
	fieldNode := NewASTNode(FieldReferenceNode, "field", base.Line, base.Column)
	fieldNode.AddChild(base)

	// Set initial value
	if base.Value != nil {
		fieldNode.Value = base.Value
	} else {
		fieldNode.Value = base.Token
	}

	// Lanjutkan parsing untuk rest of chain
	for p.curToken.Type == token.DOT {
		p.nextToken()
		if p.curToken.Type == token.IDENTIFIER {
			identNode := NewASTNode(IdentifierNode, p.curToken.Literal, p.curToken.Line, p.curToken.Column)
			identNode.Value = p.curToken.Literal

			// Buat FieldReference baru yang membungkus yang lama
			newFieldNode := NewASTNode(FieldReferenceNode, "field", fieldNode.Line, fieldNode.Column)
			newFieldNode.AddChild(fieldNode)
			newFieldNode.AddChild(identNode)

			// Update value dengan path lengkap
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
