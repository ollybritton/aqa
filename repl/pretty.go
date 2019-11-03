package repl

import (
	"fmt"

	"github.com/ollybritton/aqa/token"

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

// PrettyToken will pretty-print a token.
func PrettyToken(t token.Token) string {
	var ttype string
	var literal = au.Bold(au.BrightYellow(t.Literal))

	if t.Type == token.ILLEGAL {
		ttype = fmt.Sprint(au.Red(t.Type))
	} else {
		ttype = fmt.Sprint(au.Green(au.Italic(t.Type)))
	}

	return fmt.Sprintf("(Lit: '%s', Type: '%s', line=%d, startcol=%d, endcol=%d)",
		literal, ttype, t.Line, t.StartCol, t.EndCol,
	)
}
