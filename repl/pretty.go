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
	var literal = t.Literal

	if t.Type == token.ILLEGAL {
		ttype = fmt.Sprint(au.Red(t.Type))
	} else {
		ttype = fmt.Sprint(au.Green(au.Italic(t.Type)))
	}

	if literal == "\n" {
		literal = "\\n"
	}

	ttliteral := au.Bold(au.BrightYellow(literal))

	return fmt.Sprintf("(Lit: '%s', Type: '%s', line=%d, startcol=%d, endcol=%d)",
		ttliteral, ttype, t.Line, t.StartCol, t.EndCol,
	)
}

// Info prints some information that is displayed when a new prompt is started.
func Info() {
	fmt.Println(
		au.Sprintf(au.Cyan("aqa++ %v %v"), au.White("v1.0.0").Bold(), au.White("repl")),
	)

	fmt.Println(
		au.Sprintf(
			au.White("Type '%v', '%v', '%v' or '%v'"),
			au.BrightWhite("%lex").Bold(),
			au.BrightWhite("%parse").Bold(),
			au.BrightWhite("%eval").Bold(),
			au.BrightWhite("%help").Bold(),
		),
	)

	fmt.Println("")
}

// Help prints the help text.
func Help() {
	fmt.Println("")
	fmt.Println("Use the following commands to change the mode:")

	fmt.Println(
		au.Sprintf(au.BrightWhite("%q -- %v"), au.White("%lex").Italic(), au.Green("Lex the input (split into tokens)")),
	)
	fmt.Println(
		au.Sprintf(au.BrightWhite("%q -- %v"), au.White("%parse").Italic(), au.Green("Parse the input (form an AST)")),
	)
	fmt.Println(
		au.Sprintf(au.BrightWhite("%q -- %v"), au.White("%eval").Italic(), au.Green("Evaluate the input (run command)")),
	)

	fmt.Println("")
	fmt.Println("The following useful commands are also avaliable")
	fmt.Println(
		au.Sprintf(au.BrightWhite("%q -- %v"), au.White("%buf").Italic(), au.Green("Open up a buffer (in vim) to enter a multiline string")),
	)

	fmt.Println("")
	fmt.Println("To exit")
	fmt.Println(
		au.Sprintf(au.BrightWhite("%q -- %v"), au.White("quit").Italic(), au.Green("Quit the repl")),
	)

	fmt.Println("")
}
