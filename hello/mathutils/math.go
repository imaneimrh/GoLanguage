package mathutils

import "errors"

func Add(a int, b int) int {
	return a + b
}

func Substract(a int, b int) int {
	return a - b
}

func Divide(a int, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("Cannot divide by zero")
	}
	return float64(a) / float64(b), nil
}
