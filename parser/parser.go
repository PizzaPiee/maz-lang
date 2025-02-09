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

var precedences = map[token.TokenType]int{
	token.PLUS:     PLUS,
	token.MINUS:    PLUS,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   PAREN,
	token.RPAREN:   PAREN,
	token.IDENT:    IDENT,
}

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

		node := p.parseExpression(LOWEST)
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

func (p *Parser) parseExpression(precedence int) ast.Node {
	tok := p.curToken
	prefixFn, ok := p.prefixFns[tok.Type]
	if !ok {
		return &ast.SyntaxError{Msg: "cannot parse current token", Token: p.curToken}
	}

	left := prefixFn()
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

	return left
}

func (p *Parser) parsePrefixExpression() ast.Node {
	prefix := p.curToken
	p.nextToken()
	expression := p.parseExpression(PREFIX)
	node := ast.PrefixExpression{Prefix: prefix, Value: expression}

	return &node
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	node := ast.InfixExpression{Left: left, Operator: p.curToken}
	p.nextToken()
	node.Right = p.parseExpression(precedences[node.Operator.Type])

	return &node
}

func (p *Parser) parseParenExpression() ast.Node {
	openParen++
	p.nextToken()
	node := p.parseExpression(LOWEST)

	if p.curToken.Type == token.RPAREN {
		openParen--
	}

	if openParen != 0 {
		return &ast.SyntaxError{Msg: "unexpected parenthesis", Token: p.curToken}
	}

	openParen = 0

	p.nextToken()
	return node
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
