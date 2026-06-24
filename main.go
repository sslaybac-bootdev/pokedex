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
	scanner := bufio.NewScanner(os.Stdin)
	for promptAndScan(scanner) {
		input := scanner.Text()
		tok_input := cleanInput(input)
		fmt.Printf("Your command was: %s\n", tok_input[0])
	}
}
