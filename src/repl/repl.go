package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"sudocoding.xyz/interpreter_in_go/src/evaluator"
	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Println("Byeee")
			return
		}

		l := lexer.New_V2(strings.NewReader(line))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluator := evaluator.Eval(program)
		if evaluator == nil {
			continue
		}

		io.WriteString(out, evaluator.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err.Error()+"\n")
	}
}
