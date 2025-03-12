package object

import (
	"fmt"

	"github.com/timur-makarov/monkey-interpreter/internal/ast"
)

type Type string

const (
	IntegerType            Type = "INTEGER"
	StringType             Type = "STRING"
	BooleanType            Type = "BOOLEAN"
	NullType               Type = "NULL"
	ReturnType             Type = "RETURN"
	ErrorType              Type = "ERROR"
	IdentifierType         Type = "IDENTIFIER"
	FunctionType           Type = "FUNCTION"
	BuiltinType            Type = "BUILTIN"
	ArrayType              Type = "ARRAY"
	HashTableType          Type = "HASHTABLE"
	AccessByExpressionType Type = "ACCESSBYEXPRESSION"
)

type Object interface {
	Type() Type
	String() string
}

type Integer struct {
	Value int
}

func (i Integer) Type() Type {
	return IntegerType
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type String struct {
	Value string
}

func (s String) Type() Type {
	return StringType
}

func (s String) String() string {
	return s.Value
}

type Boolean struct {
	Value bool
}

func (b Boolean) Type() Type {
	return BooleanType
}

func (b Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n Null) Type() Type {
	return NullType
}

func (n Null) String() string {
	return "null"
}

type Return struct {
	Value Object
}

func (r Return) Type() Type {
	return ReturnType
}

func (r Return) String() string {
	return fmt.Sprintf("%s", r.Value)
}

type Identifier struct {
	Name  string
	Value Object
}

func (i Identifier) Type() Type {
	return IdentifierType
}

func (i Identifier) String() string {
	return fmt.Sprintf("%s = %s", i.Name, i.Value)
}

type Error struct {
	Message string
}

func (r Error) Type() Type {
	return ErrorType
}

func (r Error) String() string {
	return fmt.Sprintf("ERROR: %s", r.Message)
}

type Function struct {
	Parameters []ast.Identifier
	Body       ast.BlockStatement
	Env        *Environment
}

func (f Function) Type() Type {
	return FunctionType
}

func (f Function) String() string {
	return fmt.Sprintf("fn(%+v) {%s}", f.Parameters, f.Body)
}

type Array struct {
	Items []Object
}

func (a Array) Type() Type {
	return ArrayType
}

func (a Array) String() string {
	return fmt.Sprintf("%+v", a.Items)
}

type HashTable struct {
	Items map[string]Object
}

func (ht HashTable) Type() Type {
	return HashTableType
}

func (ht HashTable) String() string {
	return fmt.Sprintf("%+v", ht.Items)
}

type AccessByExpression struct {
	Left       Object
	Expression Object
	Value      Object
}

func (a AccessByExpression) Type() Type {
	return AccessByExpressionType
}

func (a AccessByExpression) String() string {
	return fmt.Sprintf("%s[%s]", a.Left, a.Expression)
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Get(key string) (Object, bool) {
	cur := e

	for cur != nil {
		if obj, ok := cur.store[key]; ok {
			return obj, ok
		}
		cur = cur.outer
	}

	return nil, false
}

func (e *Environment) Set(key string, value Object) {
	cur := e

	for cur != nil {
		if _, ok := cur.store[key]; ok {
			cur.store[key] = value
			return
		}
		cur = cur.outer
	}

	e.store[key] = value
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Function BuiltinFunction
}

func (b Builtin) Type() Type {
	return BuiltinType
}

func (b Builtin) String() string {
	return "builtin function"
}
