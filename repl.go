package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"github.com/gutek00714/pokedexcli---Boot.dev/internal/pokecache"
	"github.com/gutek00714/pokedexcli---Boot.dev/internal/pokeapi"
	"math/rand"
)

type cliCommand struct {
	name string
	description string
	callback func(cfg *config, name string) error
}

type config struct {
	nextLocationsURL string
	previousLocationsURL string
	pokeCache *pokecache.Cache
	pokedex map[string]pokeapi.Pokemon
}

var commands map[string]cliCommand

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
		"explore": {
			name: "explore",
			description: "Display a list of all Pokemon in the location",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "Try to catch a Pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Inspect Pokemon stats",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "Show caught Pokemon",
			callback: commandPokedex,
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
				var name string
				if len(line) > 1 {
					name = line[1]
				}
				err := cmd.callback(cfg, name)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func commandExit(cfg *config, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, _ string) error {
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

func commandMap(cfg *config, _ string) error {
	var url string

	// Check if next url exists or use default one
	if cfg.nextLocationsURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		url = cfg.nextLocationsURL
	}

	//get
	var issues pokeapi.Issue
	// fmt.Println(url)
	data, found := cfg.pokeCache.Get(url)
	if found {
		if err := json.Unmarshal(data, &issues); err != nil {
			return err
		}
	} else {
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

		//add
		cfg.pokeCache.Add(url, body)
		
		// Create object from the call
		if err := json.Unmarshal(body, &issues); err != nil {
			return err
		}
	}

	for _, loc := range issues.Results {
		fmt.Println(loc.Name)
	}

	// Update pointer URL
	cfg.nextLocationsURL = issues.Next
	cfg.previousLocationsURL = issues.Previous

	return nil
}

func commandMapb(cfg *config, _ string) error {
	var url string

	// Check if next url exists or use default one
	if cfg.previousLocationsURL == "" {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		url = cfg.previousLocationsURL
	}

	//get
	var issues pokeapi.Issue
	// fmt.Println(url)
	data, found := cfg.pokeCache.Get(url)
	if found {
		if err := json.Unmarshal(data, &issues); err != nil {
			return err
		}
	} else {
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

		//add
		cfg.pokeCache.Add(url, body)
		
		// Create object from the call
		if err := json.Unmarshal(body, &issues); err != nil {
			return err
		}
	}

	for _, loc := range issues.Results {
		fmt.Println(loc.Name)
	}

	// Update pointer URL
	cfg.nextLocationsURL = issues.Next
	cfg.previousLocationsURL = issues.Previous

	return nil
}

func commandExplore(cfg *config, name string) error {
	if len(name) == 0 {
		return fmt.Errorf("you must provide a location name")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + name

	fmt.Printf("Exploring %v...\n", name)

	var locationArea pokeapi.LocationArea

	// Check cache get
	data, found := cfg.pokeCache.Get(url)
	if found {
		if err := json.Unmarshal(data, &locationArea); err != nil {
			return err
		}
	} else {
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

		// Add to cache
		cfg.pokeCache.Add(url, body)

		// Create object from the call
		if err := json.Unmarshal(body, &locationArea); err != nil {
			return err
		}
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, name string) error {
	if len(name) == 0 {
		return fmt.Errorf("you must provide a Pokemon name")
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	fmt.Printf("Throwing a Pokeball at %v...\n", name)

	var pokemon pokeapi.Pokemon

	// Check cache get
	data, found := cfg.pokeCache.Get(url)
	if found {
		if err := json.Unmarshal(data, &pokemon); err != nil {
			return err
		}
	} else {
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

		// Add to cache
		cfg.pokeCache.Add(url, body)

		// Create object from the call
		if err := json.Unmarshal(body, &pokemon); err != nil {
			return err
		}
	}

	// Calculate pokemon catch rate
	catch_treshold := 30
	random := rand.Intn(pokemon.BaseExperience)
	if random > catch_treshold {
		fmt.Printf("%v escaped!\n", name)
	} else {
		fmt.Printf("%v was caught!\n", name)
		cfg.pokedex[name] = pokemon
		fmt.Println("You may now inspect it with the inspect command.")
	}
	return nil
}

func commandInspect(cfg *config, name string) error {
	if len(name) == 0 {
		return fmt.Errorf("you must provide a Pokemon name")
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	var pokemon pokeapi.Pokemon

	// Check cache get
	data, found := cfg.pokeCache.Get(url)
	if found {
		if err := json.Unmarshal(data, &pokemon); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("you have not caught that pokemon") 
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %v\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, _ string) error {
	fmt.Println("Your Pokedex:")
	for pokemon := range cfg.pokedex {
		fmt.Printf(" - %v\n", pokemon)
	}
	return nil
}