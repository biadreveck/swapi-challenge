package planet_test

import (
	"b2w/swapi-challenge/domain"
	"b2w/swapi-challenge/domain/entity/planet"
	"b2w/swapi-challenge/domain/entity/planet/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestManagerInsert(t *testing.T) {
	dbRepo := &mocks.DbRepository{}
	swapiRepo := &mocks.SwapiRepository{}

	manager := planet.NewManager(dbRepo, swapiRepo)

	pSuccess := &planet.Planet{Name: "Success"}
	pInvalid := &planet.Planet{}
	pSwapiError := &planet.Planet{Name: "Swapi Error"}
	pGetFound := &planet.Planet{Name: "Get Found"}
	pInsertError := &planet.Planet{Name: "Insert Error"}

	swapiRepo.
		On("GetPlanetApparitions", pSuccess.Name).
		Return(int32(1), nil)

	swapiRepo.
		On("GetPlanetApparitions", pGetFound.Name).
		Return(int32(1), nil)

	swapiRepo.
		On("GetPlanetApparitions", pInsertError.Name).
		Return(int32(1), nil)

	swapiRepo.
		On("GetPlanetApparitions", pSwapiError.Name).
		Return(int32(0), errors.New("swapi error"))

	dbRepo.
		On("GetByName", pSuccess.Name).
		Return(planet.Planet{}, domain.ErrNotFound)

	dbRepo.
		On("GetByName", pInsertError.Name).
		Return(planet.Planet{}, domain.ErrNotFound)

	dbRepo.
		On("GetByName", pGetFound.Name).
		Return(planet.Planet{ID: primitive.NewObjectID()}, nil)

	dbRepo.
		On("Insert", pSuccess).
		Return(nil)

	dbRepo.
		On("Insert", pInsertError).
		Return(errors.New("insert error"))

	// Testing insertion success
	err := manager.Insert(pSuccess)
	assert.Nil(t, err)
	assert.NotEqual(t, pSuccess.ID, primitive.NilObjectID)
	assert.Equal(t, pSuccess.Apparitions, int32(1))

	// Testing invalid planet
	err = manager.Insert(pInvalid)
	assert.NotNil(t, err)
	assert.Equal(t, err, domain.ErrBadParamInput)
	assert.Equal(t, pInvalid.ID, primitive.NilObjectID)

	// Testing swapi error
	err = manager.Insert(pSwapiError)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "swapi error")
	assert.Equal(t, pSwapiError.ID, primitive.NilObjectID)

	// Testing planet found
	err = manager.Insert(pGetFound)
	assert.NotNil(t, err)
	assert.Equal(t, err, domain.ErrConflict)
	assert.Equal(t, pGetFound.ID, primitive.NilObjectID)

	// Testing insertion error
	err = manager.Insert(pInsertError)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "insert error")
}

func TestManagerFindAll(t *testing.T) {
	dbRepo := &mocks.DbRepository{}
	dbRepoErr := &mocks.DbRepository{}

	manager := planet.NewManager(dbRepo, nil)
	managerErr := planet.NewManager(dbRepoErr, nil)

	pOne := planet.Planet{Name: "One"}
	pTwo := planet.Planet{Name: "Two"}
	pThree := planet.Planet{Name: "Two"}

	pList := []planet.Planet{pOne, pTwo, pThree}

	dbRepo.
		On("FindAll").
		Return(pList, nil)

	dbRepoErr.
		On("FindAll").
		Return(nil, errors.New("find all error"))

	// Testing find all success
	result, err := manager.FindAll()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result), 3)

	// Testing find all error
	result, err = managerErr.FindAll()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "find all error")
}

func TestManagerGetById(t *testing.T) {
	dbRepo := &mocks.DbRepository{}

	manager := planet.NewManager(dbRepo, nil)

	pID := primitive.NewObjectID()
	pIDErr := primitive.NewObjectID()

	p := planet.Planet{ID: pID}
	pErr := planet.Planet{}

	dbRepo.
		On("GetById", pID).
		Return(p, nil)

	dbRepo.
		On("GetById", pIDErr).
		Return(pErr, errors.New("get error"))

	// Testing get success
	result, err := manager.GetById(pID)
	assert.Nil(t, err)
	assert.Equal(t, result.ID, pID)

	// Testing get error
	result, err = manager.GetById(pIDErr)
	assert.NotNil(t, err)
	assert.Equal(t, result.ID, primitive.NilObjectID)
	assert.Equal(t, err.Error(), "get error")
}

func TestManagerDelete(t *testing.T) {
	dbRepo := &mocks.DbRepository{}

	manager := planet.NewManager(dbRepo, nil)

	pID := primitive.NewObjectID()
	pIDNotFound := primitive.NewObjectID()
	pIDErr := primitive.NewObjectID()

	dbRepo.
		On("GetById", pID).
		Return(planet.Planet{}, nil)

	dbRepo.
		On("GetById", pIDErr).
		Return(planet.Planet{}, nil)

	dbRepo.
		On("GetById", pIDNotFound).
		Return(planet.Planet{}, domain.ErrNotFound)

	dbRepo.
		On("Delete", pID).
		Return(nil)

	dbRepo.
		On("Delete", pIDErr).
		Return(errors.New("delete error"))

	// Testing delete success
	err := manager.Delete(pID)
	assert.Nil(t, err)

	// Testing planet not found
	err = manager.Delete(pIDNotFound)
	assert.NotNil(t, err)
	assert.Equal(t, err, domain.ErrNotFound)

	// Testing delete error
	err = manager.Delete(pIDErr)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "delete error")
}
