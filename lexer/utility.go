package lexer

// isLetter returns true if the given character (byte) is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isDigit returns true if the character is a number.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isWhitespace returns true if the character is a type of whitespace (a space, a tab or a linefeed)
// Newlines are handled by the lexer as they are used for automatic semicolon insertion.
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

// isValidIdentCharacter returns true if the character is valid inside an identifier (a character or an underscore)
func isValidIdentCharacter(ch byte) bool {
	return isLetter(ch) || ch == '_'
}
