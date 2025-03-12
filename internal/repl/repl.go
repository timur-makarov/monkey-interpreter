package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/timur-makarov/monkey-interpreter/internal/evaluator"
	"github.com/timur-makarov/monkey-interpreter/internal/lexer"
	"github.com/timur-makarov/monkey-interpreter/internal/object"
	"github.com/timur-makarov/monkey-interpreter/internal/parser"
)

const PROMPT = "Monkey code -> "

func ReadUserInput(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scannedLine := scanner.Scan()
		if !scannedLine {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			log.Println(p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			_, _ = io.WriteString(out, evaluated.String()+"\n")
		}
	}
}
