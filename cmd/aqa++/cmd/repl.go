package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/logrusorgru/aurora"
	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/aqa/evaluator"
	"github.com/ollybritton/aqa/object"
	"github.com/ollybritton/aqa/repl"
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
		shouldLex, err := cmd.Flags().GetBool("lex")
		if err != nil {
			fmt.Println(au.Red("Could not fetch 'lex' flag:").Bold())
			fmt.Println(au.Red(err))
			return
		}

		shouldParse, err := cmd.Flags().GetBool("parse")
		if err != nil {
			fmt.Println(au.Red("Could not fetch 'parse' flag:").Bold())
			fmt.Println(au.Red(err))
			return
		}

		if shouldLex && shouldParse {
			fmt.Println(au.Red("Both the --lex and --parse flags are passed, please pass only one.").Bold())
			return
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(au.Red("Could not fetch 'file' flag:").Bold())
			fmt.Println(au.Red(err))
			return
		}

		env := object.NewEnvironment()

		if file != "" {
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println(au.Red("Could not read file " + file).Bold())
				fmt.Println(au.Red(err))
				return
			}

			eval, errs := evaluator.EvalString(string(bytes), env)
			if len(errs) != 0 {
				for _, err := range errs {
					fmt.Println(au.Red(err.Error()))
				}
				return
			}

			if eval.Type() == object.ERROR_OBJ {
				fmt.Println(aurora.Red("Error running starting file:").Bold())
				fmt.Println(aurora.Red(eval.Inspect()))
			}

			fmt.Println("")
		}

		r := repl.New()
		r.Env = env

		if shouldLex {
			r.Mode = "lex"
		} else if shouldParse {
			r.Mode = "parse"
		}

		r.Start()
	},
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

	replCmd.Flags().BoolP("lex", "l", false, "lex the input")
	replCmd.Flags().BoolP("parse", "p", false, "parse the input")

	replCmd.Flags().StringP("file", "f", "", "eval/lex/parse this file and then start repl")

}
