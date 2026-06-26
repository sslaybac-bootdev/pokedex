package main

import (
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
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

func commandMap(cfg *config) error {
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

func commandMapb(cfg *config) error {
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

func commandHelp(cfg *config) error {
	commands := getCommands()
	for _, com := range commands {
		fmt.Printf("%s: %s\n", com.name, com.description)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
