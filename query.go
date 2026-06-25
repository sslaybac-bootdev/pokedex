package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type config struct {
	Next     *string
	Previous *string
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

func (l LocationAreaResponse) display() {
	for _, r := range l.Results {
		fmt.Printf("%s\n", r.Name)
	}

}

func getDefaultConfig() config {
	default_url := "https://pokeapi.co/api/v2/location-area"
	return config{
		Next:     &default_url,
		Previous: nil,
	}

}

func getLocations(url string) (*LocationAreaResponse, error) {
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
