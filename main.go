package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		inputCleaned := cleanInput(userInput)
		firstWord := inputCleaned[0]
		fmt.Printf("Your command was: %s\n", firstWord)
	}

}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	splitText := strings.Fields(lowerText)
	return splitText
}
