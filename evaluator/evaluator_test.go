package evaluator

import (
	"maz-lang/lexer"
	"maz-lang/object"
	"maz-lang/parser"
	"maz-lang/token"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj *object.Integer
	}{
		{
			Expression:  "5",
			ExpectedObj: &object.Integer{Value: 5},
		},
		{
			Expression:  "100",
			ExpectedObj: &object.Integer{Value: 100},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		obj := Eval(&program)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}

func TestEvalBooleanLiteral(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj *object.Boolean
	}{
		{
			Expression:  "true",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "false",
			ExpectedObj: &object.Boolean{Value: false},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		obj := Eval(&program)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}
