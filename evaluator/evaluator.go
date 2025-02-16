package evaluator

import (
	"maz-lang/ast"
	"maz-lang/object"
)

var (
	TRUE  = object.Boolean{Value: true}
	FALSE = object.Boolean{Value: false}
	NULL  = object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return EvalStatements(node.Statements)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value == true {
			return &TRUE
		}
		return &FALSE
	case *ast.PrefixExpression:
		return EvalPrefixExpression(*node)

	}

	return nil
}

func EvalStatements(statements []ast.Node) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt)
	}

	return obj
}

func EvalPrefixExpression(node ast.PrefixExpression) object.Object {
	obj := Eval(node.Value)

	switch node.Prefix.Literal {
	case "!":
		switch obj := obj.(type) {
		case *object.Boolean:
			return &object.Boolean{Value: !obj.Value}
		}
	case "-":
		switch obj := obj.(type) {
		case *object.Integer:
			return &object.Integer{Value: -obj.Value}
		}
	}

	return nil
}
