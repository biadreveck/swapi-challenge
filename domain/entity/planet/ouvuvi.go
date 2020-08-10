package planet

func NewPlanet() *Planet {
	return &Planet{
		Name:        "Yavin IV",
		Climate:     "temperate, tropical",
		Terrain:     "jungle, rainforests",
		Apparitions: 1,
	}
}
