package builtins

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ollybritton/aqa/object"
)

// BuiltinOutput will print a single argument to stdout. It is what the OUTPUT keyword uses under the hood.
func BuiltinOutput(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	for _, a := range args {
		fmt.Printf("%s ", a.Inspect())
	}

	fmt.Print("\n")
	return &object.Null{}
}

// BuiltinPrint will prints its arguments, seperated by spaces, to stdout.
func BuiltinPrint(args ...object.Object) object.Object {
	for _, a := range args {
		fmt.Printf("%s ", a.Inspect())
	}

	fmt.Print("\n")
	return &object.Null{}
}

// BuiltinInput will take input from the user. If an argument is given, it is the text before the prompt.
func BuiltinInput(args ...object.Object) object.Object {
	if len(args) > 1 {
		return newError("wrong number of arguments. got=%d, want=0 or 1", len(args))
	}

	if len(args) == 1 {
		prompt, ok := args[0].(*object.String)
		if !ok {
			return newError("argument to `INPUT` not supported, got=%s", args[0].Type())
		}

		fmt.Print(prompt.Value)
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	return &object.String{Value: text}
}
