//Danny Radosevich
//Examples of extending the reader

package commands

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

func ScanArgs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start, first := 0, rune(0) //start is the index of the start of the token,
	// first is the first character of the token
	for width := 0; start < len(data); start += width {
		// Iterate over the input data, looking for the start of a token
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		// Decode the next rune (character) from the input data
		if !unicode.IsSpace(r) {
			// If the character is not a space, it is the start of a token
			first = r
			break
		}
	}
	if isQuote(first) {
		start++ //skip the opening quote
	}
	for width, i := 0, start; i < len(data); i += width {
		// Iterate over the input data, looking for the end of the token
		var r rune
		r, width = utf8.DecodeRune(data[i:]) // Decode the next rune (character) from the input data
		if ok := isQuote(first); !ok && unicode.IsSpace(r) || ok && r == first {
			// If we are not in a quoted token and we encounter a space,
			// or if we are in a quoted token and we encounter the closing quote,
			// then we have reached the end of the token
			return i + width, data[start:i], nil
			// Return the index of the next token, the current token, and no error
		}
	}
	if atEOF && len(data) > start {
		// If we have reached the end of the input data and there is still a token to return,
		// return it
		if isQuote(first) {
			err = fmt.Errorf("unterminated quote: %q", first)
		}
		return len(data), data[start:], err
	}
	if isQuote(first) {
		start-- // if we are at the end of the data and we have an unterminated quote,
		// we want to include the opening quote in the token
	}
	return start, nil, nil // Return the index of the next token, no token, and no error
}
