package pokeapi

type LocationAreaResponse struct {
	Count    int     "json: count"
	Next     *string "jsong: next"
	Previous *string "jsgon: previous"
	Results  []struct {
		Name string "jsaon: name"
		URL  string `json: "url"`
	} "json: results"
}
