package parser

import (
	"maz-lang/ast"
	"maz-lang/lexer"
	"maz-lang/token"
	"strconv"
)

const (
	_ = iota
	LOWEST
	PLUS
	PRODUCT
	PREFIX
	PAREN
	IDENT
)

const (
	ErrUnexpectedParenthesis = "unexpected parenthesis"
	ErrCannotParseToken      = "cannot parse current token"
	ErrExpectedIdentifier    = "expected next token to be an identifier"
	ErrExpectedAssignment    = "expected assignment"
	ErrMissingSemicolon      = "missing semicolon"
	ErrExpectedExpression    = "expected expression"
)

var precedences = map[token.TokenType]int{
	token.PLUS:     PLUS,
	token.MINUS:    PLUS,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   PAREN,
	token.RPAREN:   PAREN,
	token.IDENT:    IDENT,
}

// This is a counter that keeps track of open parenthesis.
// When parsing parenthesis at the end of it the value of this counter
// must be equal to zero.
var openParen = 0

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	curPrecedence  int
	peekPrecedence int

	prefixFns map[token.TokenType]PrefixFn
	infixFns  map[token.TokenType]InfixFn
}

type PrefixFn func() ast.Node
type InfixFn func(left ast.Node) ast.Node

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:     lexer,
		prefixFns: make(map[token.TokenType]PrefixFn),
		infixFns:  make(map[token.TokenType]InfixFn),

		curPrecedence: LOWEST,
	}

	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefixFn(token.LPAREN, p.parseParenExpression)
	p.registerPrefixFn(token.LET, p.parseLetStatement)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)

	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefixFn(key token.TokenType, fn PrefixFn) {
	p.prefixFns[key] = fn
}

func (p *Parser) registerInfixFn(key token.TokenType, fn InfixFn) {
	p.infixFns[key] = fn
}

func (p *Parser) Parse() ast.Program {
	var program ast.Program

	for {
		tok := p.curToken
		if tok.Type == token.EOF || tok.Type == token.ILLEGAL {
			return program
		}

		node := p.parseExpression(LOWEST, token.EOF)
		program.Statements = append(program.Statements, node)

		switch node.(type) {
		case *ast.SyntaxError:
			return program
		}
		p.nextToken()
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()

	p.curPrecedence = precedences[p.curToken.Type]
	p.peekPrecedence = precedences[p.peekToken.Type]
}

func (p *Parser) peekTokenIs(token token.TokenType) bool {
	if p.peekToken.Type == token {
		return true
	}

	return false
}

func (p *Parser) isError(node ast.Node) bool {
	switch node.(type) {
	case *ast.SyntaxError:
		return true
	}

	return false
}

func (p *Parser) parseExpression(precedence int, end token.TokenType) ast.Node {
	tok := p.curToken
	prefixFn, ok := p.prefixFns[tok.Type]
	if !ok {
		return &ast.SyntaxError{Msg: ErrCannotParseToken, Token: p.curToken}
	}

	left := prefixFn()

	if p.peekTokenIs(end) {
		return left
	}

	for precedence < p.peekPrecedence {
		p.nextToken()
		infixFn, ok := p.infixFns[p.curToken.Type]
		if !ok {
			switch p.curToken.Type {
			case token.RPAREN:
				openParen--
			}
			break
		}

		left = infixFn(left)
	}

	if openParen != 0 {
		openParen = 0
		left = &ast.SyntaxError{Msg: ErrUnexpectedParenthesis, Token: p.curToken}
	}

	return left
}

func (p *Parser) parsePrefixExpression() ast.Node {
	prefix := p.curToken
	p.nextToken()
	expression := p.parseExpression(PREFIX, token.EOF)
	node := ast.PrefixExpression{Prefix: prefix, Value: expression}

	return &node
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	node := ast.InfixExpression{Left: left, Operator: p.curToken}
	p.nextToken()
	node.Right = p.parseExpression(precedences[node.Operator.Type], token.EOF)

	if p.isError(node.Right) {
		return node.Right
	}

	return &node
}

func (p *Parser) parseParenExpression() ast.Node {
	openParen++
	p.nextToken()
	node := p.parseExpression(LOWEST, token.EOF)

	return node
}

func (p *Parser) parseIntegerLiteral() ast.Node {
	num, _ := strconv.Atoi(p.curToken.Literal)
	node := ast.IntegerLiteral{Value: int64(num)}

	return &node
}

func (p *Parser) parseBooleanLiteral() ast.Node {
	if p.curToken.Type == token.TRUE {
		return &ast.BooleanLiteral{Value: true}
	}
	return &ast.BooleanLiteral{Value: false}
}

func (p *Parser) parseIdentifier() ast.Node {
	return &ast.Identifier{Name: p.curToken.Literal}
}

func (p *Parser) parseLetStatement() ast.Node {
	if !p.peekTokenIs(token.IDENT) {
		return &ast.SyntaxError{Msg: ErrExpectedIdentifier, Token: p.curToken}
	}

	p.nextToken()
	ident := p.curToken.Literal

	if !p.peekTokenIs(token.ASSIGN) {
		return &ast.SyntaxError{Msg: ErrExpectedAssignment, Token: p.curToken}
	}

	p.nextToken()

	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.EOF) {
		return &ast.SyntaxError{Msg: ErrExpectedExpression, Token: p.curToken}
	}

	p.nextToken()

	exp := p.parseExpression(LOWEST, token.EOF)
	if p.isError(exp) {
		return exp
	}

	if !p.peekTokenIs(token.SEMICOLON) {
		return &ast.SyntaxError{Msg: ErrMissingSemicolon, Token: p.curToken}
	}

	p.nextToken()

	return &ast.LetStatement{Ident: ident, Value: exp}
}
