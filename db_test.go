package ddbsync

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zencoder/ddbsync/mocks"
	"github.com/zencoder/ddbsync/models"
)

const DB_VALID_NAME string = "db-name"
const DB_VALID_CREATED int64 = 1424385592
const DB_VALID_CREATED_STRING string = "1424385592"

type DBSuite struct {
	suite.Suite
	mock *mocks.AWSDynamoer
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func (s *DBSuite) SetupSuite() {

}

func (s *DBSuite) SetupTest() {
	s.mock = new(mocks.AWSDynamoer)
	db.(*database).client = s.mock
	os.Setenv("DDBSYNC_LOCKS_TABLE_NAME", "")
}

func (s *DBSuite) TearDownTest() {
	db.(*database).client = dynamodb.New(nil)
}

func (s *DBSuite) TestPut() {
	s.mock.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := db.Put(DB_VALID_NAME, DB_VALID_CREATED)

	assert.Nil(s.T(), err)
}

func (s *DBSuite) TestPutError() {
	s.mock.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return((*dynamodb.PutItemOutput)(nil), errors.New("PutItem Error"))

	err := db.Put(DB_VALID_NAME, DB_VALID_CREATED)

	assert.NotNil(s.T(), err)
}

func (s *DBSuite) TestGet() {
	qo := &dynamodb.QueryOutput{
		Count: aws.Long(1),
		Items: []*map[string]*dynamodb.AttributeValue{
			&map[string]*dynamodb.AttributeValue{
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

	i, err := db.Get(DB_VALID_NAME)

	assert.NotNil(s.T(), i)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), &models.Item{Name: DB_VALID_NAME, Created: DB_VALID_CREATED}, i)
}

func (s *DBSuite) TestGetErrorNoQueryOutput() {
	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return((*dynamodb.QueryOutput)(nil), errors.New("Query Error"))

	i, err := db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
}

func (s *DBSuite) TestGetErrorNilCount() {
	qo := &dynamodb.QueryOutput{
		Count: nil,
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Count not returned"), err)
}

func (s *DBSuite) TestGetErrorZeroCount() {
	qo := &dynamodb.QueryOutput{
		Count: aws.Long(0),
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New(fmt.Sprintf("No item for Name, %s", DB_VALID_NAME)), err)
}

func (s *DBSuite) TestGetErrorCountTooHigh() {
	qo := &dynamodb.QueryOutput{
		Count: aws.Long(2),
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("Expected only 1 item returned from Dynamo, got 2"), err)
}

func (s *DBSuite) TestGetErrorCountSetNoItems() {
	qo := &dynamodb.QueryOutput{
		Count: aws.Long(1),
	}

	s.mock.On("Query", mock.AnythingOfType("*dynamodb.QueryInput")).Return(qo, nil)

	i, err := db.Get(DB_VALID_NAME)

	assert.Nil(s.T(), i)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errors.New("No item returned, count is invalid."), err)
}

func (s *DBSuite) TestDelete() {
	s.mock.On("DeleteItem", mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := db.Delete(DB_VALID_NAME)

	assert.Nil(s.T(), err)
}

func (s *DBSuite) TestDeleteError() {
	s.mock.On("DeleteItem", mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return((*dynamodb.DeleteItemOutput)(nil), errors.New("Delete Error"))

	err := db.Delete(DB_VALID_NAME)

	assert.NotNil(s.T(), err)
}

func (s *DBSuite) TestLocksTableNameDefault() {
	n := locksTableName()
	assert.IsType(s.T(), "this is a string", n)
	assert.Equal(s.T(), DEFAULT_LOCKS_TABLE_NAME, n)
}

func (s *DBSuite) TestLocksTableNameEnvVarSet() {
	l := "CustomLocksTable"
	os.Setenv("DDBSYNC_LOCKS_TABLE_NAME", l)
	n := locksTableName()
	assert.IsType(s.T(), "this is a string", n)
	assert.Equal(s.T(), l, n)
}
