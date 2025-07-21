package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/OmarJarbou/pokedexcli/internal/pokecache"
	"io"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache) error
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

func commandExit(config *Config, cache *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *pokecache.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range Commands(config){
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap(config *Config, cache *pokecache.Cache) error {
	if config.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return fetchingLocationAreaMap(*(config.Next), config, cache)
}

func commandMapb(config *Config, cache *pokecache.Cache) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return fetchingLocationAreaMap(*(config.Previous), config, cache)
}

func fetchingLocationAreaMap(url string, config *Config, cache *pokecache.Cache) error {
	var locationAreaMap LocationAreaMap

	var Response []byte
	cachedRes, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching location areas map: %w", err)
		}
		Response, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading location areas json data from response: %w", err)
		}
		cache.Add(url, Response)
		fmt.Println("DATA FETCHED FROM INTERNET")
	} else {
		Response = cachedRes
		fmt.Println("DATA FOUND IN THE CACHE")
	}

	if err := json.Unmarshal(Response, &locationAreaMap); err != nil {
		return fmt.Errorf("error parsing location areas json-encoded data: %w", err)
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