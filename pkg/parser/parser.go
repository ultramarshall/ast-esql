package parser

import (
	"esql-ast-tool/internal/token"
	"fmt"
)

var DebugMode = false

func debugPrint(format string, args ...interface{}) {
	if DebugMode {
		fmt.Printf(format, args...)
	}
}

type Parser struct {
	l         *Lexer
	tokens    []token.Token
	position  int
	curToken  token.Token
	peekToken token.Token
	errors    []string
	inSet     bool
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

func (p *Parser) GetCurToken() token.Token {
	return p.curToken
}

func (p *Parser) GetPeekToken() token.Token {
	return p.peekToken
}

func (p *Parser) GetNextToken() {
	p.nextToken()
}

func (p *Parser) GetPosition() int {
	return p.position
}

func (p *Parser) DebugTokens() {
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

func (p *Parser) GetTokens() []token.Token {
	return p.tokens
}

func (p *Parser) DebugPrintTokens() {
	fmt.Println("=== TOKENS ===")
	for i, tok := range p.tokens {
		fmt.Printf("[%d] Type: %s, Literal: '%s', Line: %d, Col: %d\n",
			i, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// ParseProgram - entry point
func (p *Parser) ParseProgram() Program {
	program := NewProgram()
	rootNode := NewASTNode(ProgramNode, "PROGRAM", 1, 1)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt.Type != "" {
			rootNode.AddChild(stmt)
			program.AddStatement(stmt) // <-- tambahkan ini
		}

		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		} else if p.curToken.Type != token.EOF && p.curToken.Type != token.END {
			p.nextToken()
		}
	}

	program.Root = rootNode
	return program
}
