package main

import ( 
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
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