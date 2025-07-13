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
