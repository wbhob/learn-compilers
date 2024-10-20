package main

import (
	"unicode"
)

// Token represents a lexical unit in our language
type Token struct {
	Type  TokenType // The type of the token
	Value string    // The string value of the token
}

// TokenType is a string that represents the category of a token
type TokenType string

// These constants define the possible types of tokens in our language
const (
	NUMBER TokenType = "NUMBER" // Represents numeric values
	COMMA  TokenType = "COMMA"  // Represents the comma separator
	X      TokenType = "X"      // Represents the 'x' used for repetition
	LPAREN TokenType = "LPAREN" // Represents a left parenthesis
	RPAREN TokenType = "RPAREN" // Represents a right parenthesis
	EOF    TokenType = "EOF"    // Represents the end of the input
)

// lex function takes a string input and returns a slice of Tokens
// It breaks down the input string into individual lexical units (tokens)
func lex(input string) []Token {
	tokens := []Token{}
	digits := ""

	// Iterate through each character in the input string
	for _, char := range input {
		// Skip whitespace characters
		if unicode.IsSpace(char) {
			// If we've been building a number, add it as a token
			if digits != "" {
				tokens = append(tokens, Token{Type: NUMBER, Value: digits})
				digits = ""
			}
			continue
		}

		// Handle characters that could be part of a number
		if char == '.' || char == '-' || unicode.IsDigit(char) {
			digits += string(char)
			continue
		}

		// If we've been building a number, add it as a token
		if digits != "" {
			tokens = append(tokens, Token{Type: NUMBER, Value: digits})
			digits = ""
		}

		// Handle special characters
		switch char {
		case ',':
			tokens = append(tokens, Token{Type: COMMA, Value: ","})
		case 'x':
			tokens = append(tokens, Token{Type: X, Value: "x"})
		case '(':
			tokens = append(tokens, Token{Type: LPAREN, Value: "("})
		case ')':
			tokens = append(tokens, Token{Type: RPAREN, Value: ")"})
		default:
			// If we encounter an unknown character, panic
			panic("Unknown character: " + string(char))
		}
	}

	// If we've been building a number, add it as a token
	if digits != "" {
		tokens = append(tokens, Token{Type: NUMBER, Value: digits})
	}

	// Add an EOF token to signify the end of the input
	tokens = append(tokens, Token{Type: EOF, Value: ""})

	return tokens
}
