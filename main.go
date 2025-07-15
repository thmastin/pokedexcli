package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thmastin/pokedexcli/internal/pokeapi"
	"github.com/thmastin/pokedexcli/internal/pokecache"
)

var commands map[string]cliCommand
var mapConfig config
var pokeCache pokecache.Cache

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
	config      *config
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
	config := commands["map"].config

	var areaMap pokeapi.LocationAreaResponse
	var err error

	cacheKey := *config.Next
	areaMap, err = fetchLocationAreaWithCache(cacheKey)
	if err != nil {
		return err
	}
	err = processLocationAreaResponse(areaMap, config)
	if err != nil {
		return err
	}
	return nil
}

func commandMapb() error {
	config := commands["mapb"].config

	var areaMap pokeapi.LocationAreaResponse
	var err error

	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	cacheKey := *config.Previous
	areaMap, err = fetchLocationAreaWithCache(cacheKey)
	if err != nil {
		return err
	}
	err = processLocationAreaResponse(areaMap, config)
	if err != nil {
		return err
	}
	return nil
}

type config struct {
	Next     *string
	Previous *string
}

func processLocationAreaResponse(areaMap pokeapi.LocationAreaResponse, config *config) error {
	if areaMap.Next != nil {
		config.Next = areaMap.Next
	} else {
		config.Next = nil
	}
	if areaMap.Previous != nil {
		config.Previous = areaMap.Previous
	} else {
		config.Previous = nil
	}
	for _, result := range areaMap.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func fetchLocationAreaWithCache(apiURL string) (pokeapi.LocationAreaResponse, error) {
	var areaMap pokeapi.LocationAreaResponse
	var err error

	cachedData, found := pokeCache.Get(apiURL)
	if found {
		err = json.Unmarshal(cachedData, &areaMap)
		if err != nil {
			return pokeapi.LocationAreaResponse{}, fmt.Errorf("error unmarshaling cached data: %w", err)
		}
	} else {

		areaMap, err = pokeapi.FetchLocationAreas(apiURL)
		if err != nil {
			return pokeapi.LocationAreaResponse{}, err
		}

		dataToCache, marshalErr := json.Marshal(areaMap)
		if marshalErr != nil {
			return pokeapi.LocationAreaResponse{}, fmt.Errorf("error marshaling data for cache: %w", marshalErr)
		}
		pokeCache.Add(apiURL, dataToCache)
	}
	return areaMap, nil

}

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      nil,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      nil,
		},
		"map": {
			name:        "map",
			description: "Displays 20 location areas",
			callback:    commandMap,
			config:      &mapConfig,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas",
			callback:    commandMapb,
			config:      &mapConfig,
		},
		"explore": {
			name:        "explore",
			description: "Displays the poke youman you can find in the area",
			callback:    commandExplore,
			config:      nil,
		},
	}
	mapStart := "https://pokeapi.co/api/v2/location-area/"
	mapConfig = config{
		Next:     &mapStart,
		Previous: nil,
	}
	pokeCache = pokecache.NewCache(5 * time.Minute)
}
