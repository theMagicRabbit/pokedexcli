package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"math/rand"
	"encoding/json"
	"io"
	"time"
	"github.com/theMagicRabbit/pokedexcli/internal"
)

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
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
	CacheApi *internal.Cache
	Pokedex map[string]internal.PokemonResponse
}

type encounterResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandCatch(conf *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Must provide a Pokemon to catch")
	}
	baseUri := "https://pokeapi.co/api/v2/pokemon"
	jsonBytes, err := cacheAwareGet(fmt.Sprintf("%s/%s", baseUri, args[0]), conf)
	if err != nil {
		return err
	}
	var pokemon internal.PokemonResponse
	err = json.Unmarshal(jsonBytes, &pokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing Pokeball at %s...\n", args[0])
	targetScore := rand.Int() % (pokemon.BaseExperience / 2)
	pokeScore := rand.Int() % pokemon.BaseExperience
	if pokeScore <= targetScore {
		pokedex := conf.Pokedex
		pokedex[args[0]] = pokemon
		fmt.Printf("caught %s!\n", args[0])
	} else { 
		fmt.Printf("%s escaped!\n", args[0])
	}
	return nil
}

func commandExit(conf *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for k, v := range commands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return nil
}

func cacheAwareGet(uri string, conf *config) ([]byte, error) {
	var jsonBody []byte
	jsonBody, exists := conf.CacheApi.Get(uri)
	if !exists {
		res, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		jsonBody, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		conf.CacheApi.Add(uri, jsonBody)
	}
	return jsonBody, nil
}

func commandMap(conf *config, args []string) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}
	var uri string
	if conf.MapNextUrl != "" {
		uri = conf.MapNextUrl
	} else {
		uri = "https://pokeapi.co/api/v2/location-area" 
	}
	jsonBody, err := cacheAwareGet(uri, conf)
	if err != nil {
		return err
	}

	var locations mapResponse
	if err = json.Unmarshal(jsonBody, &locations); err != nil {
		fmt.Printf("Error unmarshaling json: %s\n", string(jsonBody))
		return err
	}
	conf.MapNextUrl = locations.Next
	conf.MapPreviousUrl = locations.Previous
	for _, l := range locations.Results {
		fmt.Println(l.Name)
	}

	return nil
}

func commandMapBack(conf *config, args []string) error {
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
	jsonBody, err := cacheAwareGet(uri, conf)
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

func commandExplore(conf *config, args []string) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}
	if len(args) < 1 {
		return fmt.Errorf("Must provide area name to explore")
	}
	uri := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])
	jsonBody, err := cacheAwareGet(uri, conf)
	if err != nil {
		return err
	}

	var encounters encounterResponse
	if err = json.Unmarshal(jsonBody, &encounters); err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range encounters.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

func main() {
	prompt := "Pokedex > "
	interval, intervalOk := time.ParseDuration("10s")
	if intervalOk != nil {
		fmt.Printf("invalid interval: %s\n", intervalOk)
		os.Exit(1)
	}
	cache := internal.NewCache(interval)
	scanner := bufio.NewScanner(os.Stdin)
	conf := config{
		CacheApi: cache,
		Pokedex: make(map[string]internal.PokemonResponse),
	}
	commands = map[string]cliCommand {
		"catch": {
			name:		"catch",
			description:	"Yeet a pokeball at a Pokemon",
			callback:	commandCatch,
		},
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
		"explore": {
			name: 		"explore",
			description:	"Explore an area",
			callback:	commandExplore,
		},
	}
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		var args []string
		var ok bool = true
		cmd := commands["help"]
		if len(cleanedInput) != 0 {
			cmd, ok = commands[cleanedInput[0]]
		}	
		if len(cleanedInput) > 1 {
			args = cleanedInput[1:]
		}
		if ok {
			err := cmd.callback(&conf, args)
			if err != nil {
				fmt.Printf("%s error: %s\n", cmd.name, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}

}
