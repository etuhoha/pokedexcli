package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/etuhoha/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name       string
	decription string
	callback   func(*commandConfig) error
}

type commandConfig struct {
	nextURL string
	prevURL string
}

func commands() []cliCommand {
	return []cliCommand{
		{
			name:       "help",
			decription: "Displays a help message",
			callback:   commandHelp,
		},
		{
			name:       "exit",
			decription: "Exit the Pokedex",
			callback:   commandExit,
		},
		{
			name:       "map",
			decription: "List next location areas",
			callback:   commandMap,
		},
		{
			name:       "mapb",
			decription: "List previous location areas",
			callback:   commandMapB,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var cmdConf commandConfig

	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			break
		}

		text := scanner.Text()
		tokens := cleanInput(text)

		for _, cmd := range commands() {
			if tokens[0] == cmd.name {
				err := cmd.callback(&cmdConf)
				if err != nil {
					fmt.Printf("ERROR: %v", err)
				}
				break
			}
		}
	}
}

func commandHelp(*commandConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")
	for _, cmd := range commands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.decription)
	}
	return nil
}

func commandExit(*commandConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *commandConfig) error {
	url := config.nextURL
	if len(url) == 0 {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	return listMap(url, config)
}

func commandMapB(config *commandConfig) error {
	url := config.prevURL
	if len(url) == 0 {
		fmt.Println("you're on the first page")
		return nil
	}

	return listMap(url, config)
}

func listMap(url string, config *commandConfig) error {
	next, prev, err := pokeapi.Map(url)
	if err != nil {
		return err
	}

	config.nextURL = next
	config.prevURL = prev
	return nil

}
