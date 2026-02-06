package main

import (
	"github.com/gutek00714/pokedexcli---Boot.dev/internal/pokecache"
	"time"
	"github.com/gutek00714/pokedexcli---Boot.dev/internal/pokeapi"
)

func main() {
	myCache := pokecache.NewCache(5 * time.Minute)
	cfg := &config{
		nextLocationsURL: "",
		previousLocationsURL: "",
		pokeCache: myCache,
		pokedex: make(map[string]pokeapi.Pokemon),
	}

	startRepl(cfg)
}