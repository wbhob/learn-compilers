package main

import "fmt"

// ASTNode represents a node in the Abstract Syntax Tree
// It's a flexible structure that can represent any element in our language
type ASTNode struct {
	Type     ASTNodeType            // Indicates the type of node (e.g., sequence, loop, value)
	Fields   map[string]interface{} // Stores additional data specific to the node type
	Children []*ASTNode             // Contains child nodes, if any
}

// ASTNodeType is a string that represents the type of an ASTNode
type ASTNodeType string

// These constants define the possible types of ASTNodes
const (
	SEQUENCE ASTNodeType = "sequence" // Represents a sequence of elements
	LOOP     ASTNodeType = "loop"     // Represents a loop (repeated elements)
	VALUE    ASTNodeType = "value"    // Represents a single numeric value
)

// parseSequence parses a sequence of elements
// <sequence> ::= <element> | <element> "," <sequence>
// This function is the entry point for parsing our language
func parseSequence(tokens []Token) *ASTNode {
	// Initialize a new ASTNode for the sequence
	sequence := &ASTNode{
		Type:     SEQUENCE,
		Children: make([]*ASTNode, 0),
	}

	// If there are no tokens, return an empty sequence
	if len(tokens) == 0 || tokens[0].Type == EOF {
		return sequence
	}

	// The sequence must start with a NUMBER or a LPAREN
	if tokens[0].Type != NUMBER && tokens[0].Type != LPAREN {
		panic(fmt.Sprintf("Unexpected token '%s' at the beginning of the sequence", tokens[0].Type))
	}

	// Initialize variables to keep track of the current element and parenthesis nesting
	element := make([]Token, 0)
	parenthesisCount := 0

	// Iterate through all tokens
	var nextToken Token
	for i, token := range tokens {
		// Look ahead to the next token, or use EOF if we're at the end
		if i+1 < len(tokens) {
			nextToken = tokens[i+1]
		} else {
			nextToken = Token{Type: EOF}
		}

		// Handle different token types
		switch token.Type {
		case LPAREN:
			parenthesisCount++
			// LPAREN can only be followed by a NUMBER or another LPAREN
			if nextToken.Type != NUMBER && nextToken.Type != LPAREN {
				panic(fmt.Sprintf("Unexpected token '%s' after left parenthesis", nextToken.Type))
			}
		case RPAREN:
			parenthesisCount--
			// RPAREN can only be followed by a comma, another RPAREN, EOF, or X
			if nextToken.Type != COMMA && nextToken.Type != RPAREN && nextToken.Type != EOF && nextToken.Type != X {
				panic(fmt.Sprintf("Unexpected token '%s' after right parenthesis", nextToken.Type))
			}
		case COMMA:
			// COMMA can only be followed by a NUMBER or a LPAREN
			if nextToken.Type != NUMBER && nextToken.Type != LPAREN {
				panic(fmt.Sprintf("Unexpected token '%s' after comma", nextToken.Type))
			}
		case NUMBER:
			// NUMBER can only be followed by a comma, RPAREN, EOF, or X
			if nextToken.Type != COMMA && nextToken.Type != RPAREN && nextToken.Type != EOF && nextToken.Type != X {
				panic(fmt.Sprintf("Unexpected token '%s' after number", nextToken.Type))
			}
		case X:
			// X can only be followed by a NUMBER
			if nextToken.Type != NUMBER {
				panic(fmt.Sprintf("Unexpected token '%s' after 'x'", nextToken.Type))
			}
		case EOF:
			// EOF can only be followed by EOF
			if nextToken.Type != EOF {
				panic(fmt.Sprintf("Unexpected token '%s' after EOF", nextToken.Type))
			}
		}

		// Process element when a comma is found outside of parentheses
		if token.Type == COMMA && parenthesisCount == 0 {
			// Parse the accumulated element and add it to the sequence
			sequence.Children = append(sequence.Children, parseElement(element))
			// Reset the element for the next iteration
			element = make([]Token, 0)
		} else {
			// Add the current token to the element
			element = append(element, token)
		}

		// Ensure parentheses are balanced
		if parenthesisCount < 0 {
			panic("Invalid sequence: unmatched parentheses")
		}
	}

	// Process the final element
	if len(element) > 0 {
		sequence.Children = append(sequence.Children, parseElement(element))
	}

	return sequence
}

// parseElement parses a single element
// <element> ::= <number> | <group> | <loop>
func parseElement(tokens []Token) *ASTNode {
	// Check if the element is a loop by looking for the 'x' token
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Type == X {
			return parseLoop(tokens)
		}
	}

	// Check if the element is a number (starts with a NUMBER token)
	if len(tokens) > 0 && tokens[0].Type == NUMBER {
		return parseNumber(tokens)
	}

	// If not a number or loop, parse as a group
	return parseGroup(tokens)
}

// parseLoop parses a loop element
// <loop> ::= <element> "x" <integer>
func parseLoop(tokens []Token) *ASTNode {
	// Initialize slices to hold tokens before and after the 'x'
	left, right := make([]Token, 0), make([]Token, 0)
	isRight := false

	// Keep track of parenthesis nesting
	parenthesisCount := 0
	for _, token := range tokens {
		if token.Type == LPAREN {
			parenthesisCount++
		} else if token.Type == RPAREN {
			parenthesisCount--
		}

		// When we find the 'x' token outside of parentheses, start filling the right slice
		if token.Type == X && parenthesisCount == 0 {
			isRight = true
		} else if isRight {
			right = append(right, token)
		} else {
			left = append(left, token)
		}
	}

	// Parse the element to be repeated
	var element *ASTNode
	if left[0].Type == NUMBER {
		element = parseNumber(left)
	} else {
		element = parseGroup(left)
	}

	// Create and return the loop node
	return &ASTNode{
		Type: LOOP,
		Fields: map[string]interface{}{
			"repeat": element,
			"count":  parseNumber(right),
		},
	}
}

// parseGroup parses a group element
// <group> ::= "(" <sequence> ")"
func parseGroup(tokens []Token) *ASTNode {
	// Extract the sequence within the parentheses
	sequence := make([]Token, 0)
	parenthesisCount := 0
	for _, token := range tokens {
		if token.Type == LPAREN {
			parenthesisCount++
			if parenthesisCount == 1 {
				continue // Skip the outermost opening parenthesis
			}
		} else if token.Type == RPAREN {
			parenthesisCount--
			if parenthesisCount == 0 {
				break // Stop at the outermost closing parenthesis
			}
		}

		sequence = append(sequence, token)
	}

	// Parse the extracted sequence
	return parseSequence(sequence)
}

// parseNumber parses a number element
// <value> ::= <integer> | <decimal>
func parseNumber(tokens []Token) *ASTNode {
	// Create a VALUE node with the numeric value
	return &ASTNode{
		Type: VALUE,
		Fields: map[string]interface{}{
			"value": tokens[0].Value,
		},
	}
}
