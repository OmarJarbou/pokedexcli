package main

import ( 
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	initURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	config := Config{
		Next: &initURL,
		Previous: nil,
	}
	commands := Commands(&config)
	scanner := bufio.NewScanner(os.Stdin) 
	
	for programStartingREPL(scanner) {
		var inputString string
		var inputWords [] string
		
		inputString = scanner.Text() // get the line we read as a string

		inputWords = cleanInput(inputString)

		// fmt.Println("Your command was: " + inputWords[0])

		command, ok := commands[inputWords[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		command.callback(&config)
	}
	if err := scanner.Err(); err != nil { // if err occured during scanning
		fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
	}
}

func programStartingREPL(scanner *bufio.Scanner) bool {
	fmt.Print("Pokedex > ")
	return scanner.Scan() // scan based on the rules of "scanner": read a line
}

func cleanInput(text string) []string {
	var result []string

	var temp string
	temp = ""
	for i, c := range text {
		if string(c) == " " {
			if temp != "" {
				result = append(result, strings.ToLower(temp))
				temp = ""
			}
			continue
		}

		temp = temp + string(c)

		if i + 1 == len(text) {
			if temp != "" {
				result = append(result, strings.ToLower(temp))
				temp = ""
			}
			continue
		}
	}

	return result
}