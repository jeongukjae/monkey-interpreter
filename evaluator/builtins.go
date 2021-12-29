package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{}

func init() {
	// to avoid initialization loop
	builtins["len"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to len not supported, got %s", args[0].Type())
			}
		},
	}
	builtins["type"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.ObjectTypeObject{Value: args[0].Type()}
		},
	}
	builtins["puts"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	}
	builtins["first"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	}
	builtins["last"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	}
	builtins["rest"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length-1)
			if length > 0 {
				copy(newElements, arr.Elements[1:length])
			}
			return &object.Array{Elements: newElements}
		},
	}
	builtins["push"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	}
	builtins["map"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}
			if args[1].Type() != object.FUNCTION_OBJ {
				return newError("argument to second must be FUNCTION, got %s", args[1].Type())
			}

			arr := args[0].(*object.Array)
			fn := args[1].(*object.Function)

			length := len(arr.Elements)
			newElements := make([]object.Object, length)

			for index, element := range arr.Elements {
				newElements[index] = applyFunction(fn, []object.Object{element})
			}

			return &object.Array{Elements: newElements}
		},
	}
	builtins["reduce"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}
			if args[2].Type() != object.FUNCTION_OBJ {
				return newError("argument to second must be FUNCTION, got %s", args[1].Type())
			}

			accumulated := args[1]
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length == 0 {
				return accumulated
			}

			fn := args[2].(*object.Function)

			for _, element := range arr.Elements {
				accumulated = applyFunction(fn, []object.Object{accumulated, element})
			}

			return accumulated
		},
	}
}
