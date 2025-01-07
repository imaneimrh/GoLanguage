package main

import (
	"fmt"
)

func main() {

	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	numbers := []int{1, 2, 3, 4, 5}
	for i, number := range numbers {
		fmt.Println(i, number)
	}

}
