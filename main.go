package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name       string
	decription string
	callback   func() error
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			fmt.Printf("scanner error: %v\n", scanner.Err())
			break
		}

		text := scanner.Text()
		tokens := cleanInput(text)

		for _, cmd := range commands() {
			if tokens[0] == cmd.name {
				cmd.callback()
				break
			}
		}
	}
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println("")
	for _, cmd := range commands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.decription)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
