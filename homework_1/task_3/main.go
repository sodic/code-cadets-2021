package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type PokemonDTO struct {
	Name         string
	LocationsURL string `json:"location_area_encounters"`
}

type LocationDTO struct {
	LocationAreaDTO struct {
		Name string
	} `json:"location_area"`
}

type Pokemon struct {
	Name      string
	Locations []string
}

const pokemonApi = "https://pokeapi.co/api/v2/pokemon/"

func main() {
	pokemonId, err := readPokemonIdArg()
	if err != nil {
		log.Fatal(err)
	}

	data, err := getFullPokemonData(pokemonId)
	if err != nil {
		log.Fatal(err)
	}

	displayPokemonData(data)
}

func readPokemonIdArg() (string, error) {
	arguments := os.Args[1:]
	if len(arguments) != 1 {
		return "", errors.New("the program takes a single argument, a Pokemon's name or its ID")
	}
	return arguments[0], nil
}

func getFullPokemonData(pokemonId string) (Pokemon, error) {
	pokemonInfo, err := fetchPokemonInfo(pokemonId)
	if err != nil {
		return Pokemon{}, err
	}

	locations, err := fetchLocationInfo(pokemonInfo.LocationsURL)
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		Name:      pokemonInfo.Name,
		Locations: extractLocationNames(locations),
	}, nil
}

func displayPokemonData(pokemon Pokemon) {
	fmt.Println("Pokemon:", pokemon.Name)
	fmt.Println("Locations:", strings.Join(pokemon.Locations, ", "))
}

func fetchPokemonInfo(id string) (PokemonDTO, error) {
	var pokemonInfo PokemonDTO

	body, err := fetchAndReadData(pokemonApi + id)
	if err != nil {
		return pokemonInfo, errors.WithMessage(err, "failed to get Pokemon data")
	}

	err = json.Unmarshal(body, &pokemonInfo)
	if err != nil {
		return pokemonInfo, errors.WithMessage(err, "failed to parse the Pokemon API's response")
	}

	return pokemonInfo, nil
}

func fetchLocationInfo(locationUrl string) ([]LocationDTO, error) {
	body, err := fetchAndReadData(locationUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse the location API's response")
	}

	var locationInfo []LocationDTO
	err = json.Unmarshal(body, &locationInfo)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse the location API's response")
	}

	return locationInfo, nil
}

func extractLocationNames(locations []LocationDTO) []string {
	result := make([]string, 0, len(locations))
	for _, location := range locations {
		result = append(result, location.LocationAreaDTO.Name)
	}
	return result
}

func fetchAndReadData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.WithMessage(err, "failed request towards the Pokemon API")
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("the specified resource does not exist")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read the Pokemon API's response")
	}

	return body, nil
}
