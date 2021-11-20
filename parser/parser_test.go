package parser

import (
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
	identifier, ok := statement.Expression.(*ast.Identifier)
	require.True(t, ok, "expression is not Identifier, %s", statement.Expression)
	require.Equal(t, "foobar", identifier.Value, "identifier.Value is not foobar, %s", identifier.Value)
	require.Equal(t, "foobar", identifier.TokenLiteral(), "identifier.TokenLiteral() is not foobar, %s", identifier.TokenLiteral())
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
	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	require.True(t, ok, "expression is not IntegerLiteral, %s", statement.Expression)
	require.Equal(t, int64(5), literal.Value, "literal.Value is not 5, %d", literal.Value)
	require.Equal(t, "5", literal.TokenLiteral(), "literal.TokenLiteral() is not foobar, %s", literal.TokenLiteral())
}
