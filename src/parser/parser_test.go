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
	eq(t, ok, true, "Failed to typecast il to *ast.IntegerLiteral")
	eq(t, value, integ.Value, "integ.Value doesn't match expected")
	eq(t, fmt.Sprintf("%d", value), integ.TokenLiteral(), "integ.TokenLiteral() doesn't match expected")
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	eq(t, ok, true, "Failed to typecast exp to *ast.Identifier")
	eq(t, value, ident.Value, "ident.Value doesn't match expected")
	eq(t, value, ident.TokenLiteral(), "ident.TokenLiteral() doesn't match expected")
	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	eq(t, ok, true, "Failed to typecast exp to *ast.Identifier")
	eq(t, value, boolean.Value, "ident.Value doesn't match expected")
	eq(t, fmt.Sprintf("%t", value), boolean.TokenLiteral(), "ident.TokenLiteral() doesn't match expected")
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case bool:
		return testBoolean(t, exp, bool(v))
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, int64(v))
	case string:
		return testIdentifier(t, exp, string(v))
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	eq(t, ok, true, "Failed to typecast exp to *ast.InfixExpression")
	eq(t, true, testLiteralExpression(t, opExp.Left, left))
	eq(t, operator, opExp.Operator, "Operator doesn't match expected")
	eq(t, true, testLiteralExpression(t, opExp.Right, right))
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

func Test_LetStatement(t *testing.T) {
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
		stmt := program.Statements[i].(*ast.LetStatement)

		notEq(t, nil, stmt, "Got null statement at", strconv.Itoa(i))
		eq(t, true, testLetStatement(t, stmt, expectedIdentifier))
	}
}

func Test_LetStatements(t *testing.T) {
	for _, test := range []struct {
		input         string
		expectedId    string
		expectedValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	} {
		t.Run(fmt.Sprintf("Test ran for input %s", test.input), func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)

			eq(t, 1, len(program.Statements), "Not found expected number of statements")

			stmt := program.Statements[0]
			eq(t, true, testLetStatement(t, stmt, test.expectedId))

			val := stmt.(*ast.LetStatement).Value
			eq(t, true, testLiteralExpression(t, val, test.expectedValue))
		})
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

func Test_ReturnStatements(t *testing.T) {
	for _, test := range []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	} {
		t.Run(fmt.Sprintf("Test ran for input %s", test.input), func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)

			eq(t, 1, len(program.Statements), "Not found expected number of statements")

			val := program.Statements[0].(*ast.ReturnStatement).ReturnValue
			eq(t, true, testLiteralExpression(t, val, test.expectedValue))
		})
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
		name     string
		input    string
		operator string
		value    interface{}
	}{
		{
			name:     "test for ! prefix expression",
			input:    "!5",
			operator: "!",
			value:    5,
		},
		{
			name:     "test for - prefix expression",
			input:    "-5;",
			operator: "-",
			value:    5,
		},
		{
			name:     "test for ! prefix expression for true",
			input:    "!true;",
			operator: "!",
			value:    true,
		},
		{
			name:     "test for ! prefix expression for false",
			input:    "!false;",
			operator: "!",
			value:    false,
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
			eq(t, true, testLiteralExpression(t, exp.Right, test.value), "Literal test failed")
		})
	}
}

func Test_InfixOperators(t *testing.T) {
	for _, test := range []struct {
		name     string
		input    string
		left     int64
		operator string
		right    int64
	}{
		{
			name:     "test sum operator",
			input:    "5 + 5;",
			left:     5,
			operator: "+",
			right:    5,
		},
		{
			name:     "test minus operator",
			input:    "5 - 5;",
			left:     5,
			operator: "-",
			right:    5,
		},
		{
			name:     "test product operator",
			input:    "5 * 5;",
			left:     5,
			operator: "*",
			right:    5,
		},
		{
			name:     "test divide operator",
			input:    "5 / 5;",
			left:     5,
			operator: "/",
			right:    5,
		},
		{
			name:     "test greater operator",
			input:    "5 > 5;",
			left:     5,
			operator: ">",
			right:    5,
		},
		{
			name:     "test lesser operator",
			input:    "5 < 5;",
			left:     5,
			operator: "<",
			right:    5,
		},
		{
			name:     "test greater than operator",
			input:    "5 >= 5;",
			left:     5,
			operator: ">=",
			right:    5,
		},
		{
			name:     "test lesser than operator",
			input:    "5 <= 5;",
			left:     5,
			operator: "<=",
			right:    5,
		},
		{
			name:     "test equal operator",
			input:    "5 == 5;",
			left:     5,
			operator: "==",
			right:    5,
		},
		{
			name:     "test not equal operator",
			input:    "5 != 5;",
			left:     5,
			operator: "!=",
			right:    5,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)

			eq(t, 1, len(program.Statements), "Expected 1 program statement")

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			eq(t, true, ok, "Failed to typecase program.Statement[0] to *ast.ExpressionStatement")

			exp, ok := stmt.Expression.(*ast.InfixExpression)
			eq(t, true, testInfixExpression(t, exp, test.left, test.operator, test.right))
		})
	}
}

func Test_OperatorPrecedenceParsing(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "(a * (b / c))"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 <= 4 != 3 >= 4", "((5 <= 4) != (3 >= 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + (c * (d / f))) + g))"},
	} {
		t.Run(fmt.Sprintf("Test %s to give %s", test.input, test.expected), func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)
			eq(t, test.expected, program.String(), "Failed to match expected string")
		})
	}
}

func Test_Boolean(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false", false},
	} {
		t.Run(fmt.Sprintf("test for %s boolean", test.input), func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)
			eq(t, 1, len(program.Statements), "Expected 1 statement in the program")

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			eq(t, true, ok, "Failed at typecasting program.Statement[0] to *ast.ExpressionStatement")
			eq(t, true, testBoolean(t, stmt.Expression, test.expected))
		})
	}
}

func Test_IfExpression(t *testing.T) {
	l := lexer.New_V2(strings.NewReader("if (x < y) { x }"))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)

	eq(t, 1, len(program.Statements), "Expected 1 if statement in program")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecase program.Statements[0] to *ast.ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	eq(t, true, ok, "Failed to typecase stmt.Expression to *ast.IfExpression")
	eq(t, true, testInfixExpression(t, exp.Condition, "x", "<", "y"))
	notEq(t, nil, exp.Consequence, "Got nil consequence")
	eq(t, 1, len(exp.Consequence.Statements), "Expected 1 statement in consequence")

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecase exp.Consequence.Statements[0] to *ast.ExpressionStatement")
	eq(t, true, testIdentifier(t, consequence.Expression, "x"))

	eq(t, nil, exp.Alternative, "Got non nil alternative")
}

func Test_IfElseExpression(t *testing.T) {
	l := lexer.New_V2(strings.NewReader("if (x < y) { x } else { y }"))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)

	eq(t, 1, len(program.Statements), "Expected 1 if statement in program")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecase program.Statements[0] to *ast.ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	eq(t, true, ok, "Failed to typecase stmt.Expression to *ast.IfExpression")
	eq(t, true, testInfixExpression(t, exp.Condition, "x", "<", "y"))
	notEq(t, nil, exp.Consequence, "Got nil consequence")
	eq(t, 1, len(exp.Consequence.Statements), "Expected 1 statement in consequence")

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecase exp.Consequence.Statements[0] to *ast.ExpressionStatement")
	eq(t, true, testIdentifier(t, consequence.Expression, "x"))

	notEq(t, nil, exp.Alternative, "Got nil alternative")
	eq(t, 1, len(exp.Alternative.Statements), "Expected 1 statement in alternative")

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecase exp.Alternative.Statements[0] to *ast.ExpressionStatement")
	eq(t, true, testIdentifier(t, alternative.Expression, "y"))
}

func Test_FunctionLiteralParsing(t *testing.T) {
	l := lexer.New_V2(strings.NewReader("fn(x, y) { x + y }"))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)
	eq(t, 1, len(program.Statements), "Expected 1 program statement")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecast program.Statements[0] as *ast.ExpressionStatement")

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	eq(t, true, ok, "Failed to typecast stmt.Expression as *ast.FunctionLiteral")
	eq(t, 2, len(function.Parameters), "Expected 2 params")
	eq(t, true, testLiteralExpression(t, function.Parameters[0], "x"))
	eq(t, true, testLiteralExpression(t, function.Parameters[1], "y"))
	eq(t, 1, len(function.Body.Statements), "Expected 1 function body statements")

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecast function.Body.Statements[0] as *ast.ExpressionStatement")
	eq(t, true, testInfixExpression(t, bodyStmt.Expression, "x", "+", "y"))
}

func Test_FunctionParameterParsing(t *testing.T) {
	for _, test := range []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {}", []string{}},
		{"fn(x) {}", []string{"x"}},
		{"fn(x,y) {}", []string{"x", "y"}},
	} {
		t.Run(fmt.Sprintf("Test params for fun : %s", test.input), func(t *testing.T) {
			l := lexer.New_V2(strings.NewReader(test.input))
			p := New(l)
			program := p.ParseProgram()

			checkParserErrs(t, p)
			eq(t, 1, len(program.Statements), "Expected 1 program statement")

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			eq(t, true, ok, "Failed to typecast program.Statements[0] as *ast.ExpressionStatement")

			function, ok := stmt.Expression.(*ast.FunctionLiteral)
			eq(t, true, ok, "Failed to typecast stmt.Expression as *ast.FunctionLiteral")

			eq(t, len(test.expectedParams), len(function.Parameters), "Mismatch number of expected params")
			for i, expectedParam := range test.expectedParams {
				eq(t, true, testLiteralExpression(t, function.Parameters[i], expectedParam))
			}
		})
	}
}

func Test_CallExpression(t *testing.T) {
	l := lexer.New_V2(strings.NewReader("add(1, 2 * 3, 4 + 5);"))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)

	eq(t, 1, len(program.Statements), "Expected 1 statement in program")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	eq(t, true, ok, "Failed to typecast program.Statement[0] to *ast.ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.CallExpression)
	eq(t, true, ok, "Failed to typecast stmt.Expression to *ast.CallExpression")
	eq(t, true, testIdentifier(t, exp.Function, "add"))
	eq(t, true, testLiteralExpression(t, exp.Arguments[0], 1))
	eq(t, true, testInfixExpression(t, exp.Arguments[1], 2, "*", 3))
	eq(t, true, testInfixExpression(t, exp.Arguments[2], 4, "+", 5))
}

func Test_Assignment(t *testing.T) {
	l := lexer.New_V2(strings.NewReader("a = a + b * 5;"))
	p := New(l)
	program := p.ParseProgram()

	checkParserErrs(t, p)

	eq(t, 1, len(program.Statements), "Expected 1 statement in the program")

	stmt, ok := program.Statements[0].(*ast.Assignment)
	eq(t, true, ok, "Failed to typecast program.Statements[0] to *ast.Assignment")
	notEq(t, nil, stmt.Identifier, "Found nil identifier")
	notEq(t, nil, &stmt.Value, "Found nil value")
	eq(t, "a = (a + (b * 5))", stmt.String(), "Invalid assignment")
}

func Test_AssignmentErr(t *testing.T) {
	l := lexer.New_V2(strings.NewReader("a = ;"))
	p := New(l)
	p.ParseProgram()

	eq(t, 2, len(p.Errors()), "Expected 2 error")
	eq(t, "No prefix parser function for ; found", p.Errors()[0].Error(), "1st Err msg didn't match")
	eq(t, "Got empty expression on RHS of assignment", p.Errors()[1].Error(), "2nd Err msg didn't match")
}
