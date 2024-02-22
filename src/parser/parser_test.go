package parser

import (
	"fmt"
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

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	eq(t, ok, true, "Failed to typecase il to *ast.IntegerLiteral")
	eq(t, value, integ.Value, "integ.Value doesn't match expected")
	eq(t, fmt.Sprintf("%d", value), integ.TokenLiteral(), "integ.TokenLiteral() doesn't match expected")
	return true
}

func checkParserErrs(t *testing.T, p *Parser) bool {
	errs := p.Errors()
	if len(errs) == 0 {
		return true
	}

	t.Errorf("parser found %d errors:", len(errs))
	for _, err := range errs {
		t.Errorf("parser err: %q", err.Error())
	}

	return false
}

func Test_LetStatements(t *testing.T) {
	input := `let x = 5;
  let y = 10;
  let foobar = 8723456;`

	l := lexer.New_V2(strings.NewReader(input))
	p := New(l)

	program := p.ParseProgram()

	eq(t, true, checkParserErrs(t, p))
	notEq(t, nil, program, "ParseProgram() returned nil")
	eq(t, 3, len(program.Statements), "program.Statements does not contain 3 statements")

	expectedIdentifiers := []string{"x", "y", "foobar"}
	for i, expectedIdentifier := range expectedIdentifiers {
		stmt := program.Statements[i]

		notEq(t, nil, stmt, "Got null statement at", strconv.Itoa(i))
		eq(t, true, testLetStatement(t, stmt, expectedIdentifier))
	}
}

func Test_ReturnStatement(t *testing.T) {
	input := `return 5;
  return 10;

  return add(5 ,10);`

	l := lexer.New_V2(strings.NewReader(input))
	p := New(l)

	program := p.ParseProgram()

	eq(t, true, checkParserErrs(t, p))
	notEq(t, nil, program, "ParseProgram() returned nil")
	eq(t, 3, len(program.Statements), "program.Statement does not contain 3 statements")

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		eq(t, true, ok, "Failed to typecast statement to let statement")
		eq(t, "return", returnStmt.TokenLiteral(), "Token literal of return didn't match")
	}
}

func Test_IdentifierExpression(t *testing.T) {
	l := lexer.New_V2(strings.NewReader(`foobar;`))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)

	eq(t, 1, len(program.Statements), "Expected 1 statement in the program")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed at typecasting program.Statement[0] to *ast.ExpressionStatement")

	ident, ok := stmt.Expression.(*ast.Identifier)
	eq(t, true, ok, "Failed at typecasting stmt.Epxression to *ast.Identifier")

	eq(t, "foobar", ident.Value, "Identifier value mis-match")
	eq(t, "foobar", ident.TokenLiteral(), "Identifier token literal mis-match")
}

func Test_IntegerLiteralExpression(t *testing.T) {
	l := lexer.New_V2(strings.NewReader(`5;`))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)

	eq(t, 1, len(program.Statements), "Expected 1 statement in the program")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed at typecasting program.Statement[0] to *ast.ExpressionStatement")

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	eq(t, true, ok, "Failed at typecasting stmt.Epxression to *ast.IntegerLiteral")

	eq(t, 5, ident.Value, "Identifier value mis-match")
	eq(t, "5", ident.TokenLiteral(), "Identifier token literal mis-match")
}

func Test_PrefixExpression(t *testing.T) {
	for _, test := range []struct {
		name         string
		input        string
		operator     string
		integerValue int64
	}{
		{
			name:         "test for ! prefix expression",
			input:        "!5",
			operator:     "!",
			integerValue: 5,
		},
		{
			name:         "test for - prefix expression",
			input:        "-5",
			operator:     "-",
			integerValue: 5,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)

			eq(t, 1, len(program.Statements), "Expecting 1 statements")

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			eq(t, true, ok, "Failed while typecasting program.Statement[0] to Expression Statement")

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			eq(t, true, ok, "Failed while typecasting stmt to PrefixExpression")
			eq(t, test.operator, exp.Operator, "Expression operator is not as expected")
			eq(t, true, testIntegerLiteral(t, exp.Right, test.integerValue), "Interger Literal test failed")
		})
	}
}
