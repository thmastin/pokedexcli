package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchLocationAreas(url string) (LocationAreaResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to read body: %w", err)
	}

	var res LocationAreaResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return LocationAreaResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return res, nil
}

func FetchEncounter(url string, areaName string) (EncounterResponse, error) {

	resp, err := http.Get(url)
	if err != nil {
		return EncounterResponse{}, fmt.Errorf("failed to get response %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return EncounterResponse{}, fmt.Errorf("invalid area: %v. please use the pokedex 'map' command to see valid areas", areaName)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return EncounterResponse{}, fmt.Errorf("failed to read body: %w", err)
	}

	var res EncounterResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return EncounterResponse{}, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return res, nil
}

func FetchPokemon(url string, pokemonName string) (Pokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to get response %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return Pokemon{}, fmt.Errorf("invalid pokemon: %v. please use the pokedex 'explore' command to see valid pokemon", pokemonName)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to read body: %w", err)
	}

	var res Pokemon
	if err := json.Unmarshal(body, &res); err != nil {
		return Pokemon{}, fmt.Errorf("failed to unmarshal %w", err)
	}
	return res, nil
}
