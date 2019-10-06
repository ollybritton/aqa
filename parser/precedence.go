package parser

import "github.com/ollybritton/aqa/token"

// Definitions of operator precedences.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	SHIFT       // >> or <<
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // fn(x)
	INDEX       // array[index
)

// Mappings of precedences to their token types.
var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LSHIFT:   SHIFT,
	token.RSHIFT:   SHIFT,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}
