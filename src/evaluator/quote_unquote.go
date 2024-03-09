package evaluator

import (
	"fmt"

	"sudocoding.xyz/interpreter_in_go/src/object"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
	"sudocoding.xyz/interpreter_in_go/src/token"
)

const (
	QUOTE_LITERAL   = "quote"
	UNQUOTE_LITERAL = "unquote"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquotedCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToAstNode(unquoted)
	})
}

func isUnquotedCall(node ast.Node) bool {
	if callexp, ok := node.(*ast.CallExpression); ok {
		return callexp.Function.TokenLiteral() == UNQUOTE_LITERAL
	}
	return false
}

func convertObjectToAstNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		return &ast.IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: fmt.Sprintf("%d", obj.Value)},
			Value: obj.Value,
		}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	}

	return nil
}
