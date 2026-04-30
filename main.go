package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		fmt.Printf("Your command was: %v\n", tokens[0])
	}
}
