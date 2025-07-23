package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/OmarJarbou/pokedexcli/internal/pokecache"
	"io"
	"math/rand"
	"strconv"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache, []string, *Pokedex) error
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
		"explore": {
			name:        "explore",
			description: "Displays the names of pokemons located in specific location area in the Pokemon world",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catches a Pokemon and adds it to the user's Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows details about only a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints a list of all the names of the Pokemon the user has caught",
			callback:    commandPokedex,
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

func commandExit(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 1 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 0 arguments, but found " + strconv.Itoa(foundArguments))
		return nil
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 1 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 0 arguments, but found " + strconv.Itoa(foundArguments))
		return nil
	}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range Commands(config){
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 1 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 0 arguments, but found " + strconv.Itoa(foundArguments))
		return nil
	}
	if config.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return fetchingLocationAreaMap(*(config.Next), config, cache)
}

func commandMapb(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 1 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 0 arguments, but found " + strconv.Itoa(foundArguments))
		return nil
	}
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

func commandExplore(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 2 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 1 argument, but found " + strconv.Itoa(foundArguments))
		return nil
	}
	fmt.Println("Exploring " + commandWords[1] + "...")

	fullURL := "https://pokeapi.co/api/v2/location-area/" + commandWords[1] + "/"

	var Response []byte
	cachedRes, ok := cache.Get(fullURL)
	if !ok {
		res, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error fetching this location area's data: %w", err)
		}

		Response, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading location area's json data from response: %w", err)
		}

		cache.Add(fullURL, Response)
		fmt.Println("DATA FETCHED FROM INTERNET")
	} else {
		Response = cachedRes
		fmt.Println("DATA FOUND IN THE CACHE")
	}
	

	var locationArea LocationArea
	if err := json.Unmarshal(Response, &locationArea); err != nil {
		return fmt.Errorf("error parsing this location area's json-encoded data: %w", err)
	}

	fmt.Println("Found Pokemon:")
	pokemons := locationArea.PokemonEncounters
	for _, item := range pokemons {
		fmt.Println(" - " + item.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 2 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 1 argument, but found " + strconv.Itoa(foundArguments))
		return nil
	}
	fmt.Println("Throwing a Pokeball at " + commandWords[1] + "...")

	fullURL := "https://pokeapi.co/api/v2/pokemon/" + commandWords[1] + "/"

	pokemon, err := fetchPokemon(fullURL, cache)
	if err != nil {
		return err
	}

	// pokemon.BaseExperience
	randomNumber := rand.Intn(400)
	if randomNumber > pokemon.BaseExperience {
		fmt.Println(commandWords[1] + " was caught!")
		pokedex.Add(pokemon)
	} else {
		fmt.Println(commandWords[1] + " escaped!")
	}

	return nil
}

func commandInspect(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 2 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 1 argument, but found " + strconv.Itoa(foundArguments))
		return nil
	}

	found := false
	for _, poke := range pokedex.Items {
		if poke.Name == commandWords[1] {
			found = true
			break
		}
	} 
	if !found {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fullURL := "https://pokeapi.co/api/v2/pokemon/" + commandWords[1] + "/"

	pokemon, err := fetchPokemon(fullURL, cache)
	if err != nil {
		return err
	}

	fmt.Println("Name: " + pokemon.Name)
	fmt.Println("Height: " + strconv.Itoa(pokemon.Height))
	fmt.Println("Weight: " + strconv.Itoa(pokemon.Weight))
	fmt.Println("Stats:")
	for _, item := range pokemon.Stats {
		fmt.Println("  -" + item.Stat.Name + ": " + strconv.Itoa(item.BaseStat))
	}
	fmt.Println("Types:")
	for _, item := range pokemon.Types {
		fmt.Println("  -" + item.Type.Name)
	}

	return nil
}

func fetchPokemon(fullURL string, cache *pokecache.Cache) (Pokemon, error) {
	var pokemon Pokemon

	var Response []byte
	cachedData, ok := cache.Get(fullURL)
	if !ok {
		res, err := http.Get(fullURL)
		if err != nil {
			return pokemon, fmt.Errorf("error fetching this pokemon's data: %w", err)
		}

		Response, err = io.ReadAll(res.Body)
		if err != nil {
			return pokemon, fmt.Errorf("error reading pokemons's json data from response: %w", err)
		}

		cache.Add(fullURL, Response)
		fmt.Println("DATA FETCHED FROM INTERNET")
	} else {
		Response = cachedData
		fmt.Println("DATA FOUND IN THE CACHE")
	}
	
	if err := json.Unmarshal(Response, &pokemon); err != nil {
		return pokemon, fmt.Errorf("error parsing this pokemon's json-encoded data: %w", err)
	}

	return pokemon, nil
}

func commandPokedex(config *Config, cache *pokecache.Cache, commandWords []string, pokedex *Pokedex) error {
	if len(commandWords) != 1 {
		foundArguments := len(commandWords) - 1
		fmt.Println("Expected 0 arguments, but found " + strconv.Itoa(foundArguments))
		return nil
	}

	if len(pokedex.Items) == 0 {
		fmt.Println("Your Pokedex is empty!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range pokedex.Items {
		fmt.Println(" - " + pokemon.Name)
	} 

	return nil
}