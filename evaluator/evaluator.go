package evaluator

import (
	"fmt"
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
	case *ast.SyntaxError:
		return &object.Error{Value: node}
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
	case *ast.IfStatement:
		return evalIfStatement(*node, env)
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
	case ">":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Boolean{Value: left.(*object.Integer).Value > right.(*object.Integer).Value}
		}
	case ">=":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Boolean{Value: left.(*object.Integer).Value >= right.(*object.Integer).Value}
		}

	case "<":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Boolean{Value: left.(*object.Integer).Value < right.(*object.Integer).Value}
		}
	case "<=":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Boolean{Value: left.(*object.Integer).Value <= right.(*object.Integer).Value}
		}
	case "==":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Boolean{Value: left.(*object.Integer).Value == right.(*object.Integer).Value}
		}
	case "!=":
		if (left.Type() == right.Type()) && left.Type() == object.INTEGER_OBJ {
			return &object.Boolean{Value: left.(*object.Integer).Value != right.(*object.Integer).Value}
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

func evalIfStatement(node ast.IfStatement, env *environment.Environment) object.Object {
	mainCondition := Eval(node.MainCondition, env)

	switch mainCondition := mainCondition.(type) {
	case *object.Boolean:
		if mainCondition.Value {
			currentEnv := environment.New()
			currentEnv.Extend(env)
			return evalStatements(node.MainStatements, &currentEnv)
		}
	default:
		return &object.Error{Value: fmt.Errorf("expected boolean, instead got '%s'\n", mainCondition.Inspect())}
	}

	for _, elseIf := range node.ElseIfs {
		currentEnv := environment.New()
		currentEnv.Extend(env)
		res := evalElseIf(elseIf, &currentEnv)
		if res != nil {
			return res
		}
	}

	if len(node.ElseStatements) != 0 {
		currentEnv := environment.New()
		currentEnv.Extend(env)
		return evalStatements(node.ElseStatements, &currentEnv)
	}

	return &NULL
}

func evalElseIf(node ast.ElseIf, env *environment.Environment) object.Object {
	condition := Eval(node.Condition, env)

	switch condition := condition.(type) {
	case *object.Boolean:
		if condition.Value {
			return evalStatements(node.Statements, env)
		}
	default:
		return &object.Error{Value: fmt.Errorf("expected boolean, instead got '%s'\n", condition.Inspect())}
	}

	return nil
}
