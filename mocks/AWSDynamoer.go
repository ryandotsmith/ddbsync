package mocks

import "github.com/stretchr/testify/mock"

import "github.com/aws/aws-sdk-go/service/dynamodb"

type AWSDynamoer struct {
	mock.Mock
}

// PutItem provides a mock function with given fields: _a0
func (_m *AWSDynamoer) PutItem(_a0 *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.PutItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.PutItemInput) *dynamodb.PutItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.PutItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.PutItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: _a0
func (_m *AWSDynamoer) Query(_a0 *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.QueryOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.QueryInput) *dynamodb.QueryOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.QueryOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.QueryInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteItem provides a mock function with given fields: _a0
func (_m *AWSDynamoer) DeleteItem(_a0 *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	ret := _m.Called(_a0)

	var r0 *dynamodb.DeleteItemOutput
	if rf, ok := ret.Get(0).(func(*dynamodb.DeleteItemInput) *dynamodb.DeleteItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DeleteItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dynamodb.DeleteItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
