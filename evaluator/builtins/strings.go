package builtins

import "github.com/ollybritton/aqa/object"

// BuiltinLen calculates the length of a string or array object.
func BuiltinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `LEN` not supported, got=%s", args[0].Type())
	}
}

// BuiltinPosition finds the first position of a given character within a string.
func BuiltinPosition(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}

	fnArgs := []string{}
	for _, a := range args {
		switch arg := a.(type) {
		case *object.String:
			fnArgs = append(fnArgs, arg.Value)
		default:
			return newError("argument to `POSITION` not supported, got=%s", a.Type())
		}
	}

	searchString := fnArgs[0]
	toFind := fnArgs[1]

	for i := 0; i < len(searchString); i++ {
		char := string(searchString[i])

		if char == toFind {
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
