package lexer

import (
	"maz-lang/token"
)

type Lexer struct {
	Text    string
	pos     int
	readPos int
	char    byte
}

func New(text string) Lexer {
	l := Lexer{Text: text}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var res token.Token

	l.skipWhitespace()

	switch l.char {
	case 0:
		res = newToken(token.EOF, "")
	case '!':
		if l.peekChar() == '=' {
			res = newToken(token.NEQ, string(l.char)+string(l.peekChar()))
			l.readChar()
		} else {
			res = newToken(token.BANG, string(l.char))
		}
	case '+':
		res = newToken(token.PLUS, string(l.char))
	case '-':
		res = newToken(token.MINUS, string(l.char))
	case '*':
		res = newToken(token.ASTERISK, string(l.char))
	case '/':
		res = newToken(token.SLASH, string(l.char))
	case '=':
		if l.peekChar() == '=' {
			res = newToken(token.EQ, string(l.char)+string(l.peekChar()))
			l.readChar()
		} else {
			res = newToken(token.ASSIGN, string(l.char))
		}
	case ';':
		res = newToken(token.SEMICOLON, string(l.char))
	case ',':
		res = newToken(token.COMMA, string(l.char))
	case '(':
		res = newToken(token.LPAREN, string(l.char))
	case ')':
		res = newToken(token.RPAREN, string(l.char))
	case '{':
		res = newToken(token.LBRACE, string(l.char))
	case '}':
		res = newToken(token.RBRACE, string(l.char))
	case '>':
		if l.peekChar() == '=' {
			res = newToken(token.GTEQ, string(l.char)+string(l.peekChar()))
			l.readChar()
		} else {
			res = newToken(token.GT, string(l.char))
		}
	case '<':
		if l.peekChar() == '=' {
			res = newToken(token.LTEQ, string(l.char)+string(l.peekChar()))
			l.readChar()
		} else {
			res = newToken(token.LT, string(l.char))
		}
	case '"':
		str := l.readString()
		res = newToken(token.STRING, str)
	default:
		// Check if it is a digit
		if isDigit(l.char) {
			return newToken(token.INT, l.readNumber())
		}
		// Check if it is an identifier or keyword
		word := l.readWord()
		keyword := token.Lookupkeyword(word)

		if keyword != "" {
			return newToken(keyword, word)
		} else if word != "" {
			return newToken(token.IDENT, word)
		} else {
			return newToken(token.ILLEGAL, string(l.char))
		}
	}

	l.readChar()
	return res
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.Text) {
		l.char = 0
	} else {
		l.char = l.Text[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.Text) {
		return 0
	}
	return l.Text[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	start := l.readPos
	for l.peekChar() != '"' {
		l.readChar()
	}
	l.readChar()
	end := l.pos

	return l.Text[start:end]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.char) {
		l.readChar()
	}
	return l.Text[pos:l.pos]
}

func (l *Lexer) readWord() string {
	pos := l.pos
	for isLetter(l.char) {
		l.readChar()
	}
	return l.Text[pos:l.pos]
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}
