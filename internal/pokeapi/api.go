package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/etuhoha/pokedexcli/internal/pokecache"
)

type location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

var cache = pokecache.NewCache(5 * time.Second)

func Map(url string) (next string, prev string, err error) {
	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		cache.Add(url, data)
	}

	locResp := locationResponse{}
	err = json.Unmarshal(data, &locResp)
	if err != nil {
		return "", "", err
	}

	for _, l := range locResp.Results {
		fmt.Printf("%v\n", l.Name)
	}

	return locResp.Next, locResp.Previous, nil
}
