package main

import (
	"log"
)

func mainy() {
	numbers := [5]int{17, 16, 15, 18, 14}
	slice := []int{}
	for i := 0; i < len(numbers); i++ {
		log.Println(numbers[i])
	}
	slice = append(slice, 17)
	slice = append(slice, 16)
	slice = append(slice, 18)
	//slice = append(slice[:2], slice[3:]...)
	log.Println(slice)

	nullslice := []int{1, 2, 3}
	log.Println(nullslice)
	nullslice = append(nullslice, 1)

	nums := make([]int, 5)
	nums = append(nums, 1)
	log.Println(nums)
	loopNumbers(slice)

}

func loopNumbers(numbers []int) {
	for index, value := range numbers {
		log.Println(index, "->", value)
		log.Printf("%d - %d %d %p\n", index, len(numbers), cap(numbers), numbers)
	}

}
