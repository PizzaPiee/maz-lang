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
