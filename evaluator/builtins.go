package evaluator

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/ollybritton/aqa/object"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var builtins = map[string]*object.Builtin{
	"LEN": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `LEN` not supported, got=%s", args[0].Type())
			}
		},
	},
	"POSITION": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
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

			return NULL
		},
	},
	"SUBSTRING": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
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
		},
	},

	"STRING_TO_INT": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
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
		},
	},
	"INT_TO_STRING": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			i, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `INT_TO_STRING` not supported, got=%s", args[0].Type())
			}

			conv := strconv.FormatInt(i.Value, 10)
			return &object.String{Value: conv}
		},
	},
	"CHAR_TO_CODE": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
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
		},
	},
	"CODE_TO_CHAR": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			i, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `CODE_TO_CHAR` not supported, got=%s", args[0].Type())
			}

			return &object.String{Value: string(byte(i.Value))}
		},
	},

	"RANDOM_INT": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
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
		},
	},

	"OUTPUT": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			for _, a := range args {
				fmt.Printf("%s ", a.Inspect())
			}

			fmt.Print("\n")
			return NULL
		},
	},

	"FLOOR": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
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
		},
	},
}
