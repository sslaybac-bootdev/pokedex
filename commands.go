package main

import (
	"errors"
	"fmt"
	"math/rand"
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
		"catch": {
			name:        "catch <pokemon>",
			description: "Attempts to capture the target pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Prints statistics for the target pokemon",
			callback:    commandInspect,
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

func commandCatch(cfg *config, args ...string) error {
	base_url := "https://pokeapi.co/api/v2/pokemon/"
	targetPokemon := args[0]
	full_url := base_url + targetPokemon
	pokemonEntry, err := getPokemon(full_url, cfg.cache)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonEntry.Name)

	//chance to catch is 100/(100+baseExperience)
	catch_chance := 99 + pokemonEntry.BaseExperience
	catch_roll := rand.Intn(catch_chance)
	if catch_roll < 100 {
		fmt.Printf("%s escaped!\n", pokemonEntry.Name)
	} else {
		fmt.Printf("%s was caught!\n", pokemonEntry.Name)
		cfg.pokedex[pokemonEntry.Name] = *pokemonEntry
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	targetPokemon := args[0]
	pokemonEntry, exists := cfg.pokedex[targetPokemon]
	if !exists {
		fmt.Println("You have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemonEntry.Name)
	fmt.Printf("Height: %d\n", pokemonEntry.Height)
	fmt.Printf("Weight: %d\n", pokemonEntry.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemonEntry.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemonEntry.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}

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
	fmt.Println("Pokedex contains:")
	for key := range cfg.pokedex {
		fmt.Printf("\t%s\n", key)
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
