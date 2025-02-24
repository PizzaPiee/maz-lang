package repl

import (
	"bufio"
	"fmt"
	// "maz-lang/environment"
	// "maz-lang/evaluator"
	"maz-lang/lexer"
	"maz-lang/parser"
	"maz-lang/token"
	"os"
)

func Run() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Maz REPL!")
	// env := environment.New()
	for {
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')

		l := lexer.New(input)
		p := parser.New(&l)
		program := p.Parse(token.EOF)
		// obj := evaluator.Eval(&program, &env)
		fmt.Printf("%s\n", program.String())
		// fmt.Printf("%s\n", obj.Inspect())
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
