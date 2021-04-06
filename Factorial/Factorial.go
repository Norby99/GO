package main

import (
	"fmt"
)

func factorial(num int) int {
	if num == 0 {
		return 1
	}
	return num*factorial(num-1)
}

func main() {
	var num int
	fmt.Println("Enter your number: ")
	fmt.Scanln(&num)
	if num >= 0 {
		fmt.Printf("%d! = %d\n", num, factorial(num))
	}else{
		fmt.Println("Your number is invalid") // maybe one day I'll implement a func for the neg numbers and also for floating one
	}

	fmt.Println("Press \"Enter\" to close...")	// I use this to avoid the cmd from closing itself
	fmt.Scanln(&num)
}
