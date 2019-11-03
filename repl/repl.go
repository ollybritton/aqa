package repl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/aqa/evaluator"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/object"
	"github.com/ollybritton/aqa/parser"
	"github.com/ollybritton/aqa/token"
)

// Repl represents a repl. It can be used to lex, parse or evaluate AQA++ code.
type Repl struct {
	Buffer bytes.Buffer
	Prompt *prompt.Prompt

	Env *object.Environment

	Mode  string // Either "lex", "parse" or "eval"
	Level int

	input string
}

// New returns a new, initialised REPL.
func New() *Repl {
	r := &Repl{Mode: "eval", Env: object.NewEnvironment()}
	r.Prompt = prompt.New(
		r.Execute,
		r.Completor,

		prompt.OptionLivePrefix(r.Prefix),
		prompt.OptionTitle("aqa++"),

		prompt.OptionSuggestionBGColor(prompt.Red),
		prompt.OptionSuggestionTextColor(prompt.Black),

		prompt.OptionSelectedSuggestionBGColor(prompt.Red),
		prompt.OptionSelectedDescriptionBGColor(prompt.Turquoise),
		prompt.OptionSelectedDescriptionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),

		prompt.OptionInputTextColor(prompt.Turquoise),

		prompt.OptionAddKeyBind(prompt.KeyBind{
			Fn: func(buf *prompt.Buffer) {
				if r.Level > 0 {
					r.Level--
				}
			},
			Key: prompt.Backspace,
		}),
	)

	return r
}

// Execute is what executes a command inside the REPL.
func (r *Repl) Execute(input string) {
	r.Level = 0

	if strings.HasPrefix(input, "%") {
		switch input[1:len(input)] {
		case "lex":
			r.Mode = "lex"
			fmt.Println(au.Green("Mode set to 'lex'."))
			fmt.Println("")

			return

		case "parse":
			r.Mode = "parse"
			fmt.Println(au.Green("Mode set to 'parse'."))
			fmt.Println("")

			return

		case "eval", "exec":
			r.Mode = "eval"
			fmt.Println(au.Green("Mode set to 'eval'."))
			fmt.Println("")

			return

		case "buf":
			err := ioutil.WriteFile("/tmp/.aqa-buf.aqa", []byte{}, 0777)
			if err != nil {
				fmt.Println(au.Red("Error clearing buffer:").Bold())
				fmt.Println(au.Red(err))

				return
			}

			cmd := exec.Command("vim", "/tmp/.aqa-buf.aqa")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			err = cmd.Run()

			if err != nil {
				fmt.Println(au.Red("Error opening vim buffer:").Bold())
				fmt.Println(au.Red(err))

				return
			}

			bytes, err := ioutil.ReadFile("/tmp/.aqa-buf.aqa")
			if err != nil {
				fmt.Println(au.Red("Error opening vim buffer:").Bold())
				fmt.Println(au.Red(err))

				return
			}

			input = string(bytes)

			fmt.Println("")
			fmt.Println(au.Green("Added from buffer:").Bold())
			fmt.Println(au.Yellow(input))

		default:
			message := au.Red(au.Bold(
				fmt.Sprintf("No magic command %q found.", input),
			))

			fmt.Println(message)
			fmt.Println("")

			return
		}
	}

	switch r.Mode {
	case "lex":
		r.Lex(input)
	case "parse":
		r.Parse(input)
	case "eval":
		r.Eval(input)
	}
}

// Completor is what completes input inside the REPL.
func (r *Repl) Completor(input prompt.Document) []prompt.Suggest {
	w := input.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(suggestions, w, true)
}

// Prefix is what calculates the prefix/identation level.
func (r *Repl) Prefix() (string, bool) {
	if r.Level > 0 {
		indent := strings.Repeat("  ", r.Level)
		return "(" + r.Mode + ") > " + indent, true
	}

	return "(" + r.Mode + ") > ", true
}

// Lex lexes a given input string, and displays the results to stdout.
func (r *Repl) Lex(input string) {
	l := lexer.New(input)

	tokens := []token.Token{}
	tok := l.NextToken()

	for tok.Type != token.EOF {
		tokens = append(tokens, tok)
		tok = l.NextToken()
	}

	fmt.Println("")

	for i, t := range tokens {
		num := au.Blue(fmt.Sprintf("[%d]", i))
		fmt.Printf("%v %v\n", num, PrettyToken(t))
	}

	fmt.Println("")
}

// Parse parses a given input string, and displays the results to stdout.
func (r *Repl) Parse(input string) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.Parse()
	if len(p.Errors()) != 0 {
		Errors(p.Errors())
	}

	fmt.Println(program)
	fmt.Println("")
}

// Eval evaluates a given input string, and displays the results to stdout.
func (r *Repl) Eval(input string) {
	obj, errors := evaluator.EvalString(input, r.Env)

	if len(errors) != 0 {
		Errors(errors)
	}

	if obj == nil || obj.Type() == object.NULL_OBJ {
		return
	}

	if obj.Type() == object.ERROR_OBJ {
		fmt.Println(au.Red(obj.Inspect()).Bold())
		return
	}

	fmt.Println(au.Green(obj.Inspect()))
	fmt.Println("")
}

// Start starts the REPL.
func (r *Repl) Start() {
	for {

		input := r.Prompt.Input()

		// l := lexer.New(input)
		// p := parser.New(l)

		// p.Parse()
		// if len(p.Errors()) != 0 {
		// 	isEOF := false

		// 	for _, e := range p.Errors() {
		// 		switch e.(type) {
		// 		case parser.NoPrefixParseFnError:
		// 			isEOF = true
		// 			break
		// 		}
		// 	}

		// 	if isEOF {
		// 		r.Level++
		// 		continue
		// 	}
		// }

		switch {
		case input == "exit":
			os.Exit(0)
		case input == "ping":
			fmt.Println("pong")
			fmt.Println("")
		default:
			r.Execute(input)
		}

	}
}
