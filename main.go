package main

import (
	"fmt"
	"log"
	"maz-lang/environment"
	"maz-lang/evaluator"
	"maz-lang/lexer"
	"maz-lang/parser"
	"maz-lang/repl"
	"maz-lang/token"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		readAndEvalFromFile(os.Args[1])
	} else {
		repl.Run()
	}
}

func readAndEvalFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read file: %s\n", err)
	}

	env := environment.New()
	l := lexer.New(string(data))
	p := parser.New(&l)
	program := p.Parse(token.EOF)
	obj := evaluator.Eval(&program, &env)
	fmt.Printf("%s\n", obj.Inspect())
}
