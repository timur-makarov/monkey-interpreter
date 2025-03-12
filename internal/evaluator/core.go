package evaluator

import (
	"github.com/timur-makarov/monkey-interpreter/internal/ast"
	"github.com/timur-makarov/monkey-interpreter/internal/object"
)

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch res := result.(type) {
		case object.Return:
			return res.Value
		case object.Error:
			return res
		case object.Identifier:
			result = res.Value
		case object.AccessByExpression:
			result = res.Value
		}
	}

	if result == nil {
		return NULL
	}

	return result
}

func evalBlock(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	enclosedEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range statements {
		result = Eval(statement, enclosedEnv)

		rt := result.Type()
		if rt == object.ReturnType || rt == object.ErrorType {
			return result
		}
	}

	if result == nil {
		return NULL
	}

	return result
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range expressions {
		evaluated := Eval(exp, env)
		if evaluated.Type() == object.ErrorType {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfix(operator string, left, right object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalIntegerInfixOperators(operator, left.(object.Integer), right.(object.Integer))
	case left.Type() == object.StringType && right.Type() == object.StringType:
		return evalStringInfixOperators(operator, left.(object.String), right.(object.String))
	case left.Type() == object.IdentifierType && right.Type() == object.IdentifierType:
		return evalInfix(
			operator, left.(object.Identifier).Value, right.(object.Identifier).Value, env,
		)
	case right.Type() == object.IdentifierType:
		return evalInfix(operator, left, right.(object.Identifier).Value, env)
	case right.Type() == object.AccessByExpressionType:
		return evalInfix(operator, left, right.(object.AccessByExpression).Value, env)
	case operator == "==":
		return nativeBoolToObject(left == right)
	case operator == "!=":
		return nativeBoolToObject(left != right)
	case operator == "=":
		if left.Type() == object.IdentifierType {
			env.Set(left.(object.Identifier).Name, right)
			return NULL
		} else if left.Type() == object.AccessByExpressionType {
			access := left.(object.AccessByExpression)
			switch left := access.Left.(type) {
			case object.Array:
				index := access.Expression.(object.Integer)
				left.Items[index.Value] = right
			case object.HashTable:
				key := access.Expression.(object.String)
				left.Items[key.Value] = right
			}
			return NULL
		} else {
			return newError("cannot assign value to: %s", left.Type())
		}
	case left.Type() == object.IdentifierType:
		return evalInfix(operator, left.(object.Identifier).Value, right, env)
	case left.Type() == object.AccessByExpressionType:
		return evalInfix(operator, left.(object.AccessByExpression).Value, right, env)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixOperators(operator string, left, right object.Integer) object.Object {
	switch operator {
	case "+":
		left.Value = left.Value + right.Value
		return left
	case "-":
		left.Value = left.Value - right.Value
		return left
	case "*":
		left.Value = left.Value * right.Value
		return left
	case "/":
		left.Value = left.Value / right.Value
		return left
	case ">":
		return nativeBoolToObject(left.Value > right.Value)
	case "<":
		return nativeBoolToObject(left.Value < right.Value)
	case "==":
		return nativeBoolToObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToObject(left.Value != right.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixOperators(operator string, left, right object.String) object.Object {
	switch operator {
	case "+":
		left.Value = left.Value + right.Value
		return left
	case ">":
		return nativeBoolToObject(left.Value > right.Value)
	case "<":
		return nativeBoolToObject(left.Value < right.Value)
	case "==":
		return nativeBoolToObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToObject(left.Value != right.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return FALSE
	default:
		return FALSE
	}
}

func evalMinusOperator(right object.Object) object.Object {
	switch r := right.(type) {
	case object.Integer:
		r.Value = -r.Value
		return r
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalIf(node ast.If, env *object.Environment) object.Object {
	for i, condition := range node.Conditions {
		evaluated := Eval(condition, env)

		if evaluated.Type() == object.ErrorType {
			return evaluated
		}

		if isTruthy(evaluated) {
			return Eval(node.Consequences[i], env)
		}
	}

	return Eval(node.Alternative, env)
}

func evalWhile(node ast.While, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if condition.Type() == object.ErrorType {
		return condition
	}

	if cond, ok := condition.(*object.Boolean); ok {
		for cond.Value {
			evaluated := evalBlock(node.Body.Statements, env)
			if evaluated.Type() == object.ErrorType {
				return evaluated
			}
			cond = Eval(node.Condition, env).(*object.Boolean)
		}
	}

	return NULL
}

func evalIdentifier(node ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return object.Identifier{Value: val, Name: node.Value}
	}

	if builtin, ok := builtins[node.Value]; ok {
		return object.Identifier{Value: builtin, Name: node.Value}
	}

	return newError("identifier not found: %s", node.Value)
}

func evalAccessByExpression(left object.Object, exp object.Object) object.Object {
	switch l := left.(type) {
	case object.Array:
		if index, ok := exp.(object.Integer); ok {
			if len(l.Items) <= index.Value || index.Value < 0 {
				return newError("index out of bounds: got=%d", index.Value)
			}
			return object.AccessByExpression{
				Left: l, Expression: index, Value: l.Items[index.Value],
			}
		} else {
			return newError("access expression is not integer: got %s", index.Type())
		}
	case object.HashTable:
		if key, ok := exp.(object.String); ok {
			value, ok := l.Items[key.Value]
			if !ok {
				return object.AccessByExpression{Left: l, Expression: key, Value: NULL}
			}
			return object.AccessByExpression{Left: l, Expression: key, Value: value}
		} else {
			return newError("keys in hash tables must be strings: got %s", key.Type())
		}
	case object.Identifier:
		return evalAccessByExpression(l.Value, exp)
	default:
		return newError("access by expression is not supported for this type: got %s", left.Type())
	}
}

func evalHashTable(node ast.HashTable, env *object.Environment) object.Object {
	items := make(map[string]object.Object, len(node.Items))

	for key, val := range node.Items {
		var keyString string

		switch k := key.(type) {
		case ast.String:
			keyString = k.Value
		case ast.Identifier:
			identifier := evalIdentifier(k, env)
			if identifier.Type() == object.ErrorType {
				return identifier
			}

			value := identifier.(object.Identifier).Value

			if str, ok := value.(object.String); ok {
				keyString = str.Value
			} else {
				return newError("keys in hash tables must be strings: got %s", k.Token.Type)
			}
		}

		evaluated := Eval(val, env)
		if evaluated.Type() == object.ErrorType {
			return evaluated
		}
		items[keyString] = evaluated
	}

	return object.HashTable{Items: items}
}

func evalFunction(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	function, ok := fn.(object.Function)
	if !ok {
		identifier, ok := fn.(object.Identifier)
		if !ok {
			return newError("not a function: %s", fn.String())
		}

		if value, ok := env.Get(identifier.Name); ok {
			if function, ok = value.(object.Function); !ok {
				return newError("not a function: %s", value.String())
			}
		} else {
			if builtin, ok := builtins[identifier.Name]; ok {
				return builtin.Function(args...)
			}
			return newError("no such function: %s", identifier.Name)
		}
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)

	if result, ok := evaluated.(object.Return); ok {
		return result.Value
	}

	return NULL
}
