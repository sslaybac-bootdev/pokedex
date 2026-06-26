package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedex/internal/pokecache"
	"time"
)

type config struct {
	Next     *string
	Previous *string
	cache    *pokecache.Cache
	pokedex  map[string]Pokemon
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    *int           `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationSpecificResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func (l *LocationAreaResponse) display() {
	for _, r := range l.Results {
		fmt.Printf("%s\n", r.Name)
	}

}

func (l *LocationSpecificResponse) display() {
	for _, r := range l.PokemonEncounters {
		fmt.Printf("%s\n", r.Pokemon.Name)
	}

}

func getDefaultConfig() config {
	default_url := "https://pokeapi.co/api/v2/location-area"
	cache := pokecache.NewCache(5 * time.Second)
	dex := make(map[string]Pokemon, 0)
	return config{
		Next:     &default_url,
		Previous: nil,
		cache:    cache,
		pokedex:  dex,
	}
}

func query(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status code: %d (url: %s)", res.StatusCode, url)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getLocations(url string, cache *pokecache.Cache) (*LocationAreaResponse, error) {
	data, ok := cache.Get(url)
	if ok {
		return parse_locations(data)
	}
	body, err := query(url)
	if err != nil {
		return nil, err
	}
	return parse_locations(body)
}

func getEncounters(url string, cache *pokecache.Cache) (*LocationSpecificResponse, error) {
	data, ok := cache.Get(url)
	if ok {
		return parse_encounters(data)
	}
	body, err := query(url)
	if err != nil {
		return nil, err
	}
	return parse_encounters(body)
}

func getPokemon(url string, cache *pokecache.Cache) (*Pokemon, error) {
	data, ok := cache.Get(url)
	if ok {
		return parse_pokemon(data)
	}
	body, err := query(url)
	if err != nil {
		return nil, err
	}
	return parse_pokemon(body)
}

func parse_locations(data []byte) (*LocationAreaResponse, error) {
	var locations LocationAreaResponse

	err := json.Unmarshal(data, &locations)
	if err != nil {
		return nil, err
	}

	return &locations, nil
}

func parse_encounters(data []byte) (*LocationSpecificResponse, error) {
	var encounters LocationSpecificResponse

	err := json.Unmarshal(data, &encounters)
	if err != nil {
		return nil, err
	}

	return &encounters, nil
}

func parse_pokemon(data []byte) (*Pokemon, error) {
	var pokemon Pokemon

	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return nil, err
	}

	return &pokemon, nil
}
