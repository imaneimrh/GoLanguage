package main

import (
	"log"

	"simplemath.com/utils"
	"um6p.ma/hello/mathutils"
)

func main() {
	log.Println("Hello Imane (;")
	a, b := 2000, 2003
	log.Println("Addition of ", a, " and ", b, " is ", mathutils.Add(a, b))
	log.Println("Difference of ", a, " and ", b, " is ", mathutils.Substract(a, b))
	log.Println("Square of ", a, " is ", utils.Square(a))

}
