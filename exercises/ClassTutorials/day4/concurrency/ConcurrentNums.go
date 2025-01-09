package main

import (
	"fmt"
)

func square(num int) {
	defer wg.Done()
	r := num * num
	fmt.Printf("square of %d is %d\n", num, r)
}

func main2() {
	//var wg sync.WaitGroup
	/*slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, num := range slice {
		wg.Add(1)
		go square(num)
	}
	wg.Wait()
	fmt.Println("All squares calculated")*/
	ch := make(chan int, 5)
	ch2 := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	ch2 <- 4
	ch2 <- 5
	fmt.Println("c: ", <-ch)
	fmt.Println("c: ", <-ch)
	fmt.Println("c2: ", <-ch2)
	fmt.Println("c2: ", <-ch2)

}
