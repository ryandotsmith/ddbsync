package mocks

import "github.com/stretchr/testify/mock"

import "github.com/zencoder/ddbsync/models"

type DBer struct {
	mock.Mock
}

// Put provides a mock function with given fields: _a0, _a1
func (_m *DBer) Put(_a0 string, _a1 int64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0
func (_m *DBer) Get(_a0 string) (*models.Item, error) {
	ret := _m.Called(_a0)

	var r0 *models.Item
	if rf, ok := ret.Get(0).(func(string) *models.Item); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0
func (_m *DBer) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
