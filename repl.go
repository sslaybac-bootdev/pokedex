package main

import "strings"

// Converts an input sting into a list of tokens
// splits on whitespace, then sets all to lowercase
func cleanInput(text string) []string {
	tokens := strings.Fields(strings.ToLower(text))
	return tokens
}
