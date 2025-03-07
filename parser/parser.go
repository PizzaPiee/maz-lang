package parser

import (
	"maz-lang/ast"
	"maz-lang/lexer"
	"maz-lang/token"
	"slices"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUAL
	PLUS
	PRODUCT
	PREFIX
	PAREN
)

const (
	ErrUnexpectedParenthesis     = "unexpected parenthesis"
	ErrCannotParseToken          = "cannot parse current token"
	ErrExpectedIdentifier        = "expected next token to be an identifier"
	ErrExpectedAssignment        = "expected assignment"
	ErrMissingSemicolon          = "missing semicolon"
	ErrExpectedExpression        = "expected expression"
	ErrExpectedBlock             = "expected block"
	ErrExpectedParenthesis       = "expected parenthesis"
	ErrInvalidFunctionParameters = "function has invalid parameters"
)

var precedences = map[token.TokenType]int{
	token.PLUS:     PLUS,
	token.MINUS:    PLUS,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   PAREN,
	token.EQ:       EQUAL,
	token.NEQ:      EQUAL,
	token.GT:       EQUAL,
	token.LT:       EQUAL,
	token.GTEQ:     EQUAL,
	token.LTEQ:     EQUAL,
}

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
type InfixFn func(left ast.Node, endTokens ...token.TokenType) ast.Node

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
	p.registerPrefixFn(token.STRING, p.parseStringLiteral)
	p.registerPrefixFn(token.LPAREN, p.parseParenExpression)
	p.registerPrefixFn(token.LET, p.parseLetStatement)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.IF, p.parseIfStatement)
	p.registerPrefixFn(token.RETURN, p.parseReturnStatement)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionDefinition)

	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NEQ, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GTEQ, p.parseInfixExpression)
	p.registerInfixFn(token.LTEQ, p.parseInfixExpression)
	p.registerInfixFn(token.LPAREN, p.parseFunctionCall)

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

func (p *Parser) Parse(end token.TokenType) ast.Program {
	var program ast.Program

	for {
		tok := p.curToken
		if tok.Type == end || tok.Type == token.ILLEGAL {
			return program
		}

		node := p.parseExpression(LOWEST, end)
		if p.isError(node) {
			program.Statements = []ast.Node{node}
			return program
		}

		program.Statements = append(program.Statements, node)
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

func (p *Parser) parseExpression(precedence int, endTokens ...token.TokenType) ast.Node {
	tok := p.curToken
	prefixFn, ok := p.prefixFns[tok.Type]
	if !ok {
		return &ast.SyntaxError{Msg: ErrExpectedExpression, Token: p.curToken}
	}

	left := prefixFn()

	for precedence < p.peekPrecedence && !slices.Contains(endTokens, p.peekToken.Type) {
		infixFn, ok := p.infixFns[p.peekToken.Type]
		if !ok {
			return left
		}

		p.nextToken()
		left = infixFn(left, endTokens...)
	}

	return left
}

func (p *Parser) parsePrefixExpression() ast.Node {
	prefix := p.curToken
	p.nextToken()
	// NOTE: Well you see how parsePrefixExpression() is being called here? Like it? neither do I.
	// This is like this due to a bug I found and will stay like this until I find the will to think
	// of a better solution. I actually know what is the best way to fix this but this is just way
	// easier so we will stick with this.
	expression := p.parseExpression(PREFIX, token.EOF, token.SEMICOLON, token.COMMA, token.RPAREN)
	node := ast.PrefixExpression{Prefix: prefix, Value: expression}

	return &node
}

func (p *Parser) parseInfixExpression(left ast.Node, endTokens ...token.TokenType) ast.Node {
	node := ast.InfixExpression{Left: left, Operator: p.curToken}
	p.nextToken()
	node.Right = p.parseExpression(precedences[node.Operator.Type], endTokens...)

	if p.isError(node.Right) {
		return node.Right
	}

	return &node
}

func (p *Parser) parseParenExpression() ast.Node {
	p.nextToken()
	node := p.parseExpression(LOWEST, token.EOF)
	if !p.peekTokenIs(token.RPAREN) {
		return &ast.SyntaxError{Msg: ErrUnexpectedParenthesis, Token: p.curToken}
	}
	p.nextToken()

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

func (p *Parser) parseStringLiteral() ast.Node {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

// This function does not simply parse an Identifier.
// It acts as a gateway to the operations that can happen after an identifier.
// You can reference an identifier, assign a new value to it or maybe
// the identifier refers to the name of a function you are calling.
func (p *Parser) parseIdentifier() ast.Node {
	name := p.curToken.Literal
	// if p.peekTokenIs(token.LPAREN) {
	// 	p.nextToken()
	// 	args := p.parseArguments()
	// 	return &ast.FunctionCall{Name: name, Arguments: args}
	// }

	return &ast.Identifier{Name: name}
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

	exp := p.parseExpression(LOWEST, token.SEMICOLON)
	if p.isError(exp) {
		return exp
	}

	if !p.peekTokenIs(token.SEMICOLON) {
		return &ast.SyntaxError{Msg: ErrMissingSemicolon, Token: p.curToken}
	}

	p.nextToken()

	return &ast.LetStatement{Ident: ident, Value: exp}
}

func (p *Parser) parseIfStatement() ast.Node {
	node := ast.IfStatement{}

	// Parse main condition
	p.nextToken()
	condition := p.parseExpression(LOWEST, token.LBRACE)
	if p.isError(condition) {
		return condition
	}
	node.MainCondition = condition

	// Next token must be a '{'
	if !p.peekTokenIs(token.LBRACE) {
		return &ast.SyntaxError{Msg: ErrExpectedBlock, Token: p.curToken}
	}

	// Parse body of main condition
	p.nextToken()
	p.nextToken()
	stmts := p.Parse(token.RBRACE).Statements
	for _, stmt := range stmts {
		if p.isError(stmt) {
			return stmt
		}
	}
	node.MainStatements = stmts

	var elseIfs []ast.ElseIf
	for p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if p.peekTokenIs(token.IF) {
			// Parse else if condition
			elseIf := ast.ElseIf{}
			p.nextToken()
			p.nextToken()

			condition = p.parseExpression(LOWEST, token.LBRACE)
			if p.isError(condition) {
				return condition
			}
			elseIf.Condition = condition

			// Next token must be a '{'
			if !p.peekTokenIs(token.LBRACE) {
				return &ast.SyntaxError{Msg: ErrExpectedBlock, Token: p.curToken}
			}

			// Parse body of else if condition
			p.nextToken()
			p.nextToken()
			stmts = p.Parse(token.RBRACE).Statements
			for _, stmt := range stmts {
				if p.isError(stmt) {
					return stmt
				}
			}
			elseIf.Statements = stmts

			elseIfs = append(elseIfs, elseIf)
		} else {
			// Parse body of else condition
			// Next token must be a '{'
			if !p.peekTokenIs(token.LBRACE) {
				return &ast.SyntaxError{Msg: ErrExpectedBlock, Token: p.curToken}
			}
			p.nextToken()
			p.nextToken()
			stmts = p.Parse(token.RBRACE).Statements
			for _, stmt := range stmts {
				if p.isError(stmt) {
					return stmt
				}
			}
			node.ElseStatements = stmts
		}
	}
	node.ElseIfs = elseIfs

	return &node
}

func (p *Parser) parseReturnStatement() ast.Node {
	if p.peekTokenIs(token.SEMICOLON) {
		return &ast.SyntaxError{Msg: ErrExpectedExpression, Token: p.curToken}
	}

	p.nextToken()
	node := ast.ReturnStatement{}
	node.Expression = p.parseExpression(LOWEST, token.SEMICOLON)

	if !p.peekTokenIs(token.SEMICOLON) {
		return &ast.SyntaxError{Msg: ErrMissingSemicolon, Token: p.curToken}
	}
	p.nextToken()

	return &node
}

func (p *Parser) parseFunctionCall(left ast.Node, _ ...token.TokenType) ast.Node {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return &ast.SyntaxError{Msg: ErrExpectedIdentifier, Token: p.curToken}
	}

	args := p.parseArguments()
	return &ast.FunctionCall{Name: ident.Name, Arguments: args}
}

func (p *Parser) parseFunctionDefinition() ast.Node {
	node := ast.FunctionDefinition{}

	// Parse the function's name
	if !p.peekTokenIs(token.IDENT) {
		return &ast.SyntaxError{Msg: ErrExpectedIdentifier, Token: p.curToken}
	}
	p.nextToken()
	node.Name = p.curToken.Literal

	// Parse the parameters of the function
	if !p.peekTokenIs(token.LPAREN) {
		return &ast.SyntaxError{Msg: ErrExpectedParenthesis, Token: p.curToken}
	}
	p.nextToken()

	for !p.peekTokenIs(token.RPAREN) {
		if !p.peekTokenIs(token.IDENT) {
			return &ast.SyntaxError{Msg: ErrInvalidFunctionParameters, Token: p.curToken}
		}
		p.nextToken()

		param := p.parseIdentifier()

		switch param.(type) {
		case *ast.Identifier:
			node.Parameters = append(node.Parameters, param)
		default:
			return &ast.SyntaxError{Msg: ErrInvalidFunctionParameters, Token: p.curToken}
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken()

	// Parse the body of the function
	if !p.peekTokenIs(token.LBRACE) {
		return &ast.SyntaxError{Msg: ErrExpectedBlock, Token: p.curToken}
	}
	p.nextToken()
	p.nextToken()

	body := p.Parse(token.RBRACE)
	for _, stmt := range body.Statements {
		if p.isError(stmt) {
			return stmt
		}
	}
	node.Body = body.Statements

	return &node
}

func (p *Parser) parseArguments() []ast.Node {
	node := []ast.Node{}

	for !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		arg := p.parseExpression(LOWEST, token.COMMA, token.RPAREN)
		node = append(node, arg)
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken()

	return node
}
