package repl

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/parser"
	"github.com/ollybritton/aqa/token"
)

// Repl represents a repl. It can be used to lex, parse or evaluate AQA++ code.
type Repl struct {
	Buffer bytes.Buffer
	Prompt *prompt.Prompt

	Mode  string // Either "lex", "parse" or "eval"
	Level int
}

// New returns a new, initialised REPL.
func New() *Repl {
	r := &Repl{Mode: "eval"}
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
			fmt.Println("")
			fmt.Println(au.Green("Mode set to 'lex'."))

		case "parse":
			r.Mode = "parse"
			fmt.Println("")
			fmt.Println(au.Green("Mode set to 'parse'."))

		case "eval", "exec":
			r.Mode = "eval"
			fmt.Println(au.Green("Mode set to 'eval'."))

		default:
			message := au.Red(au.Bold(
				fmt.Sprintf("No magic command %q found.", input),
			))

			fmt.Println(message)
			fmt.Println("")
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
func (r *Repl) Lex(input string) {}

// Parse parses a given input string, and displays the results to stdout.
func (r *Repl) Parse(input string) {}

// Eval evaluates a given input string, and displays the results to stdout.
func (r *Repl) Eval(input string) {}

// Start starts the REPL.
func (r *Repl) Start() {
	for {
		input := r.Prompt.Input()

		if input == "" {
			continue
		}

		if r.Buffer.String() != "" {
			input = r.Buffer.String() + input
		}

		if r.Mode == "lex" {
			r.Execute(input)
			return
		}

		l := lexer.New(input)
		p := parser.New(l)

		p.Parse()
		fatalError := true

		for _, err := range p.Errors() {
			invalid, ok := err.(parser.InvalidTokenError)
			if !ok {
				continue
			}

			if invalid.Unexpected.Type != token.EOF {
				continue
			}

			fatalError = false
		}

		if fatalError {
			Errors(p.Errors())
			r.Level = 0
			return
		}

		r.Level++
		r.Buffer.WriteString(input + "\n")

	}
}
