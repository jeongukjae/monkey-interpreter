package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 -10", 10},
		{"2 * 2 * 2* 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"50 / 2 * 2 + 10", 60},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, tt.expected, evaluated)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 > 1", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"false == true", false},
		{"false != true", true},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, tt.expected, evaluated)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, tt.expected, evaluated)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 > 2) {10}", nil},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 < 2) {10} else {20}", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, int64(integer), evaluated)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 5; 9;", 5},
		{"return 2 * 5;9;", 10},
		{"9; return 2 * 5;9;", 10},
		{`
		if (10 > 1) {
			if (10  >1 ) {
				return 10;
			}
			return 5;
		}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, tt.expected, evaluated)
	}
}

//
// Helper functions
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	result, ok := actual.(*object.Intger)
	require.True(t, ok, "object is not integer")
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testBooleanObject(t *testing.T, expected bool, actual object.Object) {
	result, ok := actual.(*object.Boolean)
	require.True(t, ok, "object is not boolean, got %s", actual)
	require.Equal(t, expected, result.Value, "object has wrong value")
}

func testNullObject(t *testing.T, actual object.Object) {
	require.Equal(t, NULL, actual, "object is not NULL")
}
