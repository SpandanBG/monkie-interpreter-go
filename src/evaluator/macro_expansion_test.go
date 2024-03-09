package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/object"
	"sudocoding.xyz/interpreter_in_go/src/parser"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
)

func testParseProgram(input string) *ast.Program {
	l := lexer.New_V2(strings.NewReader(input))
	p := parser.New(l)
	return p.ParseProgram()
}

func TestDefineMacros(t *testing.T) {
	input := `
    let number = 1;
    let function = fn(x, y) { x + y };
    let mymacro = macro(x, y) { x + y; };
    `

	env := object.NewEnvironment()
	program := testParseProgram(input)

	DefineMacros(program, env)

	eq(t, 2, len(program.Statements), "Expected 2 program statemtns")

	_, ok := env.Get("number")
	eq(t, false, ok, "number should not be defined")

	_, ok = env.Get("function")
	eq(t, false, ok, "function should not be defined")

	obj, ok := env.Get("mymacro")
	eq(t, true, ok, "macro not in environment.")

	macro, ok := obj.(*object.Macro)
	eq(t, true, ok, fmt.Sprintf("object is not Macro. got=%T (%+v)", obj, obj))
	eq(t, 2, len(macro.Parameters), "Expected 2 macro params")
	eq(t, "x", macro.Parameters[0].String(), "not expected param")
	eq(t, "y", macro.Parameters[1].String(), "not expected param")

	eq(t, "(x + y)", macro.Body.String(), "not expected body")
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
      let infixExpression = macro() { quote(1 + 2); };

      infixExpression();
      `,
			`(1 + 2)`,
		},
		{
			`
      let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };

      reverse(2 + 2, 10 - 5);
      `,
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
        let unless = macro(condition, consequence, alternative) {
            quote(if (!(unquote(condition))) {
                unquote(consequence);
            } else {
                unquote(alternative);
            });
        };

        unless(10 > 5, puts("not greater"), puts("greater"));
        `,
			`if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(tt.expected)
		program := testParseProgram(tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		eq(t, expected.String(), expanded.String(), "didn't match expected")
	}
}
