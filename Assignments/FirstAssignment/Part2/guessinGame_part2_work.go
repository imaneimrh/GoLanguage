package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var HighestScore int

func main() {

	RandomInteger := rand.Intn(100)
	var guess int
	var numberTries int
	var counter int
	var playAgain = true
	var won bool = false
	//fmt.Println(RandomInteger)

	for playAgain != false {
		fmt.Println("Hello, and welcome to the Great Guessing Game!!")

		fmt.Println("Choose the number of tries in the game: ")
		fmt.Scanln(&numberTries)
		fmt.Println("\n---Rules: \nyou guess Until You run out of chances! \nYour guesses should be integers within the range 1..100:\n---Start!")

		for won != true && counter < numberTries {
			fmt.Println("Please enter your guess: ")
			fmt.Scanln(&guess)
			counter++
			switch {
			case guess == RandomInteger:
				won = true
				fmt.Println("Congratulations!! you Won!!")
				fmt.Println("Your score is : ", numberTries)
				if HighestScore < numberTries {
					HighestScore = numberTries
				}
				fmt.Println("Your Highest score is ", HighestScore)
			case guess < 1 || guess > 100:
				fmt.Println("Invalid input, make sure that you enter a number between 1 and 100")

			case guess < RandomInteger:
				fmt.Println("Too Low")
			case guess > RandomInteger:
				fmt.Println("Too High")
			case counter < numberTries:
				fmt.Println("You lost! you exhausted your tries!")
			default:
				fmt.Println("Invalid input!")
			}
			fmt.Println("You lost! you exhausted your tries!")
		}
		fmt.Println("Do you want to play again? (Yes/No)")
		var choice string
		fmt.Scanln(&choice)
		choice = strings.ToLower(choice)
		switch choice {
		case "yes":
			playAgain = true
		case "no":
			playAgain = false
		}
	}
}
