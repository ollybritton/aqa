package repl

import (
	"fmt"

	au "github.com/logrusorgru/aurora"
)

// Errors prints a list of errors
func Errors(errs []error) {

	fmt.Println(au.Red(au.Bold("Fatal error(s) occured while parsing the input:")))
	for _, err := range errs {
		errType := fmt.Sprintf("%T", err)
		fmt.Printf("* %v: %v\n", au.Italic(au.Yellow(errType)), au.Green(err.Error()))
	}

	fmt.Println("")
}
