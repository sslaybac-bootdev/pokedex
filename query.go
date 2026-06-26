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

type Pokemon struct {
	name *string `json:name`
	url  *string `json:url`
}

func (l *LocationAreaResponse) display() {
	for _, r := range l.Results {
		fmt.Printf("%s\n", r.Name)
	}

}

func getDefaultConfig() config {
	default_url := "https://pokeapi.co/api/v2/location-area"
	cache := pokecache.NewCache(5 * time.Second)
	return config{
		Next:     &default_url,
		Previous: nil,
		cache:    cache,
	}

}

func getLocations(url string, cache *pokecache.Cache) (*LocationAreaResponse, error) {
	data, ok := cache.Get(url)
	if ok {
		return parse_locations(data)
	}
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
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var locations LocationAreaResponse

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, err
	}

	return &locations, nil

}

func parse_locations(data []byte) (*LocationAreaResponse, error) {
	var locations LocationAreaResponse

	err := json.Unmarshal(data, &locations)
	if err != nil {
		return nil, err
	}

	return &locations, nil
}
