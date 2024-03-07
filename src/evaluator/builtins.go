package evaluator

import (
	"fmt"

	"sudocoding.xyz/interpreter_in_go/src/object"
)

var builtins = map[string]*object.Builtin{
	`len`: {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d. want=1", len(args))
		}

		strObj, ok := args[0].(*object.String)
		if !ok {
			return newError("argument to `len` not supported. got %s", args[0].Type())
		}

		return &object.Integer{Value: int64(len(strObj.Value))}
	}},

	`print`: {Fn: func(args ...object.Object) object.Object {
		for _, item := range args {
			switch item := item.(type) {
			case *object.Integer:
				fmt.Print(item.Value)
			case *object.String:
				fmt.Print(item.Value)
			case *object.Boolean:
				fmt.Print(item.Value)
			case *object.Error:
				fmt.Print(item.Message)
			case *object.Null:
				fmt.Print(item.Inspect())
			}
		}

		fmt.Println()
		return NULL
	}},
}
