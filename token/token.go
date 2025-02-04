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

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	SEMICOLON = ";"
	COMMA     = ","

	LBRACE = "{"
	RBRACE = "}"
	LPAREN = "("
	RPAREN = ")"

	EQ   = "=="
	NEQ  = "!="
	LT   = "<"
	GT   = ">"
	LTEQ = "<="
	GTEQ = ">="

	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	FUNCTION = "FUNCTION"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"fn":     FUNCTION,
}

func Lookupkeyword(word string) TokenType {
	keyword, ok := keywords[word]
	if ok {
		return keyword
	}

	return ""
}
