package main

import (
	"bufio"
	"fmt"
	"os"
)

func promptAndScan(scanner *bufio.Scanner) bool {
	fmt.Print("Pokedex > ")
	return scanner.Scan()
}

func main() {
	fmt.Println("Welcome to the Pokedex!")
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for promptAndScan(scanner) {
		input := scanner.Text()
		tok_input := cleanInput(input)
		action, exists := commands[tok_input[0]]
		if !exists {
			fmt.Println("Unknown command")
		} else {
			action.callback()
		}
	}
}
