package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifier
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	// Operator
	ASSIGN = "="
	PLUS   = "+"

	// Separator
	COMMA     = ","
	SEMICOLON = ";"

	L_PAREN = "("
	R_PAREN = ")"
	L_BRACE = "{"
	R_BRACE = "}"

	// Reserved
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
