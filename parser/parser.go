package parser

import (
	"maz-lang/ast"
	"maz-lang/lexer"
	"maz-lang/token"
	"strconv"
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
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.FALSE, p.parseBooleanLiteral)

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

func (p *Parser) parseExpression() ast.Node {
	tok := p.curToken
	prefixFn := p.prefixFns[tok.Type]
	exp := prefixFn()

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Node {
	prefix := p.curToken
	p.nextToken()
	expression := p.parseExpression()
	node := ast.PrefixExpression{Prefix: prefix, Value: expression}

	return &node
}

func (p *Parser) parseIntegerLiteral() ast.Node {
	// FIX: handle the case where the literal value cannot be converted
	num, _ := strconv.Atoi(p.curToken.Literal)
	node := ast.IntegerLiteral{Value: int64(num)}

	return &node
}

func (p *Parser) parseBooleanLiteral() ast.Node {
	// FIX: use an else if and return error
	if p.curToken.Type == token.TRUE {
		return &ast.BooleanLiteral{Value: true}
	}
	return &ast.BooleanLiteral{Value: false}
}
