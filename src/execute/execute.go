package execute

import (
	"fmt"
	"os"

	"sudocoding.xyz/interpreter_in_go/src/evaluator"
	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/object"
	"sudocoding.xyz/interpreter_in_go/src/parser"
)

func Execute(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	l := lexer.New_V2(file)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		panic(p.Errors())
	}

	env := object.NewEnvironment()
	switch result := evaluator.Eval(program, env).(type) {
	case *object.Error:
		fmt.Println("Error Occured: ", result.Message)
	}
}
