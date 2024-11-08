package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	fullURL := ""
	if pageURL == nil {
		endpoint := "/location-area"
		fullURL = baseURL + endpoint
	} else {
        fullURL = *pageURL
    }

	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return LocationAreaResponse{}, fmt.Errorf("Bad status code : %v", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	locationAreaResponse := LocationAreaResponse{}

	err = json.Unmarshal(data, &locationAreaResponse)
	if err != nil {
		return LocationAreaResponse{}, nil
	}

	return locationAreaResponse, nil

}
