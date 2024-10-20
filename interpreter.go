package main

import (
	"fmt"
	"strconv"
)

func runSequence(ast *ASTNode) []float64 {
	// ast root is always a sequence
	results := make([]float64, 0)
	for _, child := range ast.Children {
		switch child.Type {
		case LOOP:
			results = append(results, runLoop(child)...)
		case VALUE:
			results = append(results, runNumber(child))
		case SEQUENCE:
			results = append(results, runSequence(child)...)
		}
	}

	return results
}

func runLoop(ast *ASTNode) []float64 {
	element := ast.Fields["repeat"].(*ASTNode)
	count := ast.Fields["count"].(*ASTNode)

	if count.Type != VALUE {
		panic(fmt.Sprintf("Count is not a number: %+v", count))
	}

	countNum, err := strconv.Atoi(count.Fields["value"].(string))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse count: %s", count.Fields["value"].(string)))
	}

	results := make([]float64, 0)
	for i := 0; i < countNum; i++ {
		switch element.Type {
		case VALUE:
			results = append(results, runNumber(element))
		case SEQUENCE:
			results = append(results, runSequence(element)...)
		default:
			panic(fmt.Sprintf("Unknown element type: %s", element.Type))
		}
	}

	return results
}

func runNumber(ast *ASTNode) float64 {
	str := ast.Fields["value"].(string)
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse number: %s", str))
	}

	return num
}
