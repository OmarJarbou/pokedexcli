package main

type LocationAreaMap struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []struct{ 
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next *string
	Previous *string
}