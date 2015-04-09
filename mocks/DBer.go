package mocks

import "github.com/stretchr/testify/mock"

import "github.com/zencoder/ddbsync/models"

type DBer struct {
	mock.Mock
}

func (m *DBer) Put(_a0 string, _a1 int64) error {
	ret := m.Called(_a0, _a1)

	r0 := ret.Error(0)

	return r0
}
func (m *DBer) Get(_a0 string) (*models.Item, error) {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(*models.Item)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *DBer) Delete(_a0 string) error {
	ret := m.Called(_a0)

	r0 := ret.Error(0)

	return r0
}
