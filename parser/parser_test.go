package parser

import (
	"maz-lang/ast"
	"maz-lang/lexer"
	"maz-lang/token"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParsePrefixExpression(t *testing.T) {
	tests := []struct {
		Expression   string
		ExpectedNode ast.PrefixExpression
	}{
		{
			Expression: "-5",
			ExpectedNode: ast.PrefixExpression{
				Prefix: token.Token{Type: token.MINUS, Literal: "-"},
				Value:  &ast.IntegerLiteral{Value: int64(5)},
			},
		},
		{
			Expression: "!true",
			ExpectedNode: ast.PrefixExpression{
				Prefix: token.Token{Type: token.BANG, Literal: "!"},
				Value:  &ast.BooleanLiteral{Value: true},
			},
		},
		{
			Expression: "!false",
			ExpectedNode: ast.PrefixExpression{
				Prefix: token.Token{Type: token.BANG, Literal: "!"},
				Value:  &ast.BooleanLiteral{Value: false},
			},
		},
		{
			Expression: "!foo",
			ExpectedNode: ast.PrefixExpression{
				Prefix: token.Token{Type: token.BANG, Literal: "!"},
				Value:  &ast.Identifier{Name: "foo"},
			},
		},
		{
			Expression: "-foo",
			ExpectedNode: ast.PrefixExpression{
				Prefix: token.Token{Type: token.MINUS, Literal: "-"},
				Value:  &ast.Identifier{Name: "foo"},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Expression)
		p := New(&l)
		program := p.Parse()

		for _, stmt := range program.Statements {
			pe, ok := stmt.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("expected node of type ast.PrefixExpression, got=%T\n", pe)
			}

			if !cmp.Equal(pe.Prefix, tt.ExpectedNode.Prefix) {
				t.Errorf("expected prefix token to be %+v, instead got %+v\n", tt.ExpectedNode.Prefix, pe.Prefix)
			}

			if !cmp.Equal(pe.Value, tt.ExpectedNode.Value) {
				t.Errorf("expected value token to be %+v, instead got %+v\n", tt.ExpectedNode.Value, pe.Value)
			}
		}
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		Expression   string
		ExpectedNode ast.Node
	}{
		{
			Expression: "5+1*2",
			ExpectedNode: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 5},
				Operator: token.Token{Type: token.PLUS, Literal: "+"},
				Right: &ast.InfixExpression{
					Left:     &ast.IntegerLiteral{Value: 1},
					Operator: token.Token{Type: token.ASTERISK, Literal: "*"},
					Right:    &ast.IntegerLiteral{Value: 2},
				},
			},
		},
		{
			Expression: "(5+1)*2",
			ExpectedNode: &ast.InfixExpression{
				Left: &ast.InfixExpression{
					Left:     &ast.IntegerLiteral{Value: 5},
					Operator: token.Token{Type: token.PLUS, Literal: "+"},
					Right:    &ast.IntegerLiteral{Value: 1},
				},
				Operator: token.Token{Type: token.ASTERISK, Literal: "*"},
				Right:    &ast.IntegerLiteral{Value: 2},
			},
		},
		{
			Expression: "foo+bar",
			ExpectedNode: &ast.InfixExpression{
				Left:     &ast.Identifier{Name: "foo"},
				Operator: token.Token{Type: token.PLUS, Literal: "+"},
				Right:    &ast.Identifier{Name: "bar"},
			},
		},
		{
			Expression: "foo*bar+1",
			ExpectedNode: &ast.InfixExpression{
				Left: &ast.InfixExpression{
					Left:     &ast.Identifier{Name: "foo"},
					Operator: token.Token{Type: token.ASTERISK, Literal: "*"},
					Right:    &ast.Identifier{Name: "bar"},
				},
				Operator: token.Token{Type: token.PLUS, Literal: "+"},
				Right:    &ast.IntegerLiteral{Value: 1},
			},
		},
		{
			Expression: "(5+1",
			ExpectedNode: &ast.SyntaxError{
				Msg:   ErrUnexpectedParenthesis,
				Token: token.Token{Type: token.INT, Literal: "1"},
			},
		},
		{
			Expression: "5+1)",
			ExpectedNode: &ast.SyntaxError{
				Msg:   ErrUnexpectedParenthesis,
				Token: token.Token{Type: token.RPAREN, Literal: ")"},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Expression)
		p := New(&l)
		program := p.Parse()

		if !cmp.Equal(program.Statements[0], tt.ExpectedNode) {
			t.Errorf("expected node: %+v, instead got: %+v\n", tt.ExpectedNode, program.Statements[0])
		}
	}
}

func TestParseBlockExpressions(t *testing.T) {
	tests := []struct {
		Expression   string
		ExpectedNode ast.Node
		End          token.TokenType
	}{
		{
			Expression: "{5+1} let a = 1+2;",
			ExpectedNode: &ast.InfixExpression{
				Left:     &ast.IntegerLiteral{Value: 5},
				Operator: token.Token{Type: token.PLUS, Literal: "+"},
				Right:    &ast.IntegerLiteral{Value: 1},
			},
			End: token.RBRACE,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Expression)
		p := New(&l)
		p.nextToken() // Skip first token
		node := p.parseExpression(LOWEST, tt.End)

		if !cmp.Equal(node, tt.ExpectedNode) {
			t.Errorf("expected node: %+v, instead got: %+v\n", tt.ExpectedNode, node)
		}
	}
}

func TestParseLetStatement(t *testing.T) {
	tests := []struct {
		Expression   string
		ExpectedNode ast.Node
	}{
		{
			Expression: "let a = 5+1;",
			ExpectedNode: &ast.LetStatement{
				Ident: "a",
				Value: &ast.InfixExpression{
					Left:     &ast.IntegerLiteral{Value: 5},
					Operator: token.Token{Type: token.PLUS, Literal: "+"},
					Right:    &ast.IntegerLiteral{Value: 1},
				},
			},
		},
		{
			Expression: "let foo = (2+5)*2;",
			ExpectedNode: &ast.LetStatement{
				Ident: "foo",
				Value: &ast.InfixExpression{
					Left: &ast.InfixExpression{
						Left:     &ast.IntegerLiteral{Value: 2},
						Operator: token.Token{Type: token.PLUS, Literal: "+"},
						Right:    &ast.IntegerLiteral{Value: 5},
					},
					Operator: token.Token{Type: token.ASTERISK, Literal: "*"},
					Right:    &ast.IntegerLiteral{Value: 2},
				},
			},
		},
		{
			Expression: "let foo = bar + 1;",
			ExpectedNode: &ast.LetStatement{
				Ident: "foo",
				Value: &ast.InfixExpression{
					Left:     &ast.Identifier{Name: "bar"},
					Operator: token.Token{Type: token.PLUS, Literal: "+"},
					Right:    &ast.IntegerLiteral{Value: 1},
				},
			},
		},
		{
			Expression: "let 0 = 5+1;",
			ExpectedNode: &ast.SyntaxError{
				Msg:   ErrExpectedIdentifier,
				Token: token.Token{Type: token.LET, Literal: "let"},
			},
		},
		{
			Expression: "let a",
			ExpectedNode: &ast.SyntaxError{
				Msg:   ErrExpectedAssignment,
				Token: token.Token{Type: token.IDENT, Literal: "a"},
			},
		},
		{
			Expression: "let a = 1",
			ExpectedNode: &ast.SyntaxError{
				Msg:   ErrMissingSemicolon,
				Token: token.Token{Type: token.INT, Literal: "1"},
			},
		},
		{
			Expression: "let a = ;",
			ExpectedNode: &ast.SyntaxError{
				Msg:   ErrExpectedExpression,
				Token: token.Token{Type: token.ASSIGN, Literal: "="},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Expression)
		p := New(&l)
		program := p.Parse()

		if !cmp.Equal(program.Statements[0], tt.ExpectedNode) {
			t.Errorf("expected node: %+v, instead got %+v\n", tt.ExpectedNode, program.Statements[0])
		}
	}
}
