package repl

import (
	"bufio"
	"fmt"
	"maz-lang/lexer"
	"maz-lang/token"
	"os"
)

func Run() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Maz REPL!")
	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')

		l := lexer.New(input)
		tokens := getTokens(&l)
		fmt.Printf("%+v\n", tokens)
	}
}

func getTokens(l *lexer.Lexer) []token.Token {
	var res []token.Token

	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			return res
		}
		res = append(res, tok)
	}
}
