package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func Commands(config *Config) map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 names of location areas in the Pokemon world",
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

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range Commands(config){
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap(config *Config) error {
	if config.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return fetchingLocationAreaMap(*(config.Next), config)
}

func commandMapb(config *Config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return fetchingLocationAreaMap(*(config.Previous), config)
}

func fetchingLocationAreaMap(url string, config *Config) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching location areas map: %w", err)
	}

	var locationAreaMap LocationAreaMap
	if err := json.NewDecoder(res.Body).Decode(&locationAreaMap); err != nil {
		return fmt.Errorf("error decoding location areas json data: %w", err)
	}

	for _, locationArea := range locationAreaMap.Results {
		fmt.Println(locationArea.Name)
	}

	if locationAreaMap.Previous != nil {
		config.Previous = locationAreaMap.Previous
	} else {
		config.Previous = nil
	}
	if locationAreaMap.Next != nil {
		config.Next = locationAreaMap.Next
	} else {
		config.Next = nil
	}
	return nil
}