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

	return Eval(program)
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

func Test_EvalIntegerExpression(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
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
