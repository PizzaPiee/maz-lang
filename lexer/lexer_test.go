package lexer

import (
	"maz-lang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	+-*/=;,(){}
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
	2 >= 1
	3 <= 1
	`

	tests := []struct {
		ExpectedType    token.TokenType
		ExpectedLiteral string
	}{
		{ExpectedType: token.PLUS, ExpectedLiteral: "+"},
		{ExpectedType: token.MINUS, ExpectedLiteral: "-"},
		{ExpectedType: token.ASTERISK, ExpectedLiteral: "*"},
		{ExpectedType: token.SLASH, ExpectedLiteral: "/"},
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
		{ExpectedType: token.LET, ExpectedLiteral: "let"},
		{ExpectedType: token.IDENT, ExpectedLiteral: "a"},
		{ExpectedType: token.ASSIGN, ExpectedLiteral: "="},
		{ExpectedType: token.INT, ExpectedLiteral: "10"},
		{ExpectedType: token.SEMICOLON, ExpectedLiteral: ";"},
		{ExpectedType: token.LET, ExpectedLiteral: "let"},
		{ExpectedType: token.IDENT, ExpectedLiteral: "b"},
		{ExpectedType: token.ASSIGN, ExpectedLiteral: "="},
		{ExpectedType: token.FUNCTION, ExpectedLiteral: "fn"},
		{ExpectedType: token.LPAREN, ExpectedLiteral: "("},
		{ExpectedType: token.IDENT, ExpectedLiteral: "x"},
		{ExpectedType: token.RPAREN, ExpectedLiteral: ")"},
		{ExpectedType: token.LBRACE, ExpectedLiteral: "{"},
		{ExpectedType: token.IDENT, ExpectedLiteral: "x"},
		{ExpectedType: token.SEMICOLON, ExpectedLiteral: ";"},
		{ExpectedType: token.RBRACE, ExpectedLiteral: "}"},
		{ExpectedType: token.SEMICOLON, ExpectedLiteral: ";"},
		{ExpectedType: token.IF, ExpectedLiteral: "if"},
		{ExpectedType: token.INT, ExpectedLiteral: "5"},
		{ExpectedType: token.GT, ExpectedLiteral: ">"},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
		{ExpectedType: token.LBRACE, ExpectedLiteral: "{"},
		{ExpectedType: token.RETURN, ExpectedLiteral: "return"},
		{ExpectedType: token.INT, ExpectedLiteral: "10"},
		{ExpectedType: token.SEMICOLON, ExpectedLiteral: ";"},
		{ExpectedType: token.RBRACE, ExpectedLiteral: "}"},
		{ExpectedType: token.ELSE, ExpectedLiteral: "else"},
		{ExpectedType: token.LBRACE, ExpectedLiteral: "{"},
		{ExpectedType: token.RETURN, ExpectedLiteral: "return"},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
		{ExpectedType: token.SEMICOLON, ExpectedLiteral: ";"},
		{ExpectedType: token.RBRACE, ExpectedLiteral: "}"},
		{ExpectedType: token.INT, ExpectedLiteral: "10"},
		{ExpectedType: token.EQ, ExpectedLiteral: "=="},
		{ExpectedType: token.INT, ExpectedLiteral: "10"},
		{ExpectedType: token.INT, ExpectedLiteral: "5"},
		{ExpectedType: token.NEQ, ExpectedLiteral: "!="},
		{ExpectedType: token.INT, ExpectedLiteral: "2"},
		{ExpectedType: token.INT, ExpectedLiteral: "2"},
		{ExpectedType: token.LT, ExpectedLiteral: "<"},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
		{ExpectedType: token.INT, ExpectedLiteral: "3"},
		{ExpectedType: token.GT, ExpectedLiteral: ">"},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
		{ExpectedType: token.INT, ExpectedLiteral: "2"},
		{ExpectedType: token.GTEQ, ExpectedLiteral: ">="},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
		{ExpectedType: token.INT, ExpectedLiteral: "3"},
		{ExpectedType: token.LTEQ, ExpectedLiteral: "<="},
		{ExpectedType: token.INT, ExpectedLiteral: "1"},
	}

	l := New(input)

	for i, tt := range tests {
		token := l.NextToken()

		if token.Type != tt.ExpectedType {
			t.Errorf("#%d invalid token type, expected='%s' got='%s'", i, tt.ExpectedType, token.Type)
		}

		if token.Literal != tt.ExpectedLiteral {
			t.Errorf("#%d invalid token literal, expected='%s' got='%s'", i, tt.ExpectedLiteral, token.Literal)
		}
	}

}
