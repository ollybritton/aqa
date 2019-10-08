package builtins

import (
	"os"

	"github.com/ollybritton/aqa/object"
)

// BuiltinExit will exit the program. If there are no arguments specified, it exits with a code of 0.
// If there is one argument, it exits with that code. If there are more, it errors.
func BuiltinExit(args ...object.Object) object.Object {
	switch len(args) {
	case 0:
		os.Exit(0)
		return &object.Null{}

	case 1:
		arg, ok := args[0].(*object.Integer)
		if !ok {
			return newError("argument to EXIT not supported: %s", args[0].Type())
		}

		os.Exit(int(arg.Value))
		return &object.Null{}

	default:
		return newError("wrong number of arguments. got=%d, want=0|1", len(args))
	}
}
