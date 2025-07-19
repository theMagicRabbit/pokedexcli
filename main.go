package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

var commands map[string]cliCommand = map[string]cliCommand {
	"exit": {
		name:		"exit",
		description:	"Exit the program",
		callback:	commandExit,
	},
}
func main() {
	prompt := "Pokedex > "
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		cmd, ok := commands[cleanedInput[0]]
		if ok {
			cmd.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}


