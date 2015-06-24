package mocks

import "github.com/stretchr/testify/mock"

import "github.com/aws/aws-sdk-go/service/dynamodb"

type AWSDynamoer struct {
	mock.Mock
}

func (m *AWSDynamoer) PutItem(_a0 *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	ret := m.Called(_a0)

	var r0 *dynamodb.PutItemOutput
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dynamodb.PutItemOutput)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *AWSDynamoer) Query(_a0 *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	ret := m.Called(_a0)

	var r0 *dynamodb.QueryOutput
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dynamodb.QueryOutput)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *AWSDynamoer) DeleteItem(_a0 *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	ret := m.Called(_a0)

	var r0 *dynamodb.DeleteItemOutput
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*dynamodb.DeleteItemOutput)
	}
	r1 := ret.Error(1)

	return r0, r1
}
