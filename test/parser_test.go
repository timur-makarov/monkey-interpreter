package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/timur-makarov/monkey-interpreter/internal/ast"
	"github.com/timur-makarov/monkey-interpreter/internal/lexer"
	"github.com/timur-makarov/monkey-interpreter/internal/parser"
	"github.com/timur-makarov/monkey-interpreter/internal/token"
)

func TestLetStatements(t *testing.T) {
	input := `
		let five = 5;
		let десять = 10 + 10;
	`

	tests := []struct {
		expectedIdentifier string
		expectedValue      string
	}{{"five", "5"}, {"десять", "10 + 10"}}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testLetStatement(t, statement, test.expectedIdentifier, test.expectedValue)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name, value string) {
	assert.Equal(t, "let", s.TokenLiteral())

	letStmt, ok := s.(ast.LetStatement)
	assert.Equal(t, true, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
	assert.Equal(t, value, letStmt.Value.String())
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10 - 10;
	`

	tests := []struct {
		expectedValue string
	}{{"5"}, {"10 - 10"}}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testReturnStatement(t, statement, test.expectedValue)
	}
}

func testReturnStatement(t *testing.T, s ast.Statement, value string) {
	assert.Equal(t, "return", s.TokenLiteral())

	rs, ok := s.(ast.ReturnStatement)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, rs.Value.String())
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	for _, e := range p.Errors() {
		t.Fatalf("Error has occured: %s", e.Message)
	}
	assert.Len(t, p.Errors(), 0)
}

func TestIdentifiers(t *testing.T) {
	input := `
		five;
	`

	tests := []struct {
		expectedIdentifier string
	}{{"five"}}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testIdentifierExpression(t, statement, test.expectedIdentifier)
	}
}

func testIdentifierExpression(t *testing.T, s ast.Statement, value string) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	ident, ok := statement.Expression.(ast.Identifier)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, ident.Value)
	assert.Equal(t, value, ident.TokenLiteral())
}

func TestIntegers(t *testing.T) {
	input := `
		5;
	`

	tests := []struct{ expectedInteger int }{{5}}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testIntegerExpression(t, statement, test.expectedInteger)
	}
}

func testIntegerExpression(t *testing.T, s ast.Statement, value int) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	integer, ok := statement.Expression.(ast.Integer)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, integer.Value)
	assert.Equal(t, strconv.Itoa(value), integer.TokenLiteral())
}

func TestStrings(t *testing.T) {
	input := `
		"hello world";
	`

	tests := []struct{ expected string }{{"hello world"}}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testStringExpression(t, statement, test.expected)
	}
}

func testStringExpression(t *testing.T, s ast.Statement, value string) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	str, ok := statement.Expression.(ast.String)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, str.Value)
}

func TestPrefixes(t *testing.T) {
	input := `
		!5;
		-10;
		!true;
		!false;
		!(false == false);
	`

	tests := []struct {
		expectedOperator string
		expectedValue    string
	}{
		{token.BANG, "5"}, {token.MINUS, "10"}, {token.BANG, "true"}, {token.BANG, "false"},
		{token.BANG, "false == false"},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testPrefixExpression(t, statement, test.expectedOperator, test.expectedValue)
	}
}

func testPrefixExpression(
	t *testing.T,
	s ast.Statement,
	operator string,
	value string,
) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	exp, ok := statement.Expression.(ast.Prefix)
	assert.Equal(t, true, ok)
	assert.Equal(t, operator, exp.Operator)
	assert.Equal(t, value, exp.Right.String())
}

func TestInfixes(t *testing.T) {
	input := `
		10 + 10 / 2;
		10 - 10 / 5;
		10 / 10 * 100 / 100;
		10 / 10 / 10 * 100;
		10 > 10;
		10 < 10;
		10 == 10;
		10 != 10;
		true == true;
		(1 + 2) * 6 + 4;
		x = x * y > x * zz;
	`

	tests := []struct {
		expectedOperator   string
		expectedLeftValue  string
		expectedRightValue string
	}{
		{token.PLUS, "10", "10 / 2"},
		{token.MINUS, "10", "10 / 5"},
		{token.DIVIDE, "10 / 10 * 100", "100"},
		{token.MULTIPLY, "10 / 10 / 10", "100"},
		{token.GT, "10", "10"},
		{token.LT, "10", "10"},
		{token.EQ, "10", "10"},
		{token.NEQ, "10", "10"},
		{token.EQ, "true", "true"},
		{token.PLUS, "1 + 2 * 6", "4"},
		{token.ASSIGN, "x", "x * y > x * zz"},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testInfixExpression(
			t, statement, test.expectedOperator, test.expectedLeftValue, test.expectedRightValue,
		)
	}
}

func testInfixExpression(
	t *testing.T,
	s ast.Statement,
	operator string,
	leftValue string,
	rightValue string,
) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	exp, ok := statement.Expression.(ast.Infix)
	assert.Equal(t, true, ok)
	assert.Equal(t, operator, exp.Operator)
	assert.Equal(t, leftValue, exp.Left.String())
	assert.Equal(t, rightValue, exp.Right.String())
}

func TestBooleans(t *testing.T) {
	input := `
		true;
		false;
	`

	tests := []struct {
		expectedBoolean string
	}{{"true"}, {"false"}}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testBooleanExpression(t, statement, test.expectedBoolean)
	}
}

func testBooleanExpression(t *testing.T, s ast.Statement, value string) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	b, ok := statement.Expression.(ast.Boolean)
	assert.Equal(t, true, ok)
	assert.Equal(t, value, strconv.FormatBool(b.Value))
	assert.Equal(t, value, b.TokenLiteral())
}

func TestIfs(t *testing.T) {
	input := `
		if (x < y) { x }

		if (y > x) { x } 
		else if (yy > xx) { yy } 
		else if (yy > xx) { yyy } 
		else if (yy > xx) { yyyy } 
		else { y }
	`

	tests := []struct {
		expectedConditions   []map[string]string
		expectedConsequences []string
		expectedAlternative  string
	}{
		{
			expectedConditions: []map[string]string{
				{
					"operator": token.LT,
					"left":     "x",
					"right":    "y",
				},
			},
			expectedConsequences: []string{"x"},
		},
		{
			expectedConditions: []map[string]string{
				{
					"operator": token.GT,
					"left":     "y",
					"right":    "x",
				},
				{
					"operator": token.GT,
					"left":     "yy",
					"right":    "xx",
				},
				{
					"operator": token.GT,
					"left":     "yy",
					"right":    "xx",
				},
				{
					"operator": token.GT,
					"left":     "yy",
					"right":    "xx",
				},
			},
			expectedConsequences: []string{"x", "yy", "yyy", "yyyy"},
			expectedAlternative:  "y",
		},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testIfExpression(
			t, statement, test.expectedConditions, test.expectedConsequences,
			test.expectedAlternative,
		)
	}
}

func testIfExpression(
	t *testing.T,
	s ast.Statement,
	conditions []map[string]string,
	consequences []string,
	alternative string,
) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	i, ok := statement.Expression.(ast.If)
	assert.Equal(t, true, ok)
	assert.Len(t, i.Conditions, len(consequences))
	assert.Len(t, i.Consequences, len(consequences))

	for j := range i.Conditions {
		condition := conditions[j]
		testInfixExpression(
			t, ast.ExpressionStatement{Expression: i.Conditions[j]}, condition["operator"],
			condition["left"], condition["right"],
		)

		testIdentifierExpression(t, i.Consequences[j].Statements[0], consequences[j])
	}

	if alternative != "" {
		assert.Len(t, i.Alternative.Statements, 1)
		testIdentifierExpression(t, i.Alternative.Statements[0], alternative)
	}
}

func TestFunctions(t *testing.T) {
	input := `
		fn(x, y) { x + y; };
		fn(y, x) { y - x };
		fn(y, x, z) { y - x };
        fn() { 1 - 2 };
	`

	tests := []struct {
		expectedBody       map[string]string
		expectedParameters []string
	}{
		{
			expectedBody: map[string]string{
				"operator": token.PLUS,
				"left":     "x",
				"right":    "y",
			},
			expectedParameters: []string{"x", "y"},
		},
		{
			expectedBody: map[string]string{
				"operator": token.MINUS,
				"left":     "y",
				"right":    "x",
			},
			expectedParameters: []string{"y", "x"},
		},
		{
			expectedBody: map[string]string{
				"operator": token.MINUS,
				"left":     "y",
				"right":    "x",
			},
			expectedParameters: []string{"y", "x", "z"},
		},
		{
			expectedBody: map[string]string{
				"operator": token.MINUS,
				"left":     "1",
				"right":    "2",
			},
			expectedParameters: []string{},
		},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testFunctionExpression(t, statement, test.expectedBody, test.expectedParameters)
	}
}

func testFunctionExpression(
	t *testing.T,
	s ast.Statement,
	body map[string]string,
	parameters []string,
) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	i, ok := statement.Expression.(ast.Function)
	assert.Equal(t, true, ok)

	assert.Len(t, i.Body.Statements, 1)
	testInfixExpression(t, i.Body.Statements[0], body["operator"], body["left"], body["right"])

	assert.Len(t, i.Parameters, len(parameters))

	for index, p := range parameters {
		testIdentifierExpression(t, ast.ExpressionStatement{Expression: i.Parameters[index]}, p)
	}
}

func TestCalls(t *testing.T) {
	input := `
		add(1, 3);
	`

	tests := []struct {
		expectedIdentifier string
		expectedArguments  []string
	}{
		{
			expectedIdentifier: "add",
			expectedArguments:  []string{"1", "3"},
		},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testCallExpression(t, statement, test.expectedIdentifier, test.expectedArguments)
	}
}

func testCallExpression(
	t *testing.T,
	s ast.Statement,
	identifier string,
	arguments []string,
) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	i, ok := statement.Expression.(ast.Call)
	assert.Equal(t, true, ok)

	f, ok := i.Function.(ast.Identifier)
	assert.Equal(t, true, ok)
	assert.Equal(t, identifier, f.Value)
	assert.Len(t, i.Arguments, len(arguments))

	for index, a := range arguments {
		val, _ := strconv.Atoi(a)
		testIntegerExpression(t, ast.ExpressionStatement{Expression: i.Arguments[index]}, val)
	}
}

func TestWhiles(t *testing.T) {
	input := `
		while (x < y) { x = y > 1 }
		while (y > x) { x = y == 1 }
	`

	tests := []struct {
		expectedCondition map[string]string
		expectedBody      string
	}{
		{
			expectedCondition: map[string]string{
				"operator": token.LT,
				"left":     "x",
				"right":    "y",
			},
			expectedBody: "x = y > 1",
		},
		{
			expectedCondition: map[string]string{
				"operator": token.GT,
				"left":     "y",
				"right":    "x",
			},
			expectedBody: "x = y == 1",
		},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testWhileExpression(t, statement, test.expectedCondition, test.expectedBody)
	}
}

func testWhileExpression(
	t *testing.T,
	s ast.Statement,
	condition map[string]string,
	body string,
) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	i, ok := statement.Expression.(ast.While)
	assert.Equal(t, true, ok)

	testInfixExpression(
		t, ast.ExpressionStatement{Expression: i.Condition}, condition["operator"],
		condition["left"], condition["right"],
	)

	assert.Equal(t, body, i.Body.String())
}

func TestArrays(t *testing.T) {
	input := `
		[1, 2];
		[3, 4, 5];
	`

	tests := []struct {
		expected []int
	}{
		{[]int{1, 2}},
		{[]int{3, 4, 5}},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testArrayExpression(t, statement, test.expected)
	}
}

func testArrayExpression(t *testing.T, s ast.Statement, items []int) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	array, ok := statement.Expression.(ast.Array)
	assert.Equal(t, true, ok)

	for i, item := range array.Items {
		integer, ok := item.(ast.Integer)
		assert.Equal(t, true, ok)
		assert.Equal(t, items[i], integer.Value)
	}
}

func TestAccessByExpression(t *testing.T) {
	input := `
		myArray[1 + 1]
	`

	tests := []struct {
		left  string
		index string
	}{{"myArray", "1 + 1"}}

	program := getProgram(t, input)

	for i, test := range tests {
		statement := program.Statements[i]
		testAccessByExpression(t, statement, test.left, test.index)
	}
}

func testAccessByExpression(t *testing.T, s ast.Statement, left, index string) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	indx, ok := statement.Expression.(ast.AccessByExpression)
	assert.Equal(t, true, ok)

	l, ok := indx.Left.(ast.Identifier)
	assert.Equal(t, true, ok)
	assert.Equal(t, left, l.String())

	infx, ok := indx.Index.(ast.Infix)
	assert.Equal(t, true, ok)
	assert.Equal(t, index, infx.String())
}

func TestHashTables(t *testing.T) {
	input := `
		{"hello": 1, "world": 2};
	`

	tests := []struct {
		expected map[string]int
	}{
		{map[string]int{"hello": 1, "world": 2}},
	}

	program := getProgram(t, input)
	assert.Len(t, program.Statements, len(tests))

	for i, test := range tests {
		statement := program.Statements[i]
		testHashTableExpression(t, statement, test.expected)
	}
}

func testHashTableExpression(t *testing.T, s ast.Statement, items map[string]int) {
	statement, ok := s.(ast.ExpressionStatement)
	assert.Equal(t, true, ok)

	hashTable, ok := statement.Expression.(ast.HashTable)
	assert.Equal(t, true, ok)

	for key, value := range hashTable.Items {
		k, ok := key.(ast.String)
		assert.Equal(t, true, ok)

		val, ok := value.(ast.Integer)
		assert.Equal(t, true, ok)
		assert.Equal(t, items[k.Value], val.Value)
	}
}

func getProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	assert.NotEqual(t, nil, program)

	return program
}
