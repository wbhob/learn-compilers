package main

import "fmt"

func ParseNumberSequenceShorthand(input string) []float64 {
	fmt.Print("==========================================\n")
	fmt.Printf("Parsing number sequence shorthand: %s\n", input)
	tokens := lex(input)
	ast := parseSequence(tokens)
	return runSequence(ast)
}

func ValidateNumberSequenceShorthand(input string) error {
	tokens := lex(input)
	parseSequence(tokens)
	return nil
}

func main() {
	fmt.Println("Run the test cases with `make test`.")
}
