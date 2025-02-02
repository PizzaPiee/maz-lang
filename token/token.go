package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"

	SEMICOLON = ";"
	COMMA     = ","

	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"
)
