package main

import (
	"fmt"
)

func generatePrimes(n int, ch chan<- int) {
	for i := 2; i < n; i++ {
		ch <- i
	}
	close(ch)
}

func PrintPrimes(ch <-chan int) {
	for prime := range ch {
		fmt.Println(prime)
	}
}

func main3() {
	ch := make(chan int)

	ch <- 2

	/*
		select {
		case:
			time.Sleep(2*time.Second)
			prime := <-ch
			fmt.Println(prime)
		case <-time.After(2*time.Second):
			fmt.Println("Timeout")
		}*/

	go generatePrimes(10, ch)
	go PrintPrimes(ch)
}
