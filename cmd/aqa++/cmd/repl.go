package cmd

import (
	"fmt"

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
		fmt.Println("repl called")
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
}
