package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var commands map[string]cliCommand
var mapConfig config

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

func commandMap() error {
	res, err := http.Get(mapConfig.Next)
	if err != nil {
		return fmt.Errorf("failed to get response from endpoint: %v", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
	}
	if err != nil {
		return fmt.Errorf("error reading response %v", err)
	}
	areaMap := locationArea{}
	if err := json.Unmarshal(body, &areaMap); err != nil {
		return fmt.Errorf("error unmarshalling data: %v", err)
	}
	mapConfig.Previous = mapConfig.Next
	mapConfig.Next = areaMap.Next

	for _, result := range areaMap.Results {
		fmt.Println(result.Name)
	}
	return nil
}

type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
	Next     string
	Previous any
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
		"map": {
			name:        "map",
			description: "Displays 20 location areas",
			callback:    commandMap,
		},
	}

	mapConfig = config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}
}
