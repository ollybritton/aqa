package parser

import (
	"fmt"

	"github.com/ollybritton/aqa/token"
)

// UnexpectedTokenError represents an error that occurs when the parser expects the next token to be something, but it isn't.
type UnexpectedTokenError struct {
	Message string

	CurTok      token.Token
	PeekTok     token.Token
	ExpectedTok token.Type
}

func (e UnexpectedTokenError) Error() string {
	return e.Message
}

// NewUnexpectedTokenError returns a new UnexpectedTokenError.
func NewUnexpectedTokenError(curTok, peekTok token.Token, expected token.Type) UnexpectedTokenError {
	msg := fmt.Sprintf("expected next token to be '%s', got '%s' instead. (line=%d, col=%d)", expected, peekTok.Type, peekTok.Line, peekTok.Column)

	return UnexpectedTokenError{
		Message: msg,

		CurTok:      curTok,
		PeekTok:     peekTok,
		ExpectedTok: expected,
	}
}

// InvalidTokenError occurs when there is a token where it shouldn't be. It is like UnexpectedTokenError, but it doesn't
// need an Expected field.
type InvalidTokenError struct {
	Message string

	CurTok     token.Token
	PeekTok    token.Token
	Unexpected token.Token
}

func (e InvalidTokenError) Error() string {
	return e.Message
}

// NewInvalidTokenError returns a new InvalidTokenError.
func NewInvalidTokenError(curTok, peekTok token.Token, unexpected token.Token) InvalidTokenError {
	msg := fmt.Sprintf("unexpected token '%s', invalid in context (line=%d, col=%d)", unexpected, unexpected.Line, unexpected.Column)

	return InvalidTokenError{
		Message: msg,

		CurTok:     curTok,
		PeekTok:    peekTok,
		Unexpected: unexpected,
	}
}

// IntegerParseError represents an error that occurs when trying to parse an string into a 64-bit integer.
type IntegerParseError struct {
	Message string

	CurTok  token.Token
	PeekTok token.Token
	Value   string
}

func (e IntegerParseError) Error() string {
	return e.Message
}

// NewIntegerParseError returns a new IntegerParseError.
func NewIntegerParseError(curTok, peekTok token.Token, value string) IntegerParseError {
	msg := fmt.Sprintf("could not parse %q as integer", value)

	return IntegerParseError{
		Message: msg,

		CurTok:  curTok,
		PeekTok: peekTok,
		Value:   value,
	}
}

// NoPrefixParseFnError represents an error that occurs when the parser encounters a token it doesn't have a prefix parse
// function for.
type NoPrefixParseFnError struct {
	Message string

	CurTok      token.Token
	PeekTok     token.Token
	UnknownType token.Type
}

func (e NoPrefixParseFnError) Error() string {
	return e.Message
}

// NewNoPrefixParseFnError returns a new NoPrefixParseFnError
func NewNoPrefixParseFnError(curTok, peekTok token.Token, unknown token.Type) NoPrefixParseFnError {
	msg := fmt.Sprintf("no prefix parse function for '%s' found (line=%d, col=%d)", unknown, curTok.Line, curTok.Column)

	return NoPrefixParseFnError{
		Message: msg,

		CurTok:      curTok,
		PeekTok:     peekTok,
		UnknownType: unknown,
	}
}
