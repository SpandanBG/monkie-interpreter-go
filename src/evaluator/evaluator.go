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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		if value, ok := expectEval(node.ReturnValue, env); ok {
			return &object.ReturnValue{Value: value}
		} else {
			return value
		}
	case *ast.LetStatement:
		if value, ok := expectEval(node.Value, env); ok {
			env.Set(node.Name.Value, value)
		} else {
			return value
		}
	case *ast.Assignment:
		if _, ok := env.Get(node.Identifier.Value); !ok {
			return newError("variable %v hasn't been initialized", node.Identifier.Value)
		}
		if value, ok := expectEval(node.Value, env); ok {
			env.Set(node.Identifier.Value, value)
		} else {
			return value
		}

		// Expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node.Elements, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node.Pairs, env)
	case *ast.IndexExpression:
		left, ok := expectEval(node.Left, env)
		if !ok {
			return left
		}

		index, ok := expectEval(node.Index, env)
		if !ok {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		if right, ok := expectEval(node.Right, env); ok {
			return evalPrefixExpression(node.Operator, right)
		} else {
			return right
		}
	case *ast.InfixExpression:
		left, ok := expectEval(node.Left, env)
		if !ok {
			return left
		}

		right, ok := expectEval(node.Right, env)
		if !ok {
			return right
		}

		return evalInfixExpression(left, node.Operator, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}
	case *ast.CallExpression:
		if node.Function.TokenLiteral() == QUOTE_LITERAL {
			if len(node.Arguments) != 1 {
				return newError("wrong number of arguments. got=%d, want=2", len(node.Arguments))
			}
			return quote(node.Arguments[0])
		}

		fn, ok := expectEval(node.Function, env)
		if !ok {
			return fn
		}

		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return err
		}

		return applyFn(fn, args)
	}

	return NULL
}

func expectEval(node ast.Node, env *object.Environment) (object.Object, bool) {
	evaluated := Eval(node, env)
	_, ok := evaluated.(*object.Error)
	return evaluated, !ok
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		// Stop processing block if return statement reached
		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func evalExpressions(expressons []ast.Expression, env *object.Environment) ([]object.Object, object.Object) {
	var result []object.Object

	for _, e := range expressons {
		if evaluated, ok := expectEval(e, env); ok {
			result = append(result, evaluated)
		} else {
			return nil, evaluated
		}
	}

	return result, nil
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

	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(left, operator, right)
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

func evalStringInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	lVal := left.(*object.String).Value
	rVal := right.(*object.String).Value

	switch token.TokenType(operator) {
	case token.PLUS:
		return &object.String{Value: lVal + rVal}
	case token.EQ:
		return nativeBoolToBooleanObj(lVal == rVal)
	case token.NOT_EQ:
		return nativeBoolToBooleanObj(lVal != rVal)
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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition, ok := expectEval(ie.Condition, env)
	if !ok {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(id.Value); ok {
		return val
	}

	if fn, ok := builtins[id.Value]; ok {
		return fn
	}

	return newError("identifier not found: %s", id.Value)
}

func evalArrayLiteral(elements []ast.Expression, env *object.Environment) object.Object {
	elms, err := evalExpressions(elements, env)
	if err != nil {
		return err
	}

	return &object.Array{Elements: elms}
}

func evalHashLiteral(pairs map[ast.Expression]ast.Expression, env *object.Environment) object.Object {
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}

	for keyExp, valueExp := range pairs {
		key, ok := expectEval(keyExp, env)
		if !ok {
			return key
		}

		hashable, ok := key.(object.Hashable)
		if !ok {
			return newError("key of type %s is not hashable", key.Type())
		}

		value, ok := expectEval(valueExp, env)
		if !ok {
			return value
		}

		hash.Pairs[hashable.Hash()] = object.HashPair{Key: key, Value: value}
	}

	return hash
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExp(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExp(left, index)
	default:
		return newError("index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExp(left object.Object, index object.Object) object.Object {
	arr := left.(*object.Array)
	idx := index.(*object.Integer).Value

	max := int64(len(arr.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}

	return arr.Elements[idx]
}

func evalHashIndexExp(left object.Object, index object.Object) object.Object {
	hash := left.(*object.Hash)
	hashable, ok := index.(object.Hashable)
	if !ok {
		return newError("index of type %s cannot be used as hash index", index.Type())
	}

	if value, ok := hash.Pairs[hashable.Hash()]; ok {
		return value.Value
	}

	return NULL
}

func applyFn(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFuncEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	}

	return newError("not a function: %s", fn.Type())
}

func extendFuncEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnv(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	switch obj := obj.(type) {
	case *object.ReturnValue:
		return obj.Value
	}
	return obj
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
