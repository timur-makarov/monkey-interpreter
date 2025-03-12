package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/timur-makarov/monkey-interpreter/internal/evaluator"
	"github.com/timur-makarov/monkey-interpreter/internal/object"
)

func TestEvaluatedIntegers(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5", 5}, {"10", 10}, {"20", 20},
		{"-5", -5}, {"-10", -10}, {"-20", -20},
		{"11 * 5", 55}, {"0 - 10", -10}, {"40 - 20", 20}, {"-20 * 5", -100}, {"(1 + 2) * 4", 12},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func testIntegerObject(t *testing.T, o object.Object, value int) {
	obj, ok := o.(object.Integer)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, obj.Value)
	assert.Equal(t, object.IntegerType, obj.Type())
}

func TestEvaluatedStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello\"", "hello"},
		{"let x = \"hello \" + \"world\"; x", "hello world"},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testStringObject(t, evaluated, test.expected)
	}
}

func testStringObject(t *testing.T, o object.Object, value string) {
	obj, ok := o.(object.String)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, obj.Value)
	assert.Equal(t, object.StringType, obj.Type())
}

func TestEvaluatedBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true}, {"false", false}, {"1 > 2", false}, {"1 + 1 == 2", true},
		{"(1 - 2) * 4 < 10", true}, {"true == true", true}, {"false != true", true},
		{"(1 < 2) == true", true}, {"(1 > 2) == true", false},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func testBooleanObject(t *testing.T, o object.Object, value bool) {
	obj, ok := o.(*object.Boolean)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, obj.Value)
	assert.Equal(t, object.BooleanType, obj.Type())
}

func TestEvaluatedBangs(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false}, {"!false", true}, {"!5", false}, {"!!true", true}, {"!!false", false},
		{"!!5", true}, {"!!!true", false}, {"!!!false", true}, {"!!!5", false},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestEvaluatedIfElse(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (1 > 2) { 5 }", nil},
		{"if (1 > 2) { 5 } else { 1 }", 1},
		{"if (1 > 2) { 5 } else if (2 > 1) { 15 } else { 1 }", 15},
		{"if (1 > 2) { 5 } else if (2 == 1) { 15 } else if (2 > 1) { 25 } else { 1 }", 25},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		switch val := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, val)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, o object.Object) {
	assert.Equal(t, evaluator.NULL, o)
}

func TestEvaluatedReturns(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 5", 5}, {"return 10; 5", 10}, {"5; return 2 * 5; 5", 10},
		{
			`
			if (2 > 1) {
				if (2 > 1) {
					return 10
				}
				return 1
			}
		`, 10,
		},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvaluatedErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (2 > 1) { true + false }", "unknown operator: BOOLEAN + BOOLEAN"},
		{
			"if (2 > 1) { if (2 > 1) { true + false} return 1 }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"z", "identifier not found: z"},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testErrorObject(t, evaluated, test.expected)
	}
}

func testErrorObject(t *testing.T, o object.Object, message string) {
	obj, ok := o.(object.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, message, obj.Message)
	assert.Equal(t, object.ErrorType, obj.Type())
}

func TestEvaluatedLets(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let x = 5; x", 5},
		{"let y = 15 + 5; y", 20},
		{"let z = -11 * (10 * -1); let zz = z + z; zz", 220},
		{"let z = -11 * (10 * -1); let zz = z + z; zz = 10; zz", 10},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvaluatedFunctions(t *testing.T) {
	tests := []struct {
		input      string
		parameters []string
		body       string
	}{
		{"fn(x) {return x + 2}", []string{"x"}, "return x + 2"},
		{"fn(x, y) {return x * y}", []string{"x", "y"}, "return x * y"},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		testFunctionObject(t, evaluated, test.parameters, test.body)
	}
}

func testFunctionObject(t *testing.T, o object.Object, parameters []string, body string) {
	obj, ok := o.(object.Function)
	assert.Equal(t, true, ok)
	assert.Equal(t, object.FunctionType, obj.Type())

	assert.Len(t, obj.Parameters, len(parameters))
	for i, ident := range obj.Parameters {
		assert.Equal(t, ident.Value, parameters[i])
	}

	assert.Equal(t, obj.Body.String(), body)
}

func TestEvaluatedFunctionCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let add = fn(x) {return x + 2}; add(2)", 4},
		{"let mul = fn(x, y) {return x * y}; mul(3, 3)", 9},
	}

	for _, test := range tests {
		evaluated := testEvalWithError(t, test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvaluatedBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"len(\"\")", 0},
		{"len(\"four\")", 4},
		{`len([1, 2, 3, 4])`, 4},
		{"len(1)", "argument type is not supported: got INTEGER"},
		{"len(\"four\", \"three\")", "wrong number of arguments: got=2, want=1"},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)
		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}
}

func TestEvaluatedArrays(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"[1, 2, 3]", []any{1, 2, 3}},
		{"[10, 22, 33]", []any{10, 22, 33}},
		{`["hello", "world"]`, []any{"hello", "world"}},
	}

	for _, test := range tests {
		evaluated := testEval(t, test.input)

		array, ok := evaluated.(object.Array)
		assert.Equal(t, true, ok)

		for i, item := range array.Items {
			switch it := item.(type) {
			case object.Integer:
				testIntegerObject(t, it, test.expected[i].(int))
			case object.String:
				testStringObject(t, it, test.expected[i].(string))
			}
		}
	}
}

func TestEvaluatedAccessByExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"let x = [1, 2, 3]; x[1]", 2},
		{"let y = [10, 22, 33]; y[1 + 1]", 33},
		{`let words = ["hello", "world"]; words[(100 + 100) * 0]`, "hello"},
	}

	for _, test := range tests {
		evaluated := testEvalWithError(t, test.input)
		switch ev := evaluated.(type) {
		case object.Integer:
			testIntegerObject(t, ev, test.expected.(int))
		case object.String:
			testStringObject(t, ev, test.expected.(string))
		}
	}
}

func TestEvaluatedHashTables(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]any
	}{
		{
			`let hello = "variable"; {hello: 1 + 1 * 100, "world": 2}`,
			map[string]any{"variable": 101, "world": 2},
		},
	}

	for _, test := range tests {
		evaluated := testEvalWithError(t, test.input)

		hashTable, ok := evaluated.(object.HashTable)
		assert.Equal(t, true, ok)

		for key, val := range hashTable.Items {
			switch it := val.(type) {
			case object.Integer:
				testIntegerObject(t, it, test.expected[key].(int))
			case object.String:
				testStringObject(t, it, test.expected[key].(string))
			}
		}
	}
}

func testEval(t *testing.T, input string) object.Object {
	program := getProgram(t, input)
	env := object.NewEnvironment()

	return evaluator.Eval(program, env)
}

func testEvalWithError(t *testing.T, input string) object.Object {
	program := getProgram(t, input)
	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if err, ok := evaluated.(object.Error); ok {
		t.Fatalf("Error evaluating: %s", err.Message)
	}

	return evaluated
}
