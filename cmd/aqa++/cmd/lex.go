package cmd

import (
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/token"

	au "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// lexCmd represents the lex command
var lexCmd = &cobra.Command{
	Use:   "lex",
	Short: "lex performs lexical analysis on the input specified",
	Long:  `lex performs lexical analysis on the input specified, and will output a set of tokens.`,
	Run: func(cmd *cobra.Command, args []string) {

		l, err := readline.NewEx(&readline.Config{
			Prompt:            "\033[31m»\033[0m ",
			HistoryFile:       "/tmp/aqa-lex.hist.tmp",
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

			lex(line)
			fmt.Println("")
		}

	},
}

func prettyToken(t token.Token) string {
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

func lex(input string) {
	l := lexer.New(input)

	tokens := []token.Token{}
	tok := l.NextToken()

	for tok.Type != token.EOF {
		tokens = append(tokens, tok)
		tok = l.NextToken()
	}

	for i, t := range tokens {
		num := au.Blue(fmt.Sprintf("[%d]", i))
		fmt.Printf("%v %v\n", num, prettyToken(t))
	}
}

func init() {
	replCmd.AddCommand(lexCmd)
}
