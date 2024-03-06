package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/object"
	"sudocoding.xyz/interpreter_in_go/src/parser"
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

func testEval(input string) object.Object {
	l := lexer.New_V2(strings.NewReader(input))
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObj(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	eq(t, true, ok, "Failed to typecast obj to object.Integer")
	eq(t, expected, result.Value, "Expected int64 didn't match")
	return true
}

func testBooleanObj(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	eq(t, true, ok, "Failed to typecast obj to object.Boolean")
	eq(t, expected, result.Value, "Expected boolean didn't match")
	return true
}

func testNullObj(t *testing.T, obj object.Object) bool {
	result, ok := obj.(*object.Null)
	eq(t, true, ok, "Failed to typecast obj to object.Null")
	eq(t, NULL, result, "Expected NULL. Didn't match")
	return true
}

func Test_EvalIntegerExpression(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	} {
		t.Run(fmt.Sprintf("Tests for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)
			eq(t, true, testIntegerObj(t, evaluated, test.expected), "Did not match expected")
		})
	}
}

func Test_EvalBooleanExpression(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"true == true", true},
		{"true != true", false},
		{"false == false", true},
		{"false != false", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(5 * 5) == 25 == true", true},
	} {
		t.Run(fmt.Sprintf("Tests for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)
			eq(t, true, testBooleanObj(t, evaluated, test.expected), "Did not match expected")
		})
	}
}

func Test_BangOperator(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!5", true},
	} {
		t.Run(fmt.Sprintf("Test for input %s", test.input), func(t *testing.T) {
			eq(t, true, testBooleanObj(t, testEval(test.input), test.expected))
		})
	}
}

func Test_IfElseExpression(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	} {
		t.Run(fmt.Sprintf("Test if-else for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)

			if integer, ok := test.expected.(int); ok {
				eq(t, true, testIntegerObj(t, evaluated, int64(integer)))
			} else {
				eq(t, true, testNullObj(t, evaluated))
			}
		})
	}
}

func Test_ReturnStatement(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int16
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	} {
		t.Run(fmt.Sprintf("Test return statement for %s", test.input), func(t *testing.T) {
			eq(t, true, testIntegerObj(t, testEval(test.input), int64(test.expected)))
		})
	}
}

func Test_ErrorHandling(t *testing.T) {
	for _, test := range []struct {
		input  string
		errMsg string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false; } return 1; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	} {
		t.Run(fmt.Sprintf("Test error for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)

			errObj, ok := evaluated.(*object.Error)
			eq(t, true, ok, "Failed to typecast evaulated to *object.Error")
			eq(t, test.errMsg, errObj.Message, "Error message mismatch")
		})
	}
}

func Test_LetStatements(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	} {
		t.Run(fmt.Sprintf("Test let statement for %s", test.input), func(t *testing.T) {
			eq(t, true, testIntegerObj(t, testEval(test.input), test.expected))
		})
	}
}
