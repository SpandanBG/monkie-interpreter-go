package evaluator

import (
	"fmt"

	"sudocoding.xyz/interpreter_in_go/src/object"
)

var builtins = map[string]*object.Builtin{
	`len`: {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` not supported. got %s", args[0].Type())
		}
	}},

	`first`: {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		if args[0].Type() != object.ARRAY_OBJ {
			return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array).Elements
		if len(arr) > 0 {
			return arr[0]
		}
		return NULL
	}},

	`last`: {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
		}

		length := len(arr.Elements)
		if length > 0 {
			return arr.Elements[length-1]
		}

		return NULL
	}},

	`rest`: {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}

		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
		}

		length := len(arr.Elements)
		if length > 0 {
			newElements := make([]object.Object, length-1, length-1)
			copy(newElements, arr.Elements[1:length])
			return &object.Array{Elements: newElements}
		}

		return NULL
	}},

	`push`: {Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=2", len(args))
		}

		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
		}

		arr.Elements = append(arr.Elements, args[1])
		return NULL
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
