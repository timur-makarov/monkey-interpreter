package main

import (
	"log"
	"os"

	"github.com/timur-makarov/monkey-interpreter/internal/evaluator"
	"github.com/timur-makarov/monkey-interpreter/internal/lexer"
	"github.com/timur-makarov/monkey-interpreter/internal/object"
	"github.com/timur-makarov/monkey-interpreter/internal/parser"
	"github.com/timur-makarov/monkey-interpreter/internal/repl"
)

func main() {
	args := os.Args

	if len(args) > 1 {
		data, err := os.ReadFile(args[1])
		if err != nil {
			panic(err)
		}

		l := lexer.New(string(data))
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)

		if evaluated.Type() == object.ErrorType {
			log.Fatalln(evaluated)
		}
	} else {
		log.Println("Enter your Monkey code:")
		repl.ReadUserInput(os.Stdin, os.Stdout)
	}
}
