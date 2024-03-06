package evaluator

import (
	"fmt"

	"sudocoding.xyz/interpreter_in_go/src/object"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
	"sudocoding.xyz/interpreter_in_go/src/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		value, ok := expectEval(node.ReturnValue)
		if !ok {
			return value
		}

		return &object.ReturnValue{Value: value}

		// Expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.PrefixExpression:
		right, ok := expectEval(node.Right)
		if !ok {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, ok := expectEval(node.Left)
		if !ok {
			return left
		}

		right, ok := expectEval(node.Right)
		if !ok {
			return right
		}

		return evalInfixExpression(left, node.Operator, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return NULL
}

func expectEval(node ast.Node) (object.Object, bool) {
	evaluated := Eval(node)
	_, ok := evaluated.(*object.Error)
	return evaluated, !ok
}

func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		// Stop processing block if return statement reached
		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExp(right)
	case "-":
		return evalMinusPrefixOpExp(right)
	}

	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalBangOperatorExp(right object.Object) object.Object {
	if right == FALSE {
		return TRUE
	}

	if right == FALSE || right == NULL {
		return TRUE
	}

	if intLit, ok := right.(*object.Integer); ok && intLit.Value == 0 {
		return TRUE
	}

	return FALSE
}

func evalMinusPrefixOpExp(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(left, operator, right)
	}

	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(left, operator, right)
	}

	return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	lVal := left.(*object.Integer).Value
	rVal := right.(*object.Integer).Value

	switch token.TokenType(operator) {
	case token.PLUS:
		return &object.Integer{Value: lVal + rVal}
	case token.MINUS:
		return &object.Integer{Value: lVal - rVal}
	case token.ASTERISK:
		return &object.Integer{Value: lVal * rVal}
	case token.SLASH:
		return &object.Integer{Value: lVal / rVal}
	case token.EQ:
		return nativeBoolToBooleanObj(lVal == rVal)
	case token.NOT_EQ:
		return nativeBoolToBooleanObj(lVal != rVal)
	case token.GT:
		return nativeBoolToBooleanObj(lVal > rVal)
	case token.LT:
		return nativeBoolToBooleanObj(lVal < rVal)
	case token.GTE:
		return nativeBoolToBooleanObj(lVal >= rVal)
	case token.LTE:
		return nativeBoolToBooleanObj(lVal <= rVal)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalBooleanInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	lVal := left.(*object.Boolean).Value
	rVal := right.(*object.Boolean).Value

	switch token.TokenType(operator) {
	case token.EQ:
		return nativeBoolToBooleanObj(lVal == rVal)
	case token.NOT_EQ:
		return nativeBoolToBooleanObj(lVal != rVal)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition, ok := expectEval(ie.Condition)
	if !ok {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return NULL
}

func nativeBoolToBooleanObj(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case NULL:
		return false
	case FALSE:
		return false
	}

	if integer, ok := obj.(*object.Integer); ok && integer.Value == 0 {
		return false
	}
	return true
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
