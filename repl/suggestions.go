package repl

import "github.com/c-bata/go-prompt"

var suggestions = []prompt.Suggest{
	{Text: "<-", Description: "Assignment: a <- 10"},
	{Text: "==", Description: "Equal: a == b"},
	{Text: "!=", Description: "Not Equal: a != b"},

	{Text: "SUBROUTINE", Description: "Define a new subroutine."},
	{Text: "ENDSUBROUTINE", Description: "End a subroutine."},
	{Text: "CONSTANT", Description: "Define a constant value."},

	{Text: "IF", Description: "Start of an if statement."},
	{Text: "THEN", Description: "Goes after the condition in an if statement."},
	{Text: "ELSE", Description: "Start of an else block."},
	{Text: "ENDIF", Description: "End an if statement."},

	{Text: "true", Description: ""},
	{Text: "false", Description: ""},

	{Text: "%lex", Description: "Put the REPL into lex mode."},
	{Text: "%parse", Description: "Put the REPL into parse mode."},
	{Text: "%eval", Description: "Put the REPL into eval mode."},

	{Text: "exit", Description: "Exit the REPL."},
}
