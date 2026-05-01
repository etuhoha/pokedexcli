package pokeapi

import (
	"encoding/json"
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

type MapResponse struct {
	Areas []string
	Next  string
	Prev  string
}

func Map(url string) (MapResponse, error) {
	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return MapResponse{}, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return MapResponse{}, err
		}
		cache.Add(url, data)
	}

	locResp := locationResponse{}
	err := json.Unmarshal(data, &locResp)
	if err != nil {
		return MapResponse{}, err
	}

	mapResp := MapResponse{}
	for _, l := range locResp.Results {
		mapResp.Areas = append(mapResp.Areas, l.Name)
	}
	mapResp.Next = locResp.Next
	mapResp.Prev = locResp.Previous
	return mapResp, nil
}

type exploreResponse struct {
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
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
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

func ExploreArea(area string) ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area

	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, data)
	}

	expResp := exploreResponse{}
	err := json.Unmarshal(data, &expResp)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, pe := range expResp.PokemonEncounters {
		result = append(result, pe.Pokemon.Name)
	}

	return result, nil

}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type PokemonType struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"type"`
}

type PokemonData struct {
	Name    string
	BaseExp int
	Height  int
	Weight  int
	Stats   []PokemonStat
	Types   []PokemonType
}

func PokemonStats(name string) (PokemonData, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return PokemonData{}, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return PokemonData{}, err
		}
		cache.Add(url, data)
	}

	rawJson := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &rawJson)
	if err != nil {
		return PokemonData{}, err
	}

	var baseExp float64
	err = json.Unmarshal(rawJson["base_experience"], &baseExp)
	if err != nil {
		return PokemonData{}, err
	}

	var height float64
	err = json.Unmarshal(rawJson["height"], &height)
	if err != nil {
		return PokemonData{}, err
	}

	var weight float64
	err = json.Unmarshal(rawJson["weight"], &weight)
	if err != nil {
		return PokemonData{}, err
	}

	var stats []PokemonStat
	err = json.Unmarshal(rawJson["stats"], &stats)
	if err != nil {
		return PokemonData{}, err
	}

	var types []PokemonType
	err = json.Unmarshal(rawJson["types"], &types)
	if err != nil {
		return PokemonData{}, err
	}

	result := PokemonData{Name: name}
	result.BaseExp = int(baseExp)
	result.Height = int(height)
	result.Weight = int(weight)
	result.Stats = stats
	result.Types = types
	return result, nil
}
