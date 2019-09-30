package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/aqa/evaluator"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/object"
	"github.com/ollybritton/aqa/parser"
	"github.com/ollybritton/aqa/token"
	"github.com/spf13/cobra"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "repl [lex, parse]",
	Short: "repl creates a new REPL for running aqa++ code.",
	Long: `repl creates a new REPL for running aqa++ code.
	
	repl: ordinary repl
	repl lex: perform lexical analysis on the input text
	repl parse: parse the input text into an AST`,
	Run: func(cmd *cobra.Command, args []string) {
		l, err := readline.NewEx(&readline.Config{
			Prompt:            "\033[31m»\033[0m ",
			HistoryFile:       "/tmp/aqa-repl.hist.tmp",
			InterruptPrompt:   "^C",
			EOFPrompt:         "exit",
			HistorySearchFold: true,
		})
		if err != nil {
			fmt.Println(
				"error:", au.Red(err),
			)
		}
		defer l.Close()

		env := object.NewEnvironment()

		var lvl int
		var buffer string

		for {
			line, err := l.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}

			buffer = buffer + line
			end := eval(buffer, env)
			if end {
				buffer = ""
				lvl = 0
			} else {
				lvl++
				l.WriteStdin([]byte(
					strings.Repeat("\t", lvl),
				))
				buffer += "\n"
			}
		}
	},
}

func eval(input string, env *object.Environment) (end bool) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.Parse()
	var hadError bool
	for _, err := range p.Errors() {
		e, ok := err.(parser.InvalidTokenError)
		if !ok {
			hadError = true
			continue
		}

		if e.Unexpected.Type == token.EOF {
			return false
		}
	}

	if hadError {
		checkErrors(p)
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Println(
			au.Green(evaluated.Inspect()),
		)
	}

	fmt.Println("")
	return true
}

func init() {
	rootCmd.AddCommand(replCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// replCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// replCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
