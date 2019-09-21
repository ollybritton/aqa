package lexer

import (
	"github.com/ollybritton/aqa++/token"
)

// Lexer is a lexer for an AQA++ program.
// Its job is to translate a series of characters into chunks such as INTEGER(5) or TRUE. It also attaches information such
// as the position inside the input.
//
// BUG(me): At the moment, the lexer does not support Unicode, so only ASCII characters are supported.
type Lexer struct {
	input string

	position     int // Index to the current char the lexer is using.
	readPosition int // Index to the next char to be read.

	curLinePosition int // The position that token is located at, relative to the current line
	curLine         int // The line the current token is located on.

	startPosition int // The start position of the current token.

	ch byte // Current char under examination.

}

// New returns a new, initialised lexer.
// l.curLinePosition is set to -1 so that when one character is read, it will be set to 0.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.curLinePosition = -1
	l.readChar()

	return l
}

// readChar reads the next character in the input. If there are no characters left to read (i.e the input is finished or the
// input is blank), then the l.ch value is set to the NUL character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		l.curLinePosition++
	}

	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns the next char in the input as a byte.
// Like readChar, it returns the NUL character if there is no more input.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// skipWhitespace will skip over whitespace. If it encounters a newline, it increments the startLine and resets the startPosition.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// readIdentifier will reads a set of characters (including an underscore) and returns the string representation of that
// set of characters.
func (l *Lexer) readIdentifier() string {
	l.startPosition = l.curLinePosition
	start := l.position

	for isValidIdentCharacter(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

// readNumber reads a set of digits and returns a number as a string.
func (l *Lexer) readNumber() string {
	l.startPosition = l.curLinePosition
	start := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

// newSingleToken returns a new token from a token type.
func (l *Lexer) newSingleToken(tokenType token.Type) token.Token {
	return token.NewToken(tokenType, string(l.ch), l.curLine, l.curLinePosition)
}

// NextToken returns the next token in the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// Single characters
	case '+':
		tok = l.newSingleToken(token.PLUS)
	case '-':
		tok = l.newSingleToken(token.MINUS)
	case '*':
		tok = l.newSingleToken(token.ASTERISK)
	case '/':
		tok = l.newSingleToken(token.SLASH)
	case '>':
		tok = l.newSingleToken(token.GT)
	case ',':
		tok = l.newSingleToken(token.COMMA)
	case '(':
		tok = l.newSingleToken(token.LPAREN)
	case ')':
		tok = l.newSingleToken(token.RPAREN)
	case '[':
		tok = l.newSingleToken(token.LBRACKET)
	case ']':
		tok = l.newSingleToken(token.RBRACKET)
	case '{':
		tok = l.newSingleToken(token.LBRACE)
	case '}':
		tok = l.newSingleToken(token.RBRACE)

	// Could be single or double
	case '=': // = or ==
		if l.peekChar() == '=' {
			prev := l.ch
			l.readChar()

			tok = token.Token{
				Type:    token.EQ,
				Literal: string(prev) + string(l.ch),
				Line:    l.curLine,
				Column:  l.curLinePosition - 1, // Token started on the previous character
			}
		} else {
			// Assignment is handled using <- in AQA++
			tok = l.newSingleToken(token.ILLEGAL)
		}

	case '!': // ! or !=
		if l.peekChar() == '=' {
			prev := l.ch
			l.readChar()

			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(prev) + string(l.ch),
				Line:    l.curLine,
				Column:  l.curLinePosition - 1,
			}
		} else {
			tok = l.newSingleToken(token.BANG)
		}

	case '<': // < or <-
		if l.peekChar() == '-' {
			prev := l.ch
			l.readChar()

			tok = token.Token{
				Type:    token.ASSIGN,
				Literal: string(prev) + string(l.ch),
				Line:    l.curLine,
				Column:  l.curLinePosition - 1,
			}
		} else {
			tok = l.newSingleToken(token.LT)
		}

	// Newline handling
	case '\n':
		tok = l.newSingleToken(token.NEWLINE)
		l.readChar()

		l.curLine++
		l.curLinePosition = 0

		return tok

	// EOF
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Column = l.curLinePosition
		tok.Line = l.curLine

	// Multiple character handling
	default:
		if isValidIdentCharacter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.curLine
			tok.Column = l.startPosition
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			tok.Line = l.curLine
			tok.Column = l.startPosition
			return tok
		}

		tok = l.newSingleToken(token.ILLEGAL)
	}

	l.readChar()
	return tok
}
