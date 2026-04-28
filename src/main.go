package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	Number TokenType = iota
	Operator
	Function
	LeftParen
	RightParen
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) []Token {
	var tokens []Token
	i := 0
	expectNumber := true

	for i < len(input) {
		ch := input[i]

		if ch == ' ' {
			i++
			continue
		}

		if unicode.IsDigit(rune(ch)) || ch == '.' || (ch == '-' && expectNumber && i+1 < len(input) && (unicode.IsDigit(rune(input[i+1])) || input[i+1] == '.')) {
			start := i
			if ch == '-' {
				i++
			}
			for i < len(input) && (unicode.IsDigit(rune(input[i])) || input[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{Number, input[start:i]})
			expectNumber = false
			continue
		}

		if unicode.IsLetter(rune(ch)) {
			start := i
			for i < len(input) && unicode.IsLetter(rune(input[i])) {
				i++
			}
			tokens = append(tokens, Token{Function, input[start:i]})
			expectNumber = false
			continue
		}

		if strings.ContainsRune("+-*/^", rune(ch)) {
			tokens = append(tokens, Token{Operator, string(ch)})
			i++
			expectNumber = true
			continue
		}

		if ch == '(' {
			tokens = append(tokens, Token{LeftParen, "("})
			i++
			expectNumber = true
			continue
		}
		if ch == ')' {
			tokens = append(tokens, Token{RightParen, ")"})
			i++
			expectNumber = false
			continue
		}

		panic("unknown character: " + string(ch))
	}

	return tokens
}

var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

func toRPN(tokens []Token) []Token {
	var output []Token
	var stack []Token

	for _, tok := range tokens {
		switch tok.Type {

		case Number:
			output = append(output, tok)

		case Function:
			stack = append(stack, tok)

		case Operator:
			for len(stack) > 0 {
				top := stack[len(stack)-1]

				if top.Type == Function ||
					(top.Type == Operator && precedence[top.Value] >= precedence[tok.Value]) {

					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, tok)

		case LeftParen:
			stack = append(stack, tok)

		case RightParen:
			for len(stack) > 0 && stack[len(stack)-1].Type != LeftParen {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			stack = stack[:len(stack)-1]

			if len(stack) > 0 && stack[len(stack)-1].Type == Function {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}

func evalRPN(tokens []Token) float64 {
	var stack []float64

	for _, tok := range tokens {
		switch tok.Type {

		case Number:
			val, _ := strconv.ParseFloat(tok.Value, 64)
			stack = append(stack, val)

		case Operator:
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var res float64
			switch tok.Value {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				res = a / b
			case "^":
				res = math.Pow(a, b)
			}
			stack = append(stack, res)

		case Function:
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var res float64
			switch tok.Value {
			case "sin":
				res = math.Sin(a)
			case "cos":
				res = math.Cos(a)
			case "tan":
				res = math.Tan(a)
			case "log":
				res = math.Log(a)
			case "sqrt":
				res = math.Sqrt(a)
			case "abs":
				res = math.Abs(a)
			default:
				panic("unknown function: " + tok.Value)
			}
			stack = append(stack, res)
		}
	}

	return stack[0]
}

func evalExpression(input string) float64 {
	tokens := tokenize(input)
	rpn := toRPN(tokens)
	return evalRPN(rpn)
}

func main() {
	tests := []string{
		"2 + 3 * 4",
		"sqrt(16)",
		"sin(0)",
		"2^3 + 1",
		"abs(-5) + 3",
		"3 + 4 * 2 / (1 - 5)^2",
	}

	for _, t := range tests {
		fmt.Println(t, "=", evalExpression(t))
	}
}
