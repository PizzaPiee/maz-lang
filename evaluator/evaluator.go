package evaluator

import (
	"maz-lang/ast"
	"maz-lang/environment"
	"maz-lang/object"
)

var (
	TRUE  = object.Boolean{Value: true}
	FALSE = object.Boolean{Value: false}
	NULL  = object.Null{}
)

func Eval(node ast.Node, env *environment.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value == true {
			return &TRUE
		}
		return &FALSE
	case *ast.PrefixExpression:
		return evalPrefixExpression(*node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(*node, env)
	case *ast.LetStatement:
		return evalLetStatement(*node, env)
	case *ast.Identifier:
		return evalIdentifier(*node, env)
	}

	return nil
}

func evalStatements(statements []ast.Node, env *environment.Environment) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt, env)
	}

	return obj
}

func evalPrefixExpression(node ast.PrefixExpression, env *environment.Environment) object.Object {
	obj := Eval(node.Value, env)

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

func evalInfixExpression(node ast.InfixExpression, env *environment.Environment) object.Object {
	left := Eval(node.Left, env)
	right := Eval(node.Right, env)

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

func evalLetStatement(node ast.LetStatement, env *environment.Environment) object.Object {
	value := Eval(node.Value, env)
	env.Set(node.Ident, value)

	return &object.Boolean{Value: true}
}

func evalIdentifier(node ast.Identifier, env *environment.Environment) object.Object {
	res := env.Get(node.Name)
	if res != nil {
		return res
	}

	return &NULL
}
