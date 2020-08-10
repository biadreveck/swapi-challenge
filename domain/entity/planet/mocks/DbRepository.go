// Code generated by mockery v2.1.0. DO NOT EDIT.

package mocks

import (
	planet "b2w/swapi-challenge/domain/entity/planet"

	mock "github.com/stretchr/testify/mock"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// DbRepository is an autogenerated mock type for the DbRepository type
type DbRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *DbRepository) Delete(id primitive.ObjectID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields:
func (_m *DbRepository) FindAll() ([]planet.Planet, error) {
	ret := _m.Called()

	var r0 []planet.Planet
	if rf, ok := ret.Get(0).(func() []planet.Planet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]planet.Planet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *DbRepository) GetById(id primitive.ObjectID) (planet.Planet, error) {
	ret := _m.Called(id)

	var r0 planet.Planet
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) planet.Planet); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(planet.Planet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *DbRepository) GetByName(name string) (planet.Planet, error) {
	ret := _m.Called(name)

	var r0 planet.Planet
	if rf, ok := ret.Get(0).(func(string) planet.Planet); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(planet.Planet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: p
func (_m *DbRepository) Insert(p *planet.Planet) error {
	ret := _m.Called(p)

	var r0 error
	if rf, ok := ret.Get(0).(func(*planet.Planet) error); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
