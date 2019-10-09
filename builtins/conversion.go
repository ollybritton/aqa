package builtins

import (
	"fmt"
	"strconv"

	"github.com/ollybritton/aqa/object"
)

// BuiltinStringToInt will convert a string object into an integer object.
func BuiltinStringToInt(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return newError("argument to `STRING_TO_INT` not supported, got=%s", args[0].Type())
	}

	conv, err := strconv.ParseInt(str.Value, 0, 64)
	if err != nil {
		return newError("failed to convert %q to integer in call to `STRING_TO_INT`", str.Value)
	}

	return &object.Integer{Value: conv}
}

// BuiltinIntToString will convert an integer object into a string object.
func BuiltinIntToString(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	i, ok := args[0].(*object.Integer)
	if !ok {
		return newError("argument to `INT_TO_STRING` not supported, got=%s", args[0].Type())
	}

	conv := strconv.FormatInt(i.Value, 10)
	return &object.String{Value: conv}
}

// BuiltinStringToReal converts a string object into a real (floating point) object.
func BuiltinStringToReal(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return newError("argument to `STRING_TO_REAL` not supported, got=%s", args[0].Type())
	}

	conv, err := strconv.ParseFloat(str.Value, 64)
	if err != nil {
		return newError("failed to convert %q to integer in call to `STRING_TO_REAL`", str.Value)
	}

	return &object.Float{Value: conv}
}

// BuiltinRealToString converts a real/float to a string object.
func BuiltinRealToString(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	f, ok := args[0].(*object.Float)
	if !ok {
		return newError("argument to `REAL_TO_STRING` not supported, got=%s", args[0].Type())
	}

	conv := fmt.Sprintf("%f", f.Value)
	return &object.String{Value: conv}
}

// BuiltinCharToCode converts a character into its ASCII equivalent integer.
func BuiltinCharToCode(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	str, ok := args[0].(*object.String)
	if !ok {
		return newError("argument to `CHAR_TO_CODE` not supported, got=%s", args[0].Type())
	}

	if len(str.Value) != 1 {
		return newError("argument to `CHAR_TO_CODE` not supported, cannot convert multiple characters, got=%s", str.Value)
	}

	return &object.Integer{Value: int64(str.Value[0])}
}

// BuiltinCodeToChar converts an integer ascii code into a character.
func BuiltinCodeToChar(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	i, ok := args[0].(*object.Integer)
	if !ok {
		return newError("argument to `CODE_TO_CHAR` not supported, got=%s", args[0].Type())
	}

	return &object.String{Value: string(byte(i.Value))}
}
