package main

import "fmt"

func init() {
	fmt.Println("Hello, Imane Init!")
}

const School string = "UM6P"

func main() {
	//fmt.Println("Hello, Imane!")
	//var x int
	//var y int64
	//var z float32
	//var isStudent bool
	//var name string
	//var x, y int
	//x = 10
	//y = 20
	//fmt.Printf("x = %d y = %d", x, y)
	//var age = 21
	//var name = "Imane"
	//fmt.Printf("age = %d name = %s", age, name)
	//fmt.Printf("\nage = %T name = %T", age, name)
	x, y, z := 15, 19, 18
	fmt.Printf("\nx, y, z = %d, %d, %d", x, y, z)
	const Pi float32 = 3.14
	fmt.Printf("\nPi = %f School = %s", Pi, School)

}
