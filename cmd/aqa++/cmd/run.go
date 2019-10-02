package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/ollybritton/aqa/object"

	au "github.com/logrusorgru/aurora"
	"github.com/ollybritton/aqa/evaluator"
	"github.com/ollybritton/aqa/lexer"
	"github.com/ollybritton/aqa/parser"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [filename]",
	Args:  cobra.ExactArgs(1),
	Short: "run runs a .aqa file and displays the output",
	Long: `run will run a file containing AQA++ source code.
For now, it will also print the result of the evaluation.`,
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println(au.Bold(au.Red("Could not read file:")))
			fmt.Println(au.Red(err))
		}

		str := string(bytes)

		l := lexer.New(str)
		p := parser.New(l)

		program := p.Parse()
		if checkErrors(p) {
			return
		}

		evaluator.Eval(program, object.NewEnvironment())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
