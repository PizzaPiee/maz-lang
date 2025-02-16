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
		return evalStatements(node.Statements)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value == true {
			return &TRUE
		}
		return &FALSE
	case *ast.PrefixExpression:
		return evalPrefixExpression(*node)
	case *ast.InfixExpression:
		return evalInfixExpression(*node)
	}

	return nil
}

func evalStatements(statements []ast.Node) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt)
	}

	return obj
}

func evalPrefixExpression(node ast.PrefixExpression) object.Object {
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

func evalInfixExpression(node ast.InfixExpression) object.Object {
	left := Eval(node.Left)
	right := Eval(node.Right)

	switch node.Operator.Literal {
	case "+":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Integer{
				Value: left.(*object.Integer).Value + right.(*object.Integer).Value,
			}
		}
	case "-":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Integer{
				Value: left.(*object.Integer).Value - right.(*object.Integer).Value,
			}
		}
	case "*":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Integer{
				Value: left.(*object.Integer).Value * right.(*object.Integer).Value,
			}
		}
	case "/":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Integer{
				Value: left.(*object.Integer).Value / right.(*object.Integer).Value,
			}
		}
	}

	return nil
}
