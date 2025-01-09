package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func square(num int) {
	defer wg.Done()
	r := num * num
	fmt.Printf("square of %d is %d\n", num, r)
}

func main() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, num := range slice {
		wg.Add(1)
		go square(num)
	}
	wg.Wait()
	fmt.Println("All squares calculated")
}
