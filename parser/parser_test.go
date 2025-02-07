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
