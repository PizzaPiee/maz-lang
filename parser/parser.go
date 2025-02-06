package parser

import (
	"maz-lang/ast"
	"maz-lang/lexer"
	"maz-lang/token"
)

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	prefixFns map[token.TokenType]PrefixFn
}

type PrefixFn func() ast.Node

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, prefixFns: make(map[token.TokenType]PrefixFn)}

	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefixFn(key token.TokenType, fn PrefixFn) {
	p.prefixFns[key] = fn
}

func (p *Parser) Parse() ast.Program {
	var program ast.Program

	for {
		tok := p.curToken
		if tok.Type == token.EOF || tok.Type == token.ILLEGAL {
			return program
		}
		prefixFn := p.prefixFns[tok.Type]
		node := prefixFn()
		program.Statements = append(program.Statements, node)

		p.nextToken()
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parsePrefixExpression() ast.Node {
	prefix := p.curToken
	p.nextToken()
	value := p.curToken
	result := ast.PrefixExpression{Prefix: prefix, Value: value}

	return &result
}
