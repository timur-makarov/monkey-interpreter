package evaluator

import (
	"github.com/timur-makarov/monkey-interpreter/internal/ast"
	"github.com/timur-makarov/monkey-interpreter/internal/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)
	case ast.BlockStatement:
		return evalBlock(n.Statements, env)
	case ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case ast.ReturnStatement:
		val := Eval(n.Value, env)
		if val.Type() == object.ErrorType {
			return val
		}
		return object.Return{Value: val}
	case ast.Integer:
		return object.Integer{Value: n.Value}
	case ast.String:
		return object.String{Value: n.Value}
	case ast.Boolean:
		return nativeBoolToObject(n.Value)
	case ast.Prefix:
		right := Eval(n.Right, env)
		if right.Type() == object.ErrorType {
			return right
		}
		return evalPrefix(n.Operator, right)
	case ast.Infix:
		left := Eval(n.Left, env)
		if left.Type() == object.ErrorType {
			return left
		}
		right := Eval(n.Right, env)
		if right.Type() == object.ErrorType {
			return right
		}
		return evalInfix(n.Operator, left, right, env)
	case ast.If:
		return evalIf(n, env)
	case ast.While:
		return evalWhile(n, env)
	case ast.LetStatement:
		val := Eval(n.Value, env)
		if val.Type() == object.ErrorType {
			return val
		}
		env.Set(n.Name.Value, val)
	case ast.Identifier:
		return evalIdentifier(n, env)
	case ast.Function:
		return object.Function{Parameters: n.Parameters, Env: env, Body: n.Body}
	case ast.Call:
		function := Eval(n.Function, env)
		if function.Type() == object.ErrorType {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && args[0].Type() == object.ErrorType {
			return args[0]
		}
		return evalFunction(function, args, env)
	case ast.Array:
		items := evalExpressions(n.Items, env)
		if len(items) == 1 && items[0].Type() == object.ErrorType {
			return items[0]
		}
		return object.Array{Items: items}
	case ast.AccessByExpression:
		left := Eval(n.Left, env)
		if left.Type() == object.ErrorType {
			return left
		}
		exp := Eval(n.Index, env)
		if exp.Type() == object.ErrorType {
			return left
		}
		return evalAccessByExpression(left, exp)
	case ast.HashTable:
		return evalHashTable(n, env)
	}

	return NULL
}
