package presenter

import (
	"b2w/swapi-challenge/domain/entity/planet"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddPlanetCommand struct {
	Name    string `json:"name"`
	Climate string `json:"climate"`
	Terrain string `json:"terrain"`
}

type PlanetResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Climate     string `json:"climate"`
	Terrain     string `json:"terrain"`
	Apparitions int32  `json:"apparitions"`
}

func (p AddPlanetCommand) ToModel() planet.Planet {
	return planet.Planet{
		Name:    p.Name,
		Climate: p.Climate,
		Terrain: p.Terrain,
	}
}

func NewPlanetResult(p planet.Planet) PlanetResult {
	if p.ID == primitive.NilObjectID {
		return PlanetResult{}
	}

	return PlanetResult{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Climate:     p.Climate,
		Terrain:     p.Terrain,
		Apparitions: p.Apparitions,
	}
}

func NewPlanetResultSlice(pSlice []planet.Planet) []PlanetResult {
	if pSlice == nil {
		return make([]PlanetResult, 0)
	}

	resultSlice := make([]PlanetResult, 0, len(pSlice))
	for _, p := range pSlice {
		resultSlice = append(resultSlice, NewPlanetResult(p))
	}

	return resultSlice
}
