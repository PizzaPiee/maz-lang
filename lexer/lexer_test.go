package lexer

import (
	"maz-lang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	+=;,(){}
	10 1
	a foo fizz_buzz
	let a = 10;
	let b = fn(x) {x;};
	if 5 > 1 {
		return 10;
	} else {
		return 1;
	}
	10 == 10
	5 != 2
	2 < 1
	3 > 1
	`

	tests := []struct {
		ExpectedType    token.TokenType
		ExpectedLiteral string
	}{
		{ExpectedType: token.PLUS, ExpectedLiteral: "+"},
		{ExpectedType: token.ASSIGN, ExpectedLiteral: "="},
		{ExpectedType: token.SEMICOLON, ExpectedLiteral: ";"},
		{ExpectedType: token.COMMA, ExpectedLiteral: ","},
		{ExpectedType: token.LPAREN, ExpectedLiteral: "("},
		{ExpectedType: token.RPAREN, ExpectedLiteral: ")"},
		{ExpectedType: token.LBRACE, ExpectedLiteral: "{"},
		{ExpectedType: token.RBRACE, ExpectedLiteral: "}"},
		{ExpectedType: token.INT, ExpectedLiteral: "10"},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
		{ExpectedType: token.IDENT, ExpectedLiteral: "a"},
		{ExpectedType: token.IDENT, ExpectedLiteral: "foo"},
		{ExpectedType: token.IDENT, ExpectedLiteral: "fizz_buzz"},
	}

	l := New(input)

	for _, tt := range tests {
		token := l.NextToken()

		if token.Type != tt.ExpectedType {
			t.Errorf("invalid token type, expected=%s got=%s", tt.ExpectedType, token.Type)
		}

		if token.Literal != tt.ExpectedLiteral {
			t.Errorf("invalid token literal, expected=%s got=%s", tt.ExpectedLiteral, token.Literal)
		}
	}

}
