package parser

import (
	"strconv"
	"strings"
	"testing"

	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
)

func eq[T comparable](t *testing.T, expected T, actual T, msg ...string) {
	if expected != actual {
		t.Fatalf("%s\nexpected: %+v\nactual: %+v\n", strings.Join(msg, " "), expected, actual)
	}
}

func notEq[T comparable](t *testing.T, expected T, actual T, msg ...string) {
	if expected == actual {
		t.Fatalf("%s\nexpected not: %+v\nactual: %+v\n", strings.Join(msg, " "), expected, actual)
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedName string) bool {
		eq(t, "let", stmt.TokenLiteral(), "Got non let statement")

		letStmt, ok := stmt.(*ast.LetStatement)

		eq(t, true, ok, "Failed to typecast statement to let statement")
		eq(t, expectedName, letStmt.Name.Value, "Identifier name didn't match")
		eq(t, expectedName, letStmt.Name.TokenLiteral(), "Token literal of identifier didn't match the identifier name")
    return true
}

func Test_LetStatements(t *testing.T) {
	input := `let x = 5;
  let y = 10;
  let foobar = 8723456;`

	l := lexer.New_V2(strings.NewReader(input))
	p := New(l)

	program := p.ParseProgram()

	notEq(t, nil, program, "ParseProgram() returned nil")
	eq(t, 3, len(program.Statements), "program.Statements does not contain 3 statements")

	expectedIdentifiers := []string{"x", "y", "foobar"}
	for i, expectedIdentifier := range expectedIdentifiers {
		stmt := program.Statements[i]

		notEq(t, nil, stmt, "Got null statement at", strconv.Itoa(i))
    eq(t, true, testLetStatement(t, stmt, expectedIdentifier))
	}
}
