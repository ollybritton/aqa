package builtins

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ollybritton/aqa/object"
)

// Builtins maps the name of a builtin function within the program to the actual function.
var Builtins = make(map[string]*object.Builtin)

func newError(message string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, args...)}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	Builtins["LEN"] = &object.Builtin{Fn: BuiltinLen}
	Builtins["POSITION"] = &object.Builtin{Fn: BuiltinPosition}
	Builtins["SUBSTRING"] = &object.Builtin{Fn: BuiltinSubstring}

	Builtins["STRING_TO_INT"] = &object.Builtin{Fn: BuiltinStringToInt}
	Builtins["INT_TO_STRING"] = &object.Builtin{Fn: BuiltinIntToString}
	Builtins["STRING_TO_REAL"] = &object.Builtin{Fn: BuiltinStringToReal}
	Builtins["REAL_TO_STRING"] = &object.Builtin{Fn: BuiltinRealToString}
	Builtins["CHAR_TO_CODE"] = &object.Builtin{Fn: BuiltinCharToCode}
	Builtins["CODE_TO_CHAR"] = &object.Builtin{Fn: BuiltinCodeToChar}

	Builtins["RANDOM_INT"] = &object.Builtin{Fn: BuiltinRandomInt}

	Builtins["OUTPUT"] = &object.Builtin{Fn: BuiltinOutput}
	Builtins["PRINT"] = &object.Builtin{Fn: BuiltinPrint}
	Builtins["INPUT"] = &object.Builtin{Fn: BuiltinInput}

	Builtins["FLOOR"] = &object.Builtin{Fn: BuiltinFloor}
	Builtins["CEIL"] = &object.Builtin{Fn: BuiltinCeil}
	Builtins["SQRT"] = &object.Builtin{Fn: BuiltinSqrt}
}
