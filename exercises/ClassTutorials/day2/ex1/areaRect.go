package main

import "fmt"

func main() {
	var l, w float64
	fmt.Println("Enter length of rectangle:")
	fmt.Scanln(&l)
	fmt.Println("Enter width of rectangle:")
	fmt.Scanln(&w)
	//--------------------Function Variable
	areaRect := func(length, width float64) (float64, bool) {
		if length < 0 || width < 0 {
			return -1, false
		} else {
			return length * width, true
		}
	}
	//--------------------
	a, _ := areaRect(l, w)
	if a > 0 {
		fmt.Println("Area of rectangle with length ", l, " and width ", w, " is ", a)
	} else {
		fmt.Println("Invalid input")
	}
	//------------------------- Closure
	inc := incrementer() //inc has the variable i and the function of the closure
	inc()
	inc()

}

func incrementer() func() int {
	i := 0
	fmt.Println("\nIncrementer called, i is ", i)
	return func() int {
		fmt.Println("Closure called, before increment, i is ", i)
		i++
		fmt.Println("Closure called, after... i is ", i)
		return i
	}
}
