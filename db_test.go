package ddbsync

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zencoder/ddbsync/mocks"
	"github.com/zencoder/ddbsync/models"
)

const (
	DB_VALID_TABLE_NAME     string = "TestLockTable"
	DB_VALID_REGION         string = "us-west-2"
	DB_VALID_NO_ENDPOINT    string = ""
	DB_VALID_DISABLE_SSL_NO bool   = false
	DB_VALID_NAME           string = "db-name"
	DB_VALID_CREATED        int64  = 1424385592
	DB_VALID_CREATED_STRING string = "1424385592"
)

type DBSuite struct {
	suite.Suite
	mock *mocks.AWSDynamoer
	db   DBer
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func (s *DBSuite) SetupSuite() {

}

func (s *DBSuite) SetupTest() {
	s.mock = new(mocks.AWSDynamoer)
	s.db = NewDatabase(DB_VALID_TABLE_NAME, DB_VALID_REGION, DB_VALID_NO_ENDPOINT, DB_VALID_DISABLE_SSL_NO)
	s.db.(*database).client = s.mock
}

func (s *DBSuite) TestPut() {
	s.mock.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := s.db.Put(DB_VALID_NAME, DB_VALID_CREATED)

	assert.Nil(s.T(), err)
}

func (s *DBSuite) TestPutError() {
	s.mock.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return((*dynamodb.PutItemOutput)(nil), errors.New("PutItem Error"))

	err := s.db.Put(DB_VALID_NAME, DB_VALID_CREATED)

	assert.NotNil(s.T(), err)
}

func (s *DBSuite) TestGet() {
	one := int64(1)
	qo := &dynamodb.QueryOutput{
		Count: &one,
		Items: []map[string]*dynamodb.AttributeValue{
			map[string]*dynamodb.AttributeValue{
				"Name": &dynamodb.AttributeValue{
					S: aws.String(DB_VALID_NAME),
				},
				"Created": &dynamodb.AttributeValue{
					N: aws.String(DB_VALID_CREATED_STRING),
				},
			},
		},
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := s.db.Get(DB_VALID_NAME)

	assert.NotNil(s.T(), i)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), &models.Item{Name: DB_VALID_NAME, Created: DB_VALID_CREATED}, i)
}

func (s *DBSuite) TestGetErrorNoQueryOutput() {
	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return((*dynamodb.QueryOutput)(nil), errors.New("Query Error"))

	i, err := s.db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
}

func (s *DBSuite) TestGetErrorNilCount() {
	qo := &dynamodb.QueryOutput{
		Count: nil,
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := s.db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Count not returned"), err)
}

func (s *DBSuite) TestGetErrorZeroCount() {
	zero := int64(0)
	qo := &dynamodb.QueryOutput{
		Count: &zero,
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := s.db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New(fmt.Sprintf("No item for Name, %s", DB_VALID_NAME)), err)
}

func (s *DBSuite) TestGetErrorCountTooHigh() {
	two := int64(2)
	qo := &dynamodb.QueryOutput{
		Count: &two,
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := s.db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Expected only 1 item returned from Dynamo, got 2"), err)
}

func (s *DBSuite) TestGetErrorCountSetNoItems() {
	one := int64(1)
	qo := &dynamodb.QueryOutput{
		Count: &one,
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := s.db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("No item returned, count is invalid."), err)
}

func (s *DBSuite) TestDelete() {
	s.mock.On("DeleteItem", mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := s.db.Delete(DB_VALID_NAME)

	assert.Nil(s.T(), err)
}

func (s *DBSuite) TestDeleteError() {
	s.mock.On("DeleteItem", mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return((*dynamodb.DeleteItemOutput)(nil), errors.New("Delete Error"))

	err := s.db.Delete(DB_VALID_NAME)

	assert.NotNil(s.T(), err)
}
