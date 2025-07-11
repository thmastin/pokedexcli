package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	startREPL((scanner))

}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	splitText := strings.Fields(lowerText)
	return splitText
}

func getFirstWord(words []string) string {
	if len(words) == 0 {
		return ""
	}
	return words[0]
}

func displayOutput(word string) string {
	if word == "" {
		return "Please enter a command\n"
	}
	return fmt.Sprintf("Your command was: %s\n", word)
}

func startREPL(scanner *bufio.Scanner) {
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		inputCleaned := cleanInput(userInput)
		firstWord := getFirstWord((inputCleaned))
		fmt.Print(displayOutput((firstWord)))
	}
}
