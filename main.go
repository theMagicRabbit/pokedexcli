package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
	"io"
)

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type mapLocation struct {
	Name string
	Url  string
}

type mapResponse struct {
	Count 	 int
	Next	 string
	Previous string
	Results  []mapLocation
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func commandMap() error {
	uri := "https://pokeapi.co/api/v2/location-area"
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	jsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var locations mapResponse
	if err = json.Unmarshal(jsonBody, &locations); err != nil {
		fmt.Println(err)
		return err
	}
	for _, l := range locations.Results {
		fmt.Println(l.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

func main() {
	prompt := "Pokedex > "
	scanner := bufio.NewScanner(os.Stdin)
	commands = map[string]cliCommand {
		"exit": {
			name:		"exit",
			description:	"Exit the program",
			callback:	commandExit,
		},
		"help": {
			name:		"help",
			description:	"Print help",
			callback:	commandHelp,
		},
		"map": {
			name: 		"map",
			description:	"Page through locations in Pokemon",
			callback:	commandMap,
		},
	}
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
