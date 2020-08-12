package planet_test

import (
	"b2w/swapi-challenge/domain"
	"b2w/swapi-challenge/domain/entity/planet"
	"b2w/swapi-challenge/infra/database/mocks"
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepoInsert(t *testing.T) {
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	dbRepo := planet.NewMongoRepository(dbHelper)

	pID := primitive.NewObjectID()

	pSuccess := &planet.Planet{Name: "Success"}
	pError := &planet.Planet{Name: "Error"}

	collectionHelper.
		On("InsertOne", mock.Anything, pSuccess).
		Return(&mongo.InsertOneResult{InsertedID: pID}, nil)

	collectionHelper.
		On("InsertOne", mock.Anything, pError).
		Return(nil, errors.New("insert error"))

	dbHelper.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelper)

	// Testing insertion success
	err := dbRepo.Insert(pSuccess)
	assert.Nil(t, err)
	assert.Equal(t, pID, pSuccess.ID)

	// Testing insertion error
	err = dbRepo.Insert(pError)
	assert.NotNil(t, err)
	assert.Equal(t, "insert error", err.Error())
}

func TestRepoFindAll(t *testing.T) {
	// Testing find and cursor success
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	cursorHelper := &mocks.CursorHelper{}
	dbRepo := planet.NewMongoRepository(dbHelper)

	pOne := planet.Planet{Name: "One"}
	pTwo := planet.Planet{Name: "Two"}
	pThree := planet.Planet{Name: "Two"}

	cursorHelper.
		On("Close", mock.Anything).
		Return(nil)

	cursorHelper.
		On("All", mock.Anything, mock.AnythingOfType("*[]planet.Planet")).
		Return(func(ctx context.Context, v interface{}) error {
			list := v.(*[]planet.Planet)
			*list = append(*list, pOne)
			*list = append(*list, pTwo)
			*list = append(*list, pThree)
			return nil
		})

	collectionHelper.
		On("Find", mock.Anything, bson.M{}).
		Return(cursorHelper, nil)

	dbHelper.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelper)

	result, err := dbRepo.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(result))

	// Testing cursor error
	dbHelperCursorErr := &mocks.DatabaseHelper{}
	collectionHelperCursorErr := &mocks.CollectionHelper{}
	cursorHelperErr := &mocks.CursorHelper{}
	dbRepoCursorErr := planet.NewMongoRepository(dbHelperCursorErr)

	cursorHelperErr.
		On("Close", mock.Anything).
		Return(nil)

	cursorHelperErr.
		On("All", mock.Anything, mock.AnythingOfType("*[]planet.Planet")).
		Return(errors.New("cursor error"))

	collectionHelperCursorErr.
		On("Find", mock.Anything, bson.M{}).
		Return(cursorHelperErr, nil)

	dbHelperCursorErr.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelperCursorErr)

	result, err = dbRepoCursorErr.FindAll()
	assert.NotNil(t, err)
	assert.Equal(t, "cursor error", err.Error())
	assert.Equal(t, 0, len(result))

	// Testing find error
	dbHelperFindErr := &mocks.DatabaseHelper{}
	collectionHelperFindErr := &mocks.CollectionHelper{}
	dbRepoFindErr := planet.NewMongoRepository(dbHelperFindErr)

	collectionHelperFindErr.
		On("Find", mock.Anything, bson.M{}).
		Return(nil, errors.New("find error"))

	dbHelperFindErr.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelperFindErr)

	result, err = dbRepoFindErr.FindAll()
	assert.NotNil(t, err)
	assert.Equal(t, "find error", err.Error())
	assert.Equal(t, 0, len(result))
}

func TestRepoGetById(t *testing.T) {
	// Testing find success
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	singleResultHelper := &mocks.SingleResultHelper{}
	dbRepo := planet.NewMongoRepository(dbHelper)

	pID := primitive.NewObjectID()
	pIDNotFound := primitive.NewObjectID()
	pIDOtherErr := primitive.NewObjectID()

	singleResultHelper.
		On("Decode", mock.AnythingOfType("*planet.Planet")).
		Return(func(v interface{}) error {
			pV := v.(*planet.Planet)
			pV.ID = pID
			pV.Name = "One"
			return nil
		})

	collectionHelper.
		On("FindOne", mock.Anything, bson.M{"_id": pID}).
		Return(singleResultHelper)

	dbHelper.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelper)

	result, err := dbRepo.GetById(pID)
	assert.Nil(t, err)
	assert.Equal(t, pID, result.ID)

	// Testing item not found decode error
	singleResultHelperNotFoundErr := &mocks.SingleResultHelper{}

	singleResultHelperNotFoundErr.
		On("Decode", mock.AnythingOfType("*planet.Planet")).
		Return(mongo.ErrNoDocuments)

	collectionHelper.
		On("FindOne", mock.Anything, bson.M{"_id": pIDNotFound}).
		Return(singleResultHelperNotFoundErr)

	result, err = dbRepo.GetById(pIDNotFound)
	assert.NotNil(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Equal(t, primitive.NilObjectID, result.ID)

	// Testing other decode error
	singleResultHelperOtherErr := &mocks.SingleResultHelper{}

	singleResultHelperOtherErr.
		On("Decode", mock.AnythingOfType("*planet.Planet")).
		Return(errors.New("other decode error"))

	collectionHelper.
		On("FindOne", mock.Anything, bson.M{"_id": pIDOtherErr}).
		Return(singleResultHelperOtherErr)

	result, err = dbRepo.GetById(pIDOtherErr)
	assert.NotNil(t, err)
	assert.Equal(t, "other decode error", err.Error())
	assert.Equal(t, primitive.NilObjectID, result.ID)
}

func TestRepoGetByName(t *testing.T) {
	// Testing find success
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	singleResultHelper := &mocks.SingleResultHelper{}
	dbRepo := planet.NewMongoRepository(dbHelper)

	pName := "One"
	pNameNotFound := "Not Found"
	pNameOtherErr := "Other Error"

	singleResultHelper.
		On("Decode", mock.AnythingOfType("*planet.Planet")).
		Return(func(v interface{}) error {
			pV := v.(*planet.Planet)
			pV.ID = primitive.NewObjectID()
			pV.Name = "One"
			return nil
		})

	collectionHelper.
		On("FindOne", mock.Anything, bson.M{"name": pName}).
		Return(singleResultHelper)

	dbHelper.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelper)

	result, err := dbRepo.GetByName(pName)
	assert.Nil(t, err)
	assert.Equal(t, pName, result.Name)

	// Testing item not found decode error
	singleResultHelperNotFoundErr := &mocks.SingleResultHelper{}

	singleResultHelperNotFoundErr.
		On("Decode", mock.AnythingOfType("*planet.Planet")).
		Return(mongo.ErrNoDocuments)

	collectionHelper.
		On("FindOne", mock.Anything, bson.M{"name": pNameNotFound}).
		Return(singleResultHelperNotFoundErr)

	result, err = dbRepo.GetByName(pNameNotFound)
	assert.NotNil(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Equal(t, primitive.NilObjectID, result.ID)

	// Testing other decode error
	singleResultHelperOtherErr := &mocks.SingleResultHelper{}

	singleResultHelperOtherErr.
		On("Decode", mock.AnythingOfType("*planet.Planet")).
		Return(errors.New("other decode error"))

	collectionHelper.
		On("FindOne", mock.Anything, bson.M{"name": pNameOtherErr}).
		Return(singleResultHelperOtherErr)

	result, err = dbRepo.GetByName(pNameOtherErr)
	assert.NotNil(t, err)
	assert.Equal(t, "other decode error", err.Error())
	assert.Equal(t, primitive.NilObjectID, result.ID)
}

func TestRepoDelete(t *testing.T) {
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	dbRepo := planet.NewMongoRepository(dbHelper)

	pID := primitive.NewObjectID()
	pIDErr := primitive.NewObjectID()

	collectionHelper.
		On("DeleteOne", mock.Anything, bson.M{"_id": pID}).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

	collectionHelper.
		On("DeleteOne", mock.Anything, bson.M{"_id": pIDErr}).
		Return(nil, errors.New("delete error"))

	dbHelper.
		On("Collection", dbRepo.CollectionName()).
		Return(collectionHelper)

	// Testing deletion success
	err := dbRepo.Delete(pID)
	assert.Nil(t, err)

	// Testing deletion error
	err = dbRepo.Delete(pIDErr)
	assert.NotNil(t, err)
	assert.Equal(t, "delete error", err.Error())
}
