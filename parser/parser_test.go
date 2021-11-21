package parser

import (
	"fmt"
	"log"
	"monkey/ast"
	"monkey/lexer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 5;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	testParserErrors(t, p)
	require.NotNil(t, program, "ParseProgram() returned nil")
	require.Equal(t, 3, len(program.Statements), "program.Statements does not contain 3 statements.")

	expectedIdentifiers := []string{"x", "y", "foobar"}
	for i, expectedIdentifier := range expectedIdentifiers {
		statement := program.Statements[i]
		(func(statement ast.Statement, expectedIdentifier string) {
			require.Equal(t, "let", statement.TokenLiteral(), "Wrong TokenLiternal")
			letStatement, ok := statement.(*ast.LetStatement)
			require.True(t, ok, "Wrong type")
			require.Equal(t, expectedIdentifier, letStatement.Name.Value, "Wrong name")
			require.Equal(t, expectedIdentifier, letStatement.Name.TokenLiteral(), "Wrong token literal")
		})(statement, expectedIdentifier)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 3, len(program.Statements), "statement does not contain 3 statements, %s", program.Statements)

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		require.True(t, ok, "statement is not return statement")
		require.Equal(t, "return", returnStatement.TokenLiteral(), "Wrong TokenLiteral")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	testLiteralExpression(t, "foobar", statement.Expression)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	testLiteralExpression(t, 5, statement.Expression)
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	testLiteralExpression(t, true, statement.Expression)
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"!true;", "!", true},
		{"!false;", "!", false},
		{"-15;", "-", 15},
	}

	for _, prefixTest := range prefixTests {
		l := lexer.New(prefixTest.input)
		p := New(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		require.True(t, ok, "expression is not PrefixExpression, %s", statement.Expression)
		require.Equal(t, prefixTest.operator, expression.Operator, "Wrong operator")
		testLiteralExpression(t, prefixTest.integerValue, expression.Right)
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true != false;", true, "!=", false},
		{"false != false;", false, "!=", false},
		{"false != true;", false, "!=", true},
	}

	for _, infixTest := range infixTests {
		l := lexer.New(infixTest.input)
		p := New(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])

		testInfixExpression(t, infixTest.leftValue, infixTest.operator, infixTest.rightValue, statement.Expression)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	precedenceTests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"true", "true"},
		{"false", "false"},
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5;", "(3 + 4)((-5) * 5)"},
		{"5 > 4==3<4", "((5 > 4) == (3 < 4))"},
		{"5 < 4!=3>4", "((5 < 4) != (3 > 4))"},
		{"5 < 4!=false", "((5 < 4) != false)"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"!(true == true ==false)", "(!((true == true) == false))"},
	}
	for _, precedenceTest := range precedenceTests {
		l := lexer.New(precedenceTest.input)
		p := New(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		actual := program.String()
		require.Equal(t, precedenceTest.expected, actual, "wrong parsing")
	}
}

func TestIfExpression(t *testing.T) {
	input := `if(x<y){x}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	expression, ok := statement.Expression.(*ast.IfExpression)
	require.True(t, ok, "Expression is not if expression, %s", expression)

	testInfixExpression(t, "x", "<", "y", expression.Condition)

	require.Equal(t, 1, len(expression.Consequence.Statements), "Consequence does not contain 1 statements, %s", expression.Consequence.Statements)
	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "consequence.statements[0] is not if expression, %s", consequence)
	testIdentifier(t, "x", consequence.Expression)

	require.Nil(t, expression.Alternative, "alternative is not nil, %s", expression.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if(x<y){x}else{y}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	expression, ok := statement.Expression.(*ast.IfExpression)
	require.True(t, ok, "Expression is not if expression, %s", expression)

	testInfixExpression(t, "x", "<", "y", expression.Condition)

	require.Equal(t, 1, len(expression.Consequence.Statements), "Consequence does not contain 1 statements, %s", expression.Consequence.Statements)
	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "consequence.statements[0] is not if expression, %s", consequence)
	testIdentifier(t, "x", consequence.Expression)

	require.Equal(t, 1, len(expression.Alternative.Statements), "Alternative does not contain 1 statements, %s", expression.Alternative.Statements)
	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "alternative.statements[0] is not if expression, %s", alternative)
	testIdentifier(t, "y", alternative.Expression)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	testParserErrors(t, p)
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	expression, ok := statement.Expression.(*ast.FunctionLiteral)
	require.True(t, ok, "Expression is not FunctionLiteral, %s", expression)

	require.Equal(t, 2, len(expression.Parameters), "len(parameters) != 2, %s", expression.Parameters)
	testLiteralExpression(t, "x", expression.Parameters[0])
	testLiteralExpression(t, "y", expression.Parameters[1])

	require.Equal(t, 1, len(expression.Body.Statements), "len(body.statements) != 1, %s", expression.Parameters)
	bodyStatement, ok := expression.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "body.statements[0] is not ExpressionStatement, %s", expression.Body.Statements[0])
	testInfixExpression(t, "x", "+", "y", bodyStatement.Expression)
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn(){};", []string{}},
		{"fn(x){};", []string{"x"}},
		{"fn(x, y, z){};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		testParserErrors(t, p)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		require.Equal(t, len(tt.expectedParams), len(function.Parameters), "Wrong parameter length")
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, ident, function.Parameters[i])
		}
	}
}

// Helper functionss
//
func testParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) != 0 {
		log.Printf("parser has %d erros", len(p.Errors()))
		for _, err := range p.Errors() {
			log.Printf("parser error: \"%s\"", err)
		}
		t.FailNow()
	}
}

func testIdentifier(t *testing.T, expected string, actual ast.Expression) {
	identifier, ok := actual.(*ast.Identifier)
	require.True(t, ok, "Expression is not identifier, %s", actual)
	require.Equal(t, expected, identifier.Value, "Wrong value")
	require.Equal(t, expected, identifier.TokenLiteral(), "Wrong token literal")
}

func testIntegerLiteral(t *testing.T, expected int64, actual ast.Expression) {
	integer, ok := actual.(*ast.IntegerLiteral)
	require.True(t, ok, "Expression is not integer literal, %s", actual)
	require.Equal(t, expected, integer.Value, "Wrong value")
	require.Equal(t, fmt.Sprintf("%d", expected), integer.TokenLiteral(), "Wrong token literal")
}

func testBoolLiteral(t *testing.T, expected bool, actual ast.Expression) {
	boolean, ok := actual.(*ast.Boolean)
	require.True(t, ok, "Expression is not boolean literal, %s", actual)
	require.Equal(t, expected, boolean.Value, "Wrong value")
	require.Equal(t, fmt.Sprintf("%t", expected), boolean.TokenLiteral(), "Wrong token literal")
}

func testLiteralExpression(t *testing.T, expected interface{}, actual ast.Expression) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, int64(v), actual)
	case int64:
		testIntegerLiteral(t, v, actual)
	case string:
		testIdentifier(t, v, actual)
	case bool:
		testBoolLiteral(t, v, actual)
	default:
		t.Errorf("type of expression not handled %T", actual)
	}
}

func testInfixExpression(t *testing.T, left interface{}, operator string, right interface{}, actual ast.Expression) {
	actualOperator, ok := actual.(*ast.InfixExpression)
	require.True(t, ok, "Expression is not infix expression, %s", actual)
	testLiteralExpression(t, left, actualOperator.Left)
	require.Equal(t, operator, actualOperator.Operator, "wrong operator")
	testLiteralExpression(t, right, actualOperator.Right)
}
