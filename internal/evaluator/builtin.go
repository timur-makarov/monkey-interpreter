package evaluator

import (
	"log"

	"github.com/timur-makarov/monkey-interpreter/internal/object"
)

type BuiltinFunctions struct{}

func (bf BuiltinFunctions) len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments: got=%d, want=1", len(args))
	}

	switch item := args[0].(type) {
	case object.String:
		return object.Integer{Value: len(item.Value)}
	case object.Array:
		return object.Integer{Value: len(item.Items)}
	case object.Identifier:
		return bf.len(item.Value)
	default:
		return newError("argument type is not supported: got %s", item.Type())
	}
}

func (bf BuiltinFunctions) shift(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments: got=%d, want=1", len(args))
	}

	switch item := args[0].(type) {
	case object.Array:
		if len(item.Items) > 0 {
			return object.Array{Items: append([]object.Object{}, item.Items[1:]...)}
		} else {
			return item
		}
	case object.Identifier:
		return bf.shift(item.Value)
	default:
		return newError("argument type is not supported: got %s", item.Type())
	}
}

func (bf BuiltinFunctions) append(args ...object.Object) object.Object {
	if len(args) < 2 {
		return newError("wrong number of arguments: got=%d, want=>1", len(args))
	}

	switch item := args[0].(type) {
	case object.Array:
		if len(item.Items) > 0 {
			return object.Array{Items: append(item.Items, args[1:]...)}
		} else {
			return item
		}
	case object.Identifier:
		return bf.append(append([]object.Object{item.Value}, args[1:]...)...)
	default:
		return newError("argument type is not supported: got %s", item.Type())
	}
}

func (bf BuiltinFunctions) log(args ...object.Object) object.Object {
	strings := make([]any, len(args))

	for i, arg := range args {
		strings[i] = arg.String()
	}

	log.Println(strings...)
	return NULL
}

var bf = BuiltinFunctions{}

var builtins = map[string]object.Builtin{
	"len":    {bf.len},
	"shift":  {bf.shift},
	"append": {bf.append},
	"log":    {bf.log},
}
