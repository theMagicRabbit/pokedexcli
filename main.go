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
	callback    func(*config) error
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

type config struct {
	MapNextUrl, MapPreviousUrl string
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func commandMap(conf *config) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}
	var uri string
	if conf.MapNextUrl != "" {
		uri = conf.MapNextUrl
	} else {
		uri = "https://pokeapi.co/api/v2/location-area" 
	}
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
	conf.MapNextUrl = locations.Next
	conf.MapPreviousUrl = locations.Previous
	for _, l := range locations.Results {
		fmt.Println(l.Name)
	}

	return nil
}

func commandMapBack(conf *config) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}
	var uri string
	if conf.MapPreviousUrl != "" {
		uri = conf.MapPreviousUrl
	} else {
		fmt.Println("you're on the first page")
		return nil
	}
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
	conf.MapNextUrl = locations.Next
	conf.MapPreviousUrl = locations.Previous
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
	conf := config{}
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
		"mapb": {
			name: 		"mapb",
			description:	"Go to previous locations page in Pokemon",
			callback:	commandMapBack,
		},
	}
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		cmd, ok := commands[cleanedInput[0]]
		if ok {
			cmd.callback(&conf)
		} else {
			fmt.Println("Unknown command")
		}
	}

}
