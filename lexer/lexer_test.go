package lexer

import (
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	input := `=+(){},;`

	testCases := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.L_PAREN, "("},
		{token.R_PAREN, ")"},
		{token.L_BRACE, "{"},
		{token.R_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := New(input)
	for _, testCase := range testCases {
		token := l.NextToken()
		assert.Equal(t, token.Type, testCase.expectedType, "Wrong token type")
		assert.Equal(t, token.Literal, testCase.expectedLiteral, "Wrong literal")
	}
}
