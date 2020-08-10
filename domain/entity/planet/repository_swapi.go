package planet

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"net/url"

	"b2w/swapi-challenge/config"
)

type swapiRepo struct{}

func (swapiRepo) Endpoint() string {
	swapiConfig := config.Data.SWApi
	return fmt.Sprintf("%s/%s", swapiConfig.BaseUrl, "planets")
}

func NewSWApiRepository() *swapiRepo {
	return &swapiRepo{}
}

func (r swapiRepo) GetPlanetApparitions(name string) (int32, error) {
	escapedName := url.QueryEscape(name)
	searchUrl := fmt.Sprintf("%s/?search=%s", r.Endpoint(), escapedName)
	response, err := http.Get(searchUrl)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	data := make(map[string]interface{})
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	var results []interface{}
	var typeOk bool
	if results, typeOk = data["results"].([]interface{}); !typeOk {
		return 0, errors.New("wrong result type")
	}
	if len(results) <= 0 {
		return 0, nil
	}

	var firstResult map[string]interface{}
	if firstResult, typeOk = results[0].(map[string]interface{}); !typeOk {
		return 0, errors.New("wrong first result type")
	}

	var firstResultFilms []interface{}
	if firstResultFilms, typeOk = firstResult["films"].([]interface{}); !typeOk {
		return 0, errors.New("wrong first result films type")
	}

	apparitions := len(firstResultFilms)
	return int32(apparitions), nil
}
