package main

import (
	"bufio"
	"errors"
	. "fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type (
	RPNStack struct {
		rpnSlice []string
	}
	ResStack struct {
		resSlice []int
	}
	StackManager interface {
		Push(string)
		Pop() (string, error)
	}
)

func (s *RPNStack) Push(v string) {
	s.rpnSlice = append(s.rpnSlice, v)
}

func (s *RPNStack) Pop() (string, error) {
	if len(s.rpnSlice) == 0 {
		return "", io.EOF
	}
	v := s.rpnSlice[len(s.rpnSlice)-1]
	s.rpnSlice = s.rpnSlice[:len(s.rpnSlice)-1]
	return v, nil
}

func (rs *ResStack) Push(v int) {
	rs.resSlice = append(rs.resSlice, v)
}

func (rs *ResStack) Pop() (int, error) {
	if len(rs.resSlice) == 0 {
		return 0, io.EOF
	}
	v := rs.resSlice[len(rs.resSlice)-1]
	rs.resSlice = rs.resSlice[:len(rs.resSlice)-1]
	return v, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	mathVars := make(map[string]int)
	var mathStack RPNStack
	var finalStack ResStack

	for scanner.Scan() {
		input := scanner.Text()
		inputSlc := strings.Fields(input)

		switch {
		case input == "/exit":
			Println("Bye!")
			os.Exit(0)
		case input == "/help":
			Println("The program performs the calculator's functionality")
		case input == "":
			continue
		case strings.ContainsRune(input, '='):
			storeVars(mathVars, input)
		case len(inputSlc) >= 1 && input[0] != '/':
			if len(inputSlc) == 2 {
				Println("Invalid expression")
			} else {
				final, err := mathStack.turnIntoRPN(input)
				if err != nil {
					Println(err)
				} else {
					finalStack.getResult(final, mathVars)
				}
			}
		default:
			Println("Unknown command")
		}
	}
}

func storeVars(m map[string]int, input string) {
	equation := strings.Split(input, "=")

	for key, val := range m {
		if strings.TrimSpace(equation[1]) == key {
			m[strings.TrimSpace(equation[0])] = val
			return
		}
	}

	value, err := strconv.Atoi(strings.TrimSpace(equation[1]))
	if err != nil {
		Println("Invalid identifier")
		return
	}

	for i := range equation[0] {
		if unicode.IsDigit(rune(equation[0][i])) {
			Println("Invalid identifier")
			return
		}
	}

	m[strings.TrimSpace(equation[0])] = value
}

// TO-DO: handle exponents
func (s *RPNStack) turnIntoRPN(infix string) (string, error) {
	var b strings.Builder
	infSlice := strings.Fields(infix)
	for _, char := range infSlice {
		switch {
		case strings.Count(infix, "(") != strings.Count(infix, ")"):
			return "", errors.New("Invalid expression")

		case (unicode.IsDigit(rune(char[0])) || unicode.IsLetter(rune(char[0]))) && !strings.Contains(char, ")"):
			b.WriteString(char + " ")

		case char[0] == '+' || char[0] == '-':
			if strings.Count(char, "-")%2 == 0 {
				char = strings.ReplaceAll(char, char, "+")
			}
			if len(s.rpnSlice) == 0 ||
				s.rpnSlice[len(s.rpnSlice)-1] == "(" {
				s.Push(char)
			} else {
				popped, _ := s.Pop()
				for {
					if len(s.rpnSlice) == 0 {
						b.WriteString(popped + " ")
						s.Push(char)
						break
					} else if popped == "(" {
						s.Push(popped)
						s.Push(char)
						break
					} else {
						b.WriteString(popped + " ")
						popped, _ = s.Pop()
					}
				}
			}

		case char[0] == '*' || char[0] == '/':
			if len(char) > 1 {
				return "", errors.New("Invalid expression")
			}
			if len(s.rpnSlice) == 0 ||
				s.rpnSlice[len(s.rpnSlice)-1] == "+" ||
				s.rpnSlice[len(s.rpnSlice)-1] == "-" ||
				s.rpnSlice[len(s.rpnSlice)-1] == "(" {
				s.Push(char)
			} else {
				popped, _ := s.Pop()
				b.WriteString(popped + " ")
				s.Push(char)
			}

		case char[0] == '^':
			if len(char) > 1 {
				return "", errors.New("Invalid expression")
			}
			if s.rpnSlice[len(s.rpnSlice)-1] == "^" {
				popped, _ := s.Pop()
				b.WriteString(popped + " ")
				s.Push(char)
			} else {
				s.Push(char)
			}

		case strings.Contains(char, "("):
			var c strings.Builder
			for j := range char {
				if char[j] == '(' {
					s.Push(string(char[j]))
				} else {
					c.WriteRune(rune(char[j]))
				}
			}
			b.WriteString(c.String() + " ")

		case strings.Contains(char, ")"):
			var d strings.Builder
			for k := range char {
				if unicode.IsDigit(rune(char[k])) || unicode.IsLetter(rune(char[k])) {
					d.WriteRune(rune(char[k]))
				} else {
					s.Push(string(char[k]))
				}
			}
			b.WriteString(d.String() + " ")

			popped, _ := s.Pop()
			for {
				if popped[0] == '(' {
					break
				} else {
					if popped[0] == ')' {
						popped, _ = s.Pop()
					} else {
						b.WriteString(popped + " ")
						popped, _ = s.Pop()
					}
				}
			}
		}
	}
	for len(s.rpnSlice) > 0 {
		poppedFinal, _ := s.Pop()
		b.WriteString(poppedFinal + " ")
	}
	return b.String(), nil
}

func (rs *ResStack) getResult(postFix string, m map[string]int) {
	numSlc := strings.Fields(postFix)

	for _, element := range numSlc {
		if element == "" {
			continue
		}

		switch {
		case unicode.IsLetter(rune(element[0])):
			if value, ok := m[element]; ok {
				rs.Push(value)
			} else {
				Println("Unknown variable")
				return
			}
		case unicode.IsDigit(rune(element[0])):
			number, _ := strconv.Atoi(element)
			rs.Push(number)
		default:
			popOne, _ := rs.Pop()
			popTwo, _ := rs.Pop()
			switch element {
			case "+":
				rs.Push(popOne + popTwo)
			case "-":
				rs.Push(popTwo - popOne)
			case "*":
				rs.Push(popOne * popTwo)
			case "/":
				quotient := float64(popTwo) / float64(popOne)
				rs.Push(int(quotient))
			case "^":
				rs.Push(int(math.Pow(float64(popTwo), float64(popOne))))
			}
		}
	}

	Println(rs.resSlice[0])
	rs.resSlice = rs.resSlice[:0]
}
