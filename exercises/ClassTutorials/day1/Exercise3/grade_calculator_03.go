package main

import (
	"fmt"
)

func main() {
	var grade float64
	var average float64
	var sum float64
	var i int

	fmt.Printf("Welcome to the grade Calculator\nEnter -1 to exit the program\n")
	for grade != -1 {
		fmt.Printf("\nPlease, enter grade %d (from 0 to 100 / Enter -1 to exit the program): ", i+1)
		fmt.Scan(&grade)
		switch {
		case grade < 0 && grade != -1:
			fmt.Println("\nATTENTION: Grade can't be less than 0!!")
			continue
		case grade > 100:
			fmt.Println("\nATTENTION: Grade can't be more than 100!!")
			continue
		case grade == -1:
			fmt.Println("\nYou chose to exit the program....")
			continue
		default:
			sum += grade
			i++
		}

	}
	if grade == -1 && i == 0 {
		fmt.Println("No grades were entered, the program was existed.")
	} else {
		average = sum / float64(i)

		fmt.Println("The average grade is: ", average)
		switch {
		case average >= 90:
			fmt.Println("Score: A")
		case average < 90 && average >= 80:
			fmt.Println("Score: B")
		case average < 80 && average >= 70:
			fmt.Println("Score: C")
		case average < 70 && average >= 60:
			fmt.Println("Score: D")
		case average < 60:
			fmt.Println("Score: F")
		}
	}

}
