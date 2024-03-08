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

func testStringObj(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	eq(t, true, ok, "Failed to typecast obj to object.String")
	eq(t, expected, result.Value, "Expected string didn't match")
	return true
}

func testErrorObj(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Error)
	eq(t, true, ok, "Failed to typecast obj to object.Error")
	eq(t, expected, result.Message, "Expected string didn't match")
	return true
}

func testArrayObj(t *testing.T, obj object.Object, expected []interface{}) bool {
	result, ok := obj.(*object.Array)
	eq(t, true, ok, "Failed to typecast obj to object.Array")

	for i, exp := range expected {
		switch exp := exp.(type) {
		case int:
			eq(t, true, testIntegerObj(t, result.Elements[i], int64(exp)))
		case string:
			eq(t, true, testStringObj(t, result.Elements[i], exp))
		case bool:
			eq(t, true, testBooleanObj(t, result.Elements[i], exp))
		default:
			eq(t, true, testNullObj(t, result.Elements[i]))
		}
	}

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
		{"\"asdf\" == \"qwer\"", false},
		{"\"asdf\" != \"qwer\"", true},
	} {
		t.Run(fmt.Sprintf("Tests for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)
			eq(t, true, testBooleanObj(t, evaluated, test.expected), "Did not match expected")
		})
	}
}

func Test_EvalStringExpression(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected string
	}{
		{"\"asdf\"", "asdf"},
		{"\"asdf\" + \"qwer\"", "asdfqwer"},
	} {
		t.Run(fmt.Sprintf("Test string expresson for: %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)
			eq(t, true, testStringObj(t, evaluated, test.expected), "Did not match expected")
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
		{"if (1 > 2) { \"a\" } else { \"b\" }", "b"},
	} {
		t.Run(fmt.Sprintf("Test if-else for %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)

			switch expected := test.expected.(type) {
			case int:
				eq(t, true, testIntegerObj(t, evaluated, int64(expected)))
			case string:
				eq(t, true, testStringObj(t, evaluated, expected))
			default:
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
		{"\"asdf\" + true", "type mismatch: STRING + BOOLEAN"},
		{"\"asdf\" + 5", "type mismatch: STRING + INTEGER"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"\"asdf\" - \"asdf\";", "unknown operator: STRING - STRING"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false; } return 1; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{"a = 5;", "variable a hasn't been initialized"},
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
		expected interface{}
	}{
		{"let a = 5;", nil},
		{"let a = 5; a;", 5},
		{"let a = \"asdf\"; a;", "asdf"},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	} {
		t.Run(fmt.Sprintf("Test let statement for %s", test.input), func(t *testing.T) {
			switch expected := test.expected.(type) {
			case int:
				eq(t, true, testIntegerObj(t, testEval(test.input), int64(expected)))
			case string:
				eq(t, true, testStringObj(t, testEval(test.input), expected))
			default:
				eq(t, true, testNullObj(t, testEval(test.input)))
			}
		})
	}
}

func Test_AssigmentStatements(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a = a + 5; a", 10},
		{"let a = 5 * 5; a = 10; a", 10},
		{"let a = 5; let b = a + 10; a = b; a", 15},
		{"let a = 5; let b = a; let c = 0; c = a + b + 5; c;", 15},
	} {
		t.Run(fmt.Sprintf("Test let statement for %s", test.input), func(t *testing.T) {
			eq(t, true, testIntegerObj(t, testEval(test.input), test.expected))
		})
	}
}

func Test_FunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	eq(t, true, ok, "Failed to typecast evaluated to *object.Function")
	eq(t, 1, len(fn.Parameters), "Expected 1 param in function def")
	eq(t, "x", fn.Parameters[0].String(), "Function param name didn't match")
	eq(t, "(x + 2)", fn.Body.String(), "Function body didn't match")
}

func Test_FunctionApplication(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int64
	}{
		{"let id = fn(x) { x }; id(5);", 5},
		{"let id = fn(x) { return x; }; id(5);", 5},
		{"let double = fn(x) { x * 2; }; double(2);", 4},
		{"let add = fn(a, b) { a + b }; add(1 + 1, add(2, 3))", 7},
		{"fn(x) { x }(5)", 5},
	} {
		t.Run(fmt.Sprintf("Test function call %s", test.input), func(t *testing.T) {
			eq(t, true, testIntegerObj(t, testEval(test.input), test.expected))
		})
	}
}

func Test_Closure(t *testing.T) {
	input := `
    let add = fn(x) { fn(y) { x + y } };
    let addTwo = add(2);
    addTwo(4);
  `
	eq(t, true, testIntegerObj(t, testEval(input), 6))
}

func Test_BuiltinFunction(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("asdf")`, 4},
		{`len(1)`, "argument to `len` not supported. got INTEGER"},
		{`len("asdf", "asdf")`, "wrong number of arguments. got=2, want=1"},
		{`len([1 ,2, 3])`, 3},
		{`first([1, 2, 3])`, 1},
		{`first([1], [2, 3])`, "wrong number of arguments. got=2, want=1"},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`last([1])`, 1},
		{`last([1, 2])`, 2},
		{`last([1], [2, 3])`, "wrong number of arguments. got=2, want=1"},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`rest([1])`, []interface{}{}},
		{`rest([1, 2, 3])`, []interface{}{2, 3}},
		{`rest([1], [2, 3])`, "wrong number of arguments. got=2, want=1"},
		{`rest(1)`, "argument to `rest` must be ARRAY, got INTEGER"},
		{`let a = [1, 2]; push(a, 3); a`, []interface{}{1, 2, 3}},
		{`push([1])`, "wrong number of arguments. got=1, want=2"},
		{`push(1, 1)`, "first argument to `push` must be ARRAY, got INTEGER"},
	} {
		t.Run(fmt.Sprintf("Test built in fn: %s", test.input), func(t *testing.T) {
			evaluated := testEval(test.input)
			switch expected := test.expected.(type) {
			case int:
				eq(t, true, testIntegerObj(t, evaluated, int64(expected)))
			case string:
				if _, ok := evaluated.(*object.Error); ok {
					eq(t, true, testErrorObj(t, evaluated, expected))
				} else {
					eq(t, true, testStringObj(t, evaluated, expected))
				}
			case bool:
				eq(t, true, testBooleanObj(t, evaluated, expected))
			case []interface{}:
				eq(t, true, testArrayObj(t, evaluated, expected))
			default:
				eq(t, true, testNullObj(t, evaluated))
			}
		})
	}
}

func Test_ArrayLiterals(t *testing.T) {
	input := `[1, 2 * 3, "asdf", true]`

	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)
	eq(t, true, ok, "Failed to typecast to evaluated to *object.Array")
	eq(t, 4, len(result.Elements), "Expected 4 elements in the array")
	eq(t, true, testIntegerObj(t, result.Elements[0], 1))
	eq(t, true, testIntegerObj(t, result.Elements[1], 6))
	eq(t, true, testStringObj(t, result.Elements[2], "asdf"))
	eq(t, true, testBooleanObj(t, result.Elements[3], true))
}

func Test_IndexingExpression(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected interface{}
	}{
		{"let a = [5]; a[0]", 5},
		{"let a = [5]; a[1]", nil},
		{"let a = [5]; a[-1]", nil},
		{"let a = [1, 2, 3]; a[1 + 1]", 3},
		{"[1, 2, 3][1]", 2},
		{"let i=0; let a=[4]; a[i]", 4},
		{"let a=[1,2,3]; a[0] + a[1] + a[2]", 6},
		{`let a=[1, "a", true]; a[1]`, "a"},
		{`let a=[1, "a", true]; a[2]`, true},
	} {
		evaluated := testEval(test.input)
		switch expected := test.expected.(type) {
		case int:
			eq(t, true, testIntegerObj(t, evaluated, int64(expected)))
		case string:
			eq(t, true, testStringObj(t, evaluated, expected))
		case bool:
			eq(t, true, testBooleanObj(t, evaluated, expected))
		default:
			eq(t, true, testNullObj(t, evaluated))
		}

	}
}
