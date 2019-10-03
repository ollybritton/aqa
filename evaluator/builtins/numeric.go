package builtins

import (
	"math"
	"math/rand"

	"github.com/ollybritton/aqa/object"
)

// BuiltinRandomInt generates a random integer object between two bounds.
func BuiltinRandomInt(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	lower, ok := args[0].(*object.Integer)
	if !ok {
		return newError("argument 1 to `RANDOM_INT` not supported, got=%s", args[0].Type())
	}

	upper, ok := args[1].(*object.Integer)
	if !ok {
		return newError("argument 2 to `RANDOM_INT` not supported, got=%s", args[1].Type())
	}

	val := rand.Intn(int(upper.Value-lower.Value+1)) + int(lower.Value)
	return &object.Integer{Value: int64(val)}
}

// BuiltinFloor will floor a float. It has no effect on integers.
func BuiltinFloor(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch val := args[0].(type) {
	case *object.Float:
		return &object.Integer{Value: int64(math.Floor(val.Value))}
	case *object.Integer:
		return val
	default:
		return newError("argument to `FLOOR` not supported, got=%s", args[0].Type())
	}
}

// BuiltinCeil will round a float up. It has no effect on integers.
func BuiltinCeil(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch val := args[0].(type) {
	case *object.Float:
		return &object.Integer{Value: int64(math.Ceil(val.Value))}
	case *object.Integer:
		return val
	default:
		return newError("argument to `FLOOR` not supported, got=%s", args[0].Type())
	}
}

// BuiltinSqrt will find the square root of an integer or a float.
func BuiltinSqrt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch val := args[0].(type) {
	case *object.Float:
		return &object.Float{Value: math.Sqrt(val.Value)}
	case *object.Integer:
		return &object.Float{Value: math.Sqrt(float64(val.Value))}
	default:
		return newError("argument to `FLOOR` not supported, got=%s", args[0].Type())
	}
}
