package handler_test

import (
	"b2w/swapi-challenge/api"
	"b2w/swapi-challenge/domain"
	"b2w/swapi-challenge/domain/entity/planet"
	"b2w/swapi-challenge/domain/entity/planet/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type responseBody struct {
	Data interface{} `json:"data"`
	Err  string      `json:"error"`
}

func idMatchsParam(id string) interface{} {
	return mock.MatchedBy(func(reqId primitive.ObjectID) bool {
		return reqId.Hex() == id
	})
}

func planetMatchsName(name string) interface{} {
	return mock.MatchedBy(func(p *planet.Planet) bool {
		return p.Name == name
	})
}

func planetMatchsClimate(climate string) interface{} {
	return mock.MatchedBy(func(p *planet.Planet) bool {
		return p.Climate == climate
	})
}

func TestCreatePlanet(t *testing.T) {
	manager := &mocks.Manager{}

	router := api.SetupRouter(manager)
	ts := httptest.NewServer(router)
	defer ts.Close()

	baseUrl := fmt.Sprintf("%s/v1/planets", ts.URL)

	var baSuccess = []byte(`{"name":"Success"}`)
	var baInvalidJson = []byte(`{name:Invalid}`)
	var baInvalidPlanet = []byte(`{"climate":"temperate"}`)
	var baConflict = []byte(`{"name":"Conflict"}`)
	var baError = []byte(`{"name":"Error"}`)

	manager.
		On("Insert", planetMatchsName("Success")).
		Return(func(p *planet.Planet) error {
			p.ID = primitive.NewObjectID()
			return nil
		})

	manager.
		On("Insert", planetMatchsClimate("temperate")).
		Return(domain.ErrBadParamInput)

	manager.
		On("Insert", planetMatchsName("Conflict")).
		Return(domain.ErrConflict)

	manager.
		On("Insert", planetMatchsName("Error")).
		Return(errors.New("create error"))

	// Testing create success
	resp, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(baSuccess))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var body responseBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.Nil(t, err)

	bodyData := body.Data.(map[string]interface{})
	assert.NotNil(t, bodyData["id"])
	assert.Equal(t, "Success", bodyData["name"])

	resp.Body.Close()

	// Testing create invalid json
	resp, err = http.Post(baseUrl, "application/json", bytes.NewBuffer(baInvalidJson))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing create invalid planet
	resp, err = http.Post(baseUrl, "application/json", bytes.NewBuffer(baInvalidPlanet))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing create conflict
	resp, err = http.Post(baseUrl, "application/json", bytes.NewBuffer(baConflict))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	resp.Body.Close()

	// Testing create error
	resp, err = http.Post(baseUrl, "application/json", bytes.NewBuffer(baError))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()
}

func TestGetPlanet(t *testing.T) {
	manager := &mocks.Manager{}

	router := api.SetupRouter(manager)
	ts := httptest.NewServer(router)
	defer ts.Close()

	baseUrl := fmt.Sprintf("%s/v1/planets", ts.URL)

	pID := primitive.NewObjectID()
	pIDInvalid := "Invalid"
	pIDNotFound := primitive.NewObjectID()
	pIDErr := primitive.NewObjectID()

	p := planet.Planet{ID: pID}
	pErr := planet.Planet{}

	manager.
		On("GetById", idMatchsParam(pID.Hex())).
		Return(p, nil)

	manager.
		On("GetById", idMatchsParam(pIDNotFound.Hex())).
		Return(pErr, domain.ErrNotFound)

	manager.
		On("GetById", idMatchsParam(pIDErr.Hex())).
		Return(pErr, errors.New("get error"))

	// Testing get success
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseUrl, pID.Hex()))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body responseBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.Nil(t, err)

	bodyData := body.Data.(map[string]interface{})
	assert.Equal(t, pID.Hex(), bodyData["id"])

	resp.Body.Close()

	// Testing get invalid id
	resp, err = http.Get(fmt.Sprintf("%s/%s", baseUrl, pIDInvalid))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing get not found
	resp, err = http.Get(fmt.Sprintf("%s/%s", baseUrl, pIDNotFound.Hex()))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// Testing get error
	resp, err = http.Get(fmt.Sprintf("%s/%s", baseUrl, pIDErr.Hex()))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()
}

func TestGetPlanets(t *testing.T) {
	manager := &mocks.Manager{}

	router := api.SetupRouter(manager)
	ts := httptest.NewServer(router)
	defer ts.Close()

	baseUrl := fmt.Sprintf("%s/v1/planets", ts.URL)

	pOne := planet.Planet{Name: "One"}
	pTwo := planet.Planet{Name: "Two"}
	pThree := planet.Planet{Name: "Two"}

	pList := []planet.Planet{pOne, pTwo, pThree}

	manager.
		On("FindAll").
		Return(pList, nil)

	// Testing get success
	resp, err := http.Get(baseUrl)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body responseBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.Nil(t, err)

	bodyData := body.Data.([]interface{})
	assert.Equal(t, 3, len(bodyData))

	resp.Body.Close()
}

func TestGetPlanetsErr(t *testing.T) {
	manager := &mocks.Manager{}

	router := api.SetupRouter(manager)
	ts := httptest.NewServer(router)
	defer ts.Close()

	baseUrl := fmt.Sprintf("%s/v1/planets", ts.URL)

	manager.
		On("FindAll").
		Return(nil, errors.New("find all error"))

	// Testing get error
	resp, err := http.Get(baseUrl)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()
}

func TestGetPlanetsWithName(t *testing.T) {
	manager := &mocks.Manager{}

	router := api.SetupRouter(manager)
	ts := httptest.NewServer(router)
	defer ts.Close()

	baseUrl := fmt.Sprintf("%s/v1/planets", ts.URL)

	pName := "Success"
	pNameNotFound := "NotFound"
	pNameErr := "Error"

	pID := primitive.NewObjectID()
	p := planet.Planet{ID: pID, Name: pName}
	pErr := planet.Planet{}

	manager.
		On("GetByName", pName).
		Return(p, nil)

	manager.
		On("GetByName", pNameNotFound).
		Return(pErr, domain.ErrNotFound)

	manager.
		On("GetByName", pNameErr).
		Return(pErr, errors.New("get error"))

	// Testing get success
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", baseUrl, pName))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body responseBody
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.Nil(t, err)

	bodyData := body.Data.(map[string]interface{})

	assert.Equal(t, pName, bodyData["name"])

	resp.Body.Close()

	// Testing get not found
	resp, err = http.Get(fmt.Sprintf("%s?name=%s", baseUrl, pNameNotFound))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// Testing get error
	resp, err = http.Get(fmt.Sprintf("%s?name=%s", baseUrl, pNameErr))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()
}

func TestDelete(t *testing.T) {
	manager := &mocks.Manager{}

	router := api.SetupRouter(manager)

	pID := primitive.NewObjectID()
	pIDInvalid := "Invalid"
	pIDNotFound := primitive.NewObjectID()
	pIDErr := primitive.NewObjectID()

	manager.
		On("Delete", idMatchsParam(pID.Hex())).
		Return(nil)

	manager.
		On("Delete", idMatchsParam(pIDNotFound.Hex())).
		Return(domain.ErrNotFound)

	manager.
		On("Delete", idMatchsParam(pIDErr.Hex())).
		Return(errors.New("delete error"))

	ts := httptest.NewServer(router)
	defer ts.Close()

	baseUrl := fmt.Sprintf("%s/v1/planets", ts.URL)

	client := &http.Client{}

	// Testing get success
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", baseUrl, pID.Hex()), nil)
	assert.Nil(t, err)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Testing get invalid id
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%s", baseUrl, pIDInvalid), nil)
	assert.Nil(t, err)
	resp, err = client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Testing get not found
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%s", baseUrl, pIDNotFound.Hex()), nil)
	assert.Nil(t, err)
	resp, err = client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// Testing get error
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/%s", baseUrl, pIDErr.Hex()), nil)
	assert.Nil(t, err)
	resp, err = client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()
}
