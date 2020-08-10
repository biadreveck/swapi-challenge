package planet

import (
	"b2w/swapi-challenge/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type manager struct {
	dbRepo    DbRepository
	swapiRepo SwapiRepository
}

func NewManager(dbR DbRepository, swapiR SwapiRepository) *manager {
	return &manager{
		dbRepo:    dbR,
		swapiRepo: swapiR,
	}
}

func (m *manager) Insert(p *Planet) error {
	if err := p.Validate(); err != nil {
		return domain.ErrBadParamInput
	}

	apparitions, err := m.swapiRepo.GetPlanetApparitions(p.Name)
	if err != nil {
		return err
	}

	existingP, _ := m.GetByName(p.Name)
	if existingP.ID != primitive.NilObjectID {
		return domain.ErrConflict
	}

	p.ID = primitive.NewObjectID()
	p.Apparitions = apparitions

	return m.dbRepo.Insert(p)
}

func (m *manager) FindAll() ([]Planet, error) {
	return m.dbRepo.FindAll()
}

func (m *manager) GetById(id primitive.ObjectID) (Planet, error) {
	return m.dbRepo.GetById(id)
}

func (m *manager) GetByName(name string) (Planet, error) {
	return m.dbRepo.GetByName(name)
}

func (m *manager) Delete(id primitive.ObjectID) error {
	_, err := m.dbRepo.GetById(id)
	if err != nil {
		return err
	}
	return m.dbRepo.Delete(id)
}
