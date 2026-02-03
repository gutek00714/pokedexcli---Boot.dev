package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	return strings.Fields(lowered)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := cleanInput(scanner.Text())
		if len(line) == 0 {
			continue
		} else {
			first_word := line[0]
			fmt.Println("Your command was:", first_word)
		}
	}
}