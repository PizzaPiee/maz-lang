package repl

import (
	"bufio"
	"fmt"
	// "maz-lang/evaluator"
	"maz-lang/lexer"
	"maz-lang/parser"
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
		p := parser.New(&l)
		program := p.Parse(token.EOF)
		// obj := evaluator.Eval(&program)
		// fmt.Printf("%s\n", obj.Inspect())
		fmt.Println(program.String())
	}
}

func getTokens(l *lexer.Lexer) []token.Token {
	var res []token.Token

	for {
		tok := l.NextToken()
		if tok.Type == token.EOF || tok.Type == token.ILLEGAL {
			return res
		}
		res = append(res, tok)
	}
}
