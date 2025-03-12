package evaluator

import (
	"fmt"

	"github.com/timur-makarov/monkey-interpreter/internal/object"
)

func extendFunctionEnv(fn object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func nativeBoolToObject(value bool) object.Object {
	if value {
		return TRUE
	} else {
		return FALSE
	}
}

func isTruthy(obj object.Object) bool {
	switch o := obj.(type) {
	case object.Integer:
		return o.Value > 0
	case *object.Boolean:
		return o.Value
	case object.Null:
		return false
	default:
		return false
	}
}

func newError(format string, a ...any) object.Error {
	return object.Error{Message: fmt.Sprintf(format, a...)}
}
