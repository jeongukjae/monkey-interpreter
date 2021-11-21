package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			arg, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to len not supported, got %s", args[0].Type())
			}
			return &object.Integer{Value: int64(len(arg.Value))}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.ObjectTypeObject{Value: args[0].Type()}
		},
	},
}
