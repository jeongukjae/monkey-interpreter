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
