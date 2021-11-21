package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	//
	case *ast.Program:
		return evalStatemens(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	//
	case *ast.IntegerLiteral:
		return &object.Intger{Value: node.Value}
	}
	return nil
}

func evalStatemens(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}
