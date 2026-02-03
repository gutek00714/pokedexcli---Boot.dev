package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
	"net/http"
	"io"
	"encoding/json"
)

type cliCommand struct {
	name string
	description string
	callback func(cfg *config) error
}

var commands map[string]cliCommand

// var commands = map[string]cliCommand{
// 	"exit": {
// 		name:        "exit",
// 		description: "Exit the Pokedex",
// 		callback:    commandExit,
// 	},
// 	"help": {
// 		name: "help",
// 		description: "Displays a help message",
// 		callback: commandHelp,
// 	},
// }

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	return strings.Fields(lowered)
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Displays the names of 20 location areas",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Display the names of previous 20 location areas",
			callback: commandMapb,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := cleanInput(scanner.Text())
		if len(line) == 0 {
			continue
		} else {
			// get first word
			first_word := line[0]

			// check if the key(first_word) is in map
			cmd, exists := commands[first_word]
			if exists {
				err := cmd.callback(cfg)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	// fmt.Printf("%s: %s\n", commands["help"].name, commands["help"].description)
	// fmt.Printf("%s: %s\n", commands["exit"].name, commands["exit"].description)
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

type Issue struct {
	Next string
	Previous string
	Results []Result
}

type Result struct {
	Name string
}

type config struct {
	nextLocationsURL string
	previousLocationsURL string
}

func commandMap(cfg *config) error {
	var url string

	// Check if next url exists or use default one
	if cfg.nextLocationsURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = cfg.nextLocationsURL
	}

	// Call API
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	
	// Create object from the call
	var issues Issue
	if err := json.Unmarshal(body, &issues); err != nil {
		return err
	}

	for _, loc := range issues.Results {
		fmt.Println(loc.Name)
	}

	// Update pointer URL
	cfg.nextLocationsURL = issues.Next
	cfg.previousLocationsURL = issues.Previous

	return nil
}

func commandMapb(cfg *config) error {
	var url string

	// Check if next url exists or use default one
	if cfg.previousLocationsURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = cfg.previousLocationsURL
	}

	// Call API
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	
	// Create object from the call
	var issues Issue
	if err := json.Unmarshal(body, &issues); err != nil {
		return err
	}

	for _, loc := range issues.Results {
		fmt.Println(loc.Name)
	}

	// Update pointer URL
	cfg.nextLocationsURL = issues.Next
	cfg.previousLocationsURL = issues.Previous

	return nil
}