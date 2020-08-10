package planet

import "go.mongodb.org/mongo-driver/bson/primitive"

type DbRepository interface {
	Insert(p *Planet) error
	FindAll() ([]Planet, error)
	GetById(id primitive.ObjectID) (Planet, error)
	GetByName(name string) (Planet, error)
	Delete(id primitive.ObjectID) error
}

type SwapiRepository interface {
	GetPlanetApparitions(name string) (int32, error)
}

type Manager interface {
	DbRepository
}
