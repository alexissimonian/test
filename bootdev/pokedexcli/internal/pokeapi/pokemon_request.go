package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemon(name string) (PokemonResponse, error) {
	fullURL := baseURL + "/pokemon/" + name

	// cache read
	if data, ok := c.cache.Get(fullURL); ok {
		pokemon := PokemonResponse{}
        err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("Error parsing the data into a pokemon object: %v\n", err)
		}
        return pokemon, nil
	}
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonResponse{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return PokemonResponse{}, err
	}
	defer response.Body.Close()

    if response.StatusCode > 399 {
        return PokemonResponse{}, fmt.Errorf("Bad status code %v", response.StatusCode)
    }
	pokemon := PokemonResponse{}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return PokemonResponse{}, fmt.Errorf("Error reading pokemon data from api : %v\n", err)
	}

	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return PokemonResponse{}, fmt.Errorf("Error parsing the data into a pokemon object: %v\n", err)
	}

	c.cache.Add(fullURL, data)
	return pokemon, nil
}
