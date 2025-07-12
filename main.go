package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands map[string]cliCommand

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
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading input %v\n", err)
			}
			break
		}
		userInput := processCommand(scanner.Text())
		command, exists := commands[userInput]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Printf("Error executing exit command: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func processCommand(userInput string) string {
	inputCleaned := cleanInput(userInput)
	firstWord := getFirstWord((inputCleaned))
	return firstWord
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println(helpMessage())
	return nil
}

func helpMessage() string {
	newMessage := `Welcome to the Pokedex!
Usage:`
	newMessage += "\n\n"
	for key, value := range commands {
		commandDescription := fmt.Sprintf("%v: %v\n", key, value.description)
		newMessage += commandDescription
	}

	return newMessage

}

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
