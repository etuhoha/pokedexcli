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
	callback   func(*commandConfig, []string) error
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
		{
			name:       "explore",
			decription: "Explore an area",
			callback:   commandExplore,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var config commandConfig

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
				err := cmd.callback(&config, tokens[1:])
				if err != nil {
					fmt.Printf("ERROR: %v", err)
				}
				break
			}
		}
	}
}

func commandHelp(config *commandConfig, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")
	for _, cmd := range commands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.decription)
	}
	return nil
}

func commandExit(config *commandConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *commandConfig, args []string) error {
	url := config.nextURL
	if len(url) == 0 {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	return listMap(url, config)
}

func commandMapB(config *commandConfig, args []string) error {
	url := config.prevURL
	if len(url) == 0 {
		fmt.Println("you're on the first page")
		return nil
	}

	return listMap(url, config)
}

func listMap(url string, config *commandConfig) error {
	mapResp, err := pokeapi.Map(url)
	if err != nil {
		return err
	}

	for _, a := range mapResp.Areas {
		fmt.Printf("%v\n", a)
	}
	config.nextURL = mapResp.Next
	config.prevURL = mapResp.Prev
	return nil

}

func commandExplore(config *commandConfig, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing area name\n")
	}

	pokemons, err := pokeapi.ExploreArea(args[0])
	if err != nil {
		return err
	}

	for _, p := range pokemons {
		fmt.Printf("%v\n", p)
	}

	return nil
}
