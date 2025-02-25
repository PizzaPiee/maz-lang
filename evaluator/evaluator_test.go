package evaluator

import (
	"maz-lang/environment"
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
		env := environment.New()
		obj := Eval(&program, &env)

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
		t.Logf("evaluating '%s'\n", tt.Expression)
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		env := environment.New()
		obj := Eval(&program, &env)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj object.Object
	}{
		{
			Expression:  "!true",
			ExpectedObj: &object.Boolean{Value: false},
		},
		{
			Expression:  "!false",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "-10",
			ExpectedObj: &object.Integer{Value: -10},
		},
	}

	for _, tt := range tests {
		t.Logf("evaluating '%s'\n", tt.Expression)
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		env := environment.New()
		obj := Eval(&program, &env)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj object.Object
	}{
		{
			Expression:  "2+1*5",
			ExpectedObj: &object.Integer{Value: 7},
		},
		{
			Expression:  "(2*(1+5))-10",
			ExpectedObj: &object.Integer{Value: 2},
		},
		{
			Expression:  "((100+5)-(3*5))/5",
			ExpectedObj: &object.Integer{Value: 18},
		},
		{
			Expression:  "5 >= 5",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "1 > 2",
			ExpectedObj: &object.Boolean{Value: false},
		},
		{
			Expression:  "10 == 10",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "10 != 10",
			ExpectedObj: &object.Boolean{Value: false},
		},
		{
			Expression:  "10 < 20",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "10 <= 10",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "10 + 20 * 2",
			ExpectedObj: &object.Integer{Value: 50},
		},
		{
			Expression:  "(10 + 20) * 2",
			ExpectedObj: &object.Integer{Value: 60},
		},
		{
			Expression:  "10 / 2 + 5",
			ExpectedObj: &object.Integer{Value: 10},
		},
		{
			Expression:  "true == true",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "true != false",
			ExpectedObj: &object.Boolean{Value: true},
		},
		{
			Expression:  "2 * (3 + 4) - 5",
			ExpectedObj: &object.Integer{Value: 9},
		},
		{
			Expression:  "2 * 3 + 4 * 5",
			ExpectedObj: &object.Integer{Value: 26},
		},
		{
			Expression:  "2 * (3 + 4 * 5)",
			ExpectedObj: &object.Integer{Value: 46},
		},
		{
			Expression:  "2 * 3 + 4 * 5 - 6 / 2",
			ExpectedObj: &object.Integer{Value: 23},
		},
		{
			Expression:  "2 * (3 + 4) - 6 / 2",
			ExpectedObj: &object.Integer{Value: 11},
		},
		{
			Expression:  "2 * 3 + 4 * (5 - 6) / 2",
			ExpectedObj: &object.Integer{Value: 4},
		},
		{
			Expression:  "2 * (3 + 4 * (5 - 6)) / 2",
			ExpectedObj: &object.Integer{Value: -1},
		},
		{
			Expression:  "2 * 3 + 4 * 5 - 6 / 2 + 10",
			ExpectedObj: &object.Integer{Value: 33},
		},
		{
			Expression:  "2 * (3 + 4) - 6 / 2 + 10",
			ExpectedObj: &object.Integer{Value: 21},
		},
		{
			Expression:  "2 * 3 + 4 * (5 - 6) / 2 + 10",
			ExpectedObj: &object.Integer{Value: 14},
		},
		{
			Expression:  "2 * (3 + 4 * (5 - 6)) / 2 + 10",
			ExpectedObj: &object.Integer{Value: 9},
		},
	}

	for _, tt := range tests {
		t.Logf("evaluating '%s'\n", tt.Expression)
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		env := environment.New()
		obj := Eval(&program, &env)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}

func TestEvalLetStatement(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj object.Object
	}{
		{
			Expression:  "let a = 10; a",
			ExpectedObj: &object.Integer{Value: 10},
		},
		{
			Expression:  "let a = !true; a",
			ExpectedObj: &object.Boolean{Value: false},
		},
	}

	for _, tt := range tests {
		t.Logf("evaluating: '%s'\n", tt.Expression)
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		env := environment.New()
		obj := Eval(&program, &env)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj object.Object
	}{
		{
			Expression:  "if 5 > 2 {10} else {20}",
			ExpectedObj: &object.Integer{Value: 10},
		},
		{
			Expression:  "let a = 5; let b = 7; if a > b {10} else {20}",
			ExpectedObj: &object.Integer{Value: 20},
		},
		{
			Expression:  "let a = true; let b = false; if a == b {10} else if a != b {20} else {30}",
			ExpectedObj: &object.Integer{Value: 20},
		},
		{
			Expression:  "if 1 > 2 {let a = 10;} else {let a = 20;} a",
			ExpectedObj: &object.Null{},
		},
	}

	for _, tt := range tests {
		t.Logf("evaluating: '%s'\n", tt.Expression)
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		env := environment.New()
		obj := Eval(&program, &env)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}

func TestEvalFunction(t *testing.T) {
	tests := []struct {
		Expression  string
		ExpectedObj object.Object
	}{
		{
			Expression:  "fn sum(a, b) {return a + b;} sum(1,2)",
			ExpectedObj: &object.Return{Value: &object.Integer{Value: 3}},
		},
		{
			Expression:  "fn fib(n) { if n == 0 { return 0; } else if n == 1 { return 1; } else if n == 2 { return 1; } else { return fib(n-1) + fib(n-2); } } fib(19)",
			ExpectedObj: &object.Return{Value: &object.Integer{Value: 4181}},
		},
	}

	for _, tt := range tests {
		t.Logf("evaluating: '%s'\n", tt.Expression)
		l := lexer.New(tt.Expression)
		program := parser.New(&l).Parse(token.EOF)
		env := environment.New()
		obj := Eval(&program, &env)

		if !cmp.Equal(obj, tt.ExpectedObj) {
			t.Errorf("expected object to be %+v, instead got %+v\n", tt.ExpectedObj, obj)
		}
	}
}
