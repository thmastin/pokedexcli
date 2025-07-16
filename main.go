package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/thmastin/pokedexcli/internal/pokeapi"
	"github.com/thmastin/pokedexcli/internal/pokecache"
)

var commands map[string]cliCommand
var mapConfig config
var pokeCache pokecache.Cache
var encounterBaseUrl string
var catchBaseUrl string
var pokedex map[string]pokeapi.Pokemon
var rng *rand.Rand
var catchAttempts map[string]int

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

func getSecondWord(words []string) string {
	if len(words) < 2 {
		return ""
	}
	return words[1]
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
		userInput, seconduserInput := processCommand(scanner.Text())
		command, exists := commands[userInput]
		if exists {
			err := command.callback(seconduserInput)
			if err != nil {
				fmt.Printf("Error executing %v command: %v\n", userInput, err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func processCommand(userInput string) (string, string) {
	inputCleaned := cleanInput(userInput)
	firstWord := getFirstWord(inputCleaned)
	secondWord := getSecondWord(inputCleaned)
	return firstWord, secondWord
}

func commandExit(seconduserInput string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
	config      *config
}

func commandHelp(seconduserInput string) error {
	fmt.Println(helpMessage())
	return nil
}

func helpMessage() string {
	newMessage := `Welcome to the Pokedex!
Usage:`
	newMessage += "\n\n"
	for _, value := range commands {
		commandDescription := fmt.Sprintf("%v: %v\n", value.name, value.description)
		newMessage += commandDescription
	}

	return newMessage

}

func commandMap(seconduserInput string) error {
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

func commandMapb(seconduserInput string) error {
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

func commandExplore(seconduserInput string) error {
	if seconduserInput == "" {
		return fmt.Errorf("please enter an area name")
	}

	cacheKey := encounterBaseUrl + seconduserInput

	cachedData, found := pokeCache.Get(cacheKey)
	if found {
		var encounter pokeapi.EncounterResponse
		err := json.Unmarshal(cachedData, &encounter)
		if err != nil {
			return fmt.Errorf("invalid area: %s, please use the Pokedex 'map' command to see valid areas", seconduserInput)
		}
		err = processEncounterResponse(encounter, seconduserInput)
		if err != nil {
			return err
		}
		return nil
	}

	encounter, err := fetchEncounterWithCache(cacheKey, seconduserInput)
	if err != nil {
		return err
	}
	err = processEncounterResponse(encounter, seconduserInput)
	if err != nil {
		return err
	}
	return nil
}

func commandCatch(seconduserInput string) error {
	var pokemon pokeapi.Pokemon
	if seconduserInput == "" {
		return fmt.Errorf("please enter a pokemon name")
	}

	cacheKey := catchBaseUrl + seconduserInput

	cachedData, found := pokeCache.Get(cacheKey)
	if found {
		err := json.Unmarshal(cachedData, &pokemon)
		if err != nil {
			return fmt.Errorf("invalid pokemon: %s, please use the Pokedex 'explore' command to see valid pokemon", seconduserInput)
		}
		err = processCatchResponse(pokemon, seconduserInput)
		if err != nil {
			return err
		}
		return nil
	}

	pokemon, err := fetchCatchWithCache(cacheKey, seconduserInput)
	if err != nil {
		return err
	}
	err = processCatchResponse(pokemon, seconduserInput)
	if err != nil {
		return err
	}
	return nil
}

type config struct {
	Next     *string
	Previous *string
	Results  []pokeapi.LocationArea
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
	config.Results = areaMap.Results
	for _, result := range areaMap.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func processEncounterResponse(encounter pokeapi.EncounterResponse, areaName string) error {
	if encounter.PokemonEncounters == nil {
		fmt.Printf("No Pokemon found in %v", areaName)
		return nil
	}

	fmt.Printf("Exploring %v...\n", areaName)
	fmt.Println("Found Pokemon:")
	for _, encounterEntry := range encounter.PokemonEncounters {
		fmt.Printf(" - %v\n", encounterEntry.Pokemon.Name)
	}
	return nil
}

func processCatchResponse(pokemon pokeapi.Pokemon, pokemonName string) error {
	if _, ok := pokedex[pokemonName]; ok {
		return fmt.Errorf("you've already caught %s", pokemonName)
	} else {
		fmt.Printf("Throwing a Pokeball at %s...", pokemonName)
		pokemonExperience := pokemon.BaseExperience
		if catchAttempts[pokemon.Name] == 2 {
			pokemonCatch(pokemon, pokemonName)
		} else {
			chance := rng.Float64() * 100
			if chance > float64(pokemonExperience/2) {
				pokemonCatch(pokemon, pokemonName)
			} else {
				fmt.Printf("%s escaped!\n", pokemonName)
				catchAttempts[pokemon.Name] += 1
				fmt.Printf("Attempted catch: %v\n", catchAttempts[pokemon.Name])

			}
		}
		return nil
	}
}

func pokemonCatch(catch pokeapi.Pokemon, pokemonName string) {
	pokedex[catch.Name] = catch
	fmt.Printf("%s was caught!\n", pokemonName)
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

func fetchEncounterWithCache(apiURL string, areaName string) (pokeapi.EncounterResponse, error) {
	var encounter pokeapi.EncounterResponse
	var err error

	cachedData, found := pokeCache.Get(apiURL)
	if found {
		err = json.Unmarshal(cachedData, &encounter)
		if err != nil {
			return pokeapi.EncounterResponse{}, fmt.Errorf("error unmarshaling cached data: %w", err)
		}
	} else {
		encounter, err = pokeapi.FetchEncounter(apiURL, areaName)
		if err != nil {
			return pokeapi.EncounterResponse{}, err
		}

		dataToCache, marshalErr := json.Marshal(encounter)
		if marshalErr != nil {
			return pokeapi.EncounterResponse{}, fmt.Errorf("error marshaling data for cache: %w", marshalErr)
		}
		pokeCache.Add(apiURL, dataToCache)
	}
	return encounter, nil
}

func fetchCatchWithCache(apiURL string, pokemonName string) (pokeapi.Pokemon, error) {
	var pokemon pokeapi.Pokemon
	var err error

	cachedData, found := pokeCache.Get(apiURL)
	if found {
		err = json.Unmarshal(cachedData, &pokemon)
		if err != nil {
			return pokeapi.Pokemon{}, fmt.Errorf("error unmarshaling cached data: %w", err)
		}
	} else {
		pokemon, err = pokeapi.FetchPokemon(apiURL, pokemonName)
		if err != nil {
			return pokeapi.Pokemon{}, err
		}
		dataToCache, marshalErr := json.Marshal(pokemon)
		if marshalErr != nil {
			return pokeapi.Pokemon{}, fmt.Errorf("error marshaling data for cache: %w", marshalErr)
		}
		pokeCache.Add(apiURL, dataToCache)
	}
	return pokemon, nil
}

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays this help message",
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
			name:        "explore <area_name>",
			description: "Displays the poke youman you can find in the area",
			callback:    commandExplore,
			config:      nil,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempts to catch a Pokemon",
			callback:    commandCatch,
			config:      nil,
		},
	}
	mapStart := "https://pokeapi.co/api/v2/location-area/"
	mapConfig = config{
		Next:     &mapStart,
		Previous: nil,
		Results:  []pokeapi.LocationArea{},
	}
	pokeCache = pokecache.NewCache(5 * time.Minute)
	encounterBaseUrl = "https://pokeapi.co/api/v2/location-area/"
	catchBaseUrl = "https://pokeapi.co/api/v2/pokemon/"
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	pokedex = make(map[string]pokeapi.Pokemon)
	catchAttempts = make(map[string]int)
}
