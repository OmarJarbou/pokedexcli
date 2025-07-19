package main

import (
	"testing"
	"strings"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	} {
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input: " Charmander Bulbasaur PIKACHU ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input: " h e l l  o wo r ll d ",
			expected: []string{"h", "e", "l", "l", "o", "wo", "r", "ll", "d"},
		},
	}
	for _, c := range cases {
		result := cleanInput(c.input)
		if len(result) < len(c.expected) {
			t.Errorf("Failed parsing input: cleanInput function's result is less than expected.\nExpected: %v\nYour's: %v", c.expected, result)
			return
		}
		// slices are uncomparable:
		// if result != c.expected : wrong
		for i := 0; i < len(result); i++ {
			res := result[i]
			exp := strings.ToLower(c.expected[i])
			if res != exp {
				t.Errorf("Failed parsing input: one of cleanInput function's results is not as expected.\nExpected: %v\nYour's: %v", c.expected, result)
				return
			}
		}
	}
}
