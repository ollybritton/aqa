package builtins

import (
	"github.com/ollybritton/aqa/object"
)

// BuiltinLen calculates the length of a string or array object.
func BuiltinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `LEN` not supported, got=%s", args[0].Type())
	}
}

// BuiltinPosition finds the first position of a given character within a string.
func BuiltinPosition(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	switch args[0].(type) {
	case *object.String:
		return builtinPositionString(args...)
	case *object.Array:
		return builtinPositionArray(args...)
	default:
		return newError("argument to `POSITION` not supported, got=%s", args[0].Type())
	}
}

func builtinPositionString(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	typed := []string{}
	for _, a := range args {
		switch arg := a.(type) {
		case *object.String:
			typed = append(typed, arg.Value)
		default:
			return newError("argument to `POSITION` not supported, got=%s", a.Type())
		}
	}

	search := typed[0]
	find := typed[1]

	for i := 0; i < len(search); i++ {
		char := string(search[i])

		if char == find {
			return &object.Integer{Value: int64(i)}
		}
	}

	return &object.Null{}
}

func builtinPositionArray(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	search, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `POSITION` not supported, got=%s", args[0].Type())
	}

	find := args[1]

	for i := 0; i < len(search.Elements); i++ {
		val := search.Elements[i]

		// TODO: use hash here
		if val.Inspect() == find.Inspect() {
			return &object.Integer{Value: int64(i)}
		}
	}

	return &object.Null{}
}

// BuiltinSubstring will slice a string object.
func BuiltinSubstring(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got=%d, want=3", len(args))
	}

	start, ok := args[0].(*object.Integer)
	if !ok {
		return newError("argument 1 to `SUBSTRING` not supported, got=%s", args[0].Type())
	}

	end, ok := args[1].(*object.Integer)
	if !ok {
		return newError("argument 2 to `SUBSTRING` not supported, got=%s", args[1].Type())
	}

	str, ok := args[2].(*object.String)
	if !ok {
		return newError("argument 3 to `SUBSTRING` not supported, got=%s", args[2].Type())
	}

	if start.Value > end.Value || start.Value < 0 || end.Value < 0 || end.Value > int64(len(str.Value)) {
		return newError("invalid bounds [%d:%d] in call to SUBSTRING", start.Value, end.Value)
	}

	return &object.String{Value: str.Value[start.Value : end.Value+1]}
}

// BuiltinSlice will slice an array object.
func BuiltinSlice(args ...object.Object) object.Object {
	if len(args) != 3 {
		return newError("wrong number of arguments. got=%d, want=3", len(args))
	}

	start, ok := args[0].(*object.Integer)
	if !ok {
		return newError("argument 1 to `SUBSTRING` not supported, got=%s", args[0].Type())
	}

	end, ok := args[1].(*object.Integer)
	if !ok {
		return newError("argument 2 to `SUBSTRING` not supported, got=%s", args[1].Type())
	}

	array, ok := args[2].(*object.Array)
	if !ok {
		return newError("argument 3 to `SUBSTRING` not supported, got=%s", args[2].Type())
	}

	if start.Value > end.Value || start.Value < 0 || end.Value < 0 || end.Value > int64(len(array.Elements)) {
		return newError("invalid bounds [%d:%d] in call to SLICE", start.Value, end.Value)
	}

	return &object.Array{Elements: array.Elements[start.Value : end.Value+1]}
}
