package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResponse, error) {
	fullURL := ""
	if pageURL == nil {
		endpoint := "/location-area"
		fullURL = baseURL + endpoint
	} else {
		fullURL = *pageURL
	}

	// check cache
	cache, ok := c.cache.Get(fullURL)
	if ok {
		// cache hit
		fmt.Println("cache hit !")
		locationAreasResponse := LocationAreasResponse{}

		err := json.Unmarshal(cache, &locationAreasResponse)
		if err != nil {
			return LocationAreasResponse{}, nil
		}

		return locationAreasResponse, nil

	}

	fmt.Println("cache miss !")

	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return LocationAreasResponse{}, fmt.Errorf("Bad status code : %v", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	locationAreasResponse := LocationAreasResponse{}

	err = json.Unmarshal(data, &locationAreasResponse)
	if err != nil {
		return LocationAreasResponse{}, nil
	}

	c.cache.Add(fullURL, data)

	return locationAreasResponse, nil

}

func (c *Client) GetLocationArea(locationArea string) (LocationAreaResponse, error) {
	endpoint := "/location-area" + "/" + locationArea
    fullURL := baseURL + endpoint

	// check cache
	cache, ok := c.cache.Get(fullURL)
	if ok {
		// cache hit
		fmt.Println("cache hit !")
		locationAreaResponse := LocationAreaResponse{}

		err := json.Unmarshal(cache, &locationAreaResponse)
		if err != nil {
			return LocationAreaResponse{}, nil
		}

		return locationAreaResponse, nil

	}

	fmt.Println("cache miss !")

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

	c.cache.Add(fullURL, data)

	return locationAreaResponse, nil

}
