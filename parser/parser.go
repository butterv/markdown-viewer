package parser

import (
	"github.com/istsh/markdown-viewer/ast"
	"github.com/istsh/markdown-viewer/lexer"
	"github.com/istsh/markdown-viewer/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	//prefixParseFns map[token.TokenType]prefixParseFn
	//infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	return p.parseStringLiteral()
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	sl := &ast.StringLiteral{Token: p.curToken}
	return sl
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
