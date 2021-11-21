package parser

import (
	"fmt"
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
	require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())
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
	require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())
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
	require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	testIdentifier(t, "foobar", statement.Expression)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())
	require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])
	testIntegerLiteral(t, 5, statement.Expression)
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, prefixTest := range prefixTests {
		l := lexer.New(prefixTest.input)
		p := New(l)
		program := p.ParseProgram()
		require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())

		require.Equal(t, 1, len(program.Statements), "statement does not contain 1 statements, %s", program.Statements)
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "statements[0] is not ExpressionStatement, %s", program.Statements[0])

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		require.True(t, ok, "expression is not PrefixExpression, %s", statement.Expression)
		require.Equal(t, prefixTest.operator, expression.Operator, "Wrong operator")
		testIntegerLiteral(t, prefixTest.integerValue, expression.Right)
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, infixTest := range infixTests {
		l := lexer.New(infixTest.input)
		p := New(l)
		program := p.ParseProgram()
		require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())

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
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5;", "(3 + 4)((-5) * 5)"},
		{"5 > 4==3<4", "((5 > 4) == (3 < 4))"},
		{"5 < 4!=3>4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}
	for _, precedenceTest := range precedenceTests {
		l := lexer.New(precedenceTest.input)
		p := New(l)
		program := p.ParseProgram()
		require.Equal(t, 0, len(p.Errors()), "parser errors: %s", p.Errors())

		actual := program.String()
		require.Equal(t, precedenceTest.expected, actual, "wrong parsing")
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

func testLiteralExpression(t *testing.T, expected interface{}, actual ast.Expression) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, int64(v), actual)
	case int64:
		testIntegerLiteral(t, v, actual)
	case string:
		testIdentifier(t, v, actual)
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
