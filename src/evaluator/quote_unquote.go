package evaluator

import (
	"sudocoding.xyz/interpreter_in_go/src/object"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
)

const (
	QUOTE_LITERAL = "quote"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
