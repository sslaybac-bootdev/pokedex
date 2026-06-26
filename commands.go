package main

import (
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the next page of locations from the API",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of locations from the API",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Displays a list of pokemon that can be encountered at the target location",
			callback:    commandExplore,
		},
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

func commandMap(cfg *config, args ...string) error {
	if cfg.Next == nil {
		return errors.New("you're on the last page.")
	} else {
		locs, err := getLocations(*cfg.Next, cfg.cache)
		if err != nil {
			return err
		}

		locs.display()
		cfg.Next = locs.Next
		cfg.Previous = locs.Previous
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.Previous == nil {
		return errors.New("you're on the first page.")
	} else {
		locs, err := getLocations(*cfg.Previous, cfg.cache)
		if err != nil {
			return err
		}

		locs.display()
		cfg.Next = locs.Next
		cfg.Previous = locs.Previous
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	base_url := "https://pokeapi.co/api/v2/location-area/"
	location := args[0]
	fullUrl := base_url + location
	pokemon, err := getEncounters(fullUrl, cfg.cache)
	if err != nil {
		return err
	}
	pokemon.display()

	return nil
}

func commandHelp(cfg *config, args ...string) error {
	commands := getCommands()
	for _, com := range commands {
		fmt.Printf("%s: %s\n", com.name, com.description)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
