package handler

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"strings"
)

// метод вычисления выражения
func SolveMathExpression(expression string) (float64, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	finalResult, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("unable to convert result to float64")
	}

	return finalResult, nil
}

// Функция для подсчета символа '+'
func CountPlus(input string) int {
	count := 0
	for _, char := range input {

		if char == '+' {
			count++
		}
	}
	return count
}

// Функция для подсчета символа '-'
func CountMinus(input string) int {
	count := 0
	for _, char := range input {
		if char == '-' {
			count++
		}
	}
	return count
}

// Функция для подсчета символа '*'
func CountMultiply(input string) int {
	count := 0
	for _, char := range input {
		if char == '*' {
			count++
		}
	}
	return count
}

// Функция для подсчета символа '/'
func CountDivide(input string) int {
	count := 0
	for _, char := range input {
		if char == '/' {
			count++
		}
	}
	return count
}
func hasPlus(input string) bool {
	return strings.Contains(input, "+")
}

func hasMinus(input string) bool {
	return strings.Contains(input, "-")
}

func hasMultiply(input string) bool {
	return strings.Contains(input, "*")
}

func hasDivide(input string) bool {
	return strings.Contains(input, "/")
}
