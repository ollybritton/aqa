package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aqa++ <command> [arguments]",
	Short: "aqa++ is an implementation of the AQA pseudocode specification.",
	Long: `aqa++ is an implementation of the AQA pseudocode specification in Go.
  
  This command line application provides tools such as running .aqa files and providing a REPL. It also provides an interface
  to the internal lexer and parser system.

  Example:
    aqa++ run file.aqa
    aqa++ repl

    aqa++ repl lex
    aqa++ repl parse`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
