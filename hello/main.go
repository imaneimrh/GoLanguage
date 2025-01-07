package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/imusmanmalik/randomizer"
	mimo "simplemath.com/utils"
	"um6p.ma/hello/mathutils"
)

func main() {
	defer fmt.Println("Goodbye Imane (;")
	log.Println("Hello Imane (;")
	a, b := 2000, 2003
	log.Println("Addition of ", a, " and ", b, " is ", mathutils.Add(a, b))
	log.Println("Difference of ", a, " and ", b, " is ", mathutils.Substract(a, b))
	log.Println("Square of ", a, " is ", mimo.Square(a))
	res, err := randomizer.RandomInt(1, 100)
	log.Println("Random number: ", res, " Error: ", err)
	//Error is a pointer in the background
	dResult, err := mathutils.Divide(10, 0)
	if err != nil {
		log.Println("Error: ", err)
	} else {
		log.Println("Division result: ", dResult)
	}
	message, err_wrapping := FormatDevisionManager(10, 0)
	fmt.Println(message, err_wrapping)

}

func PerformDivision(a, b int) (float64, error) {
	return mathutils.Divide(a, b)
}
func FormatDevisionManager(a, b int) (string, error) {
	res, err := PerformDivision(a, b)
	if err != nil {
		return " ", errors.New(err.Error())
	}
	message := fmt.Sprintf("Division of %d and %d is %.2f", a, b, res)
	return message, nil

}
