package ddbsync

import (
	"errors"
	"testing"

	"github.com/awslabs/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zencoder/ddbsync/mocks"
	"github.com/zencoder/ddbsync/models"
)

const VALID_MUTEX_NAME string = "mut-test"
const VALID_MUTEX_TTL int64 = 4
const VALID_MUTEX_CREATED int64 = 1424385592

type MutexSuite struct {
	suite.Suite
	mock *mocks.DBer
}

func TestMutexSuite(t *testing.T) {
	suite.Run(t, new(MutexSuite))
}

func (s *MutexSuite) SetupSuite() {

}

func (s *MutexSuite) SetupTest() {
	s.mock = new(mocks.DBer)
	db = s.mock
}

func (s *MutexSuite) TearDownSuite() {
	db = &database{
		client: dynamodb.New(nil),
	}
}

func (s *MutexSuite) TestNew() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}
	assert.Equal(s.T(), VALID_MUTEX_NAME, underTest.Name)
	assert.Equal(s.T(), VALID_MUTEX_TTL, underTest.TTL)
}

func (s *MutexSuite) TestLock() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}

	s.mock.On("Put", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(nil)
	s.mock.On("Get", mock.AnythingOfType("string")).Return(&models.Item{Name: VALID_MUTEX_NAME, Created: VALID_MUTEX_CREATED}, nil)
	s.mock.On("Delete", mock.AnythingOfType("string")).Return(nil)

	underTest.Lock()
}

func (s *MutexSuite) TestUnlock() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}

	s.mock.On("Delete", mock.AnythingOfType("string")).Return(nil)

	underTest.Unlock()
}

func (s *MutexSuite) TestPruneExpired() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}

	s.mock.On("Get", mock.AnythingOfType("string")).Return(&models.Item{Name: VALID_MUTEX_NAME, Created: VALID_MUTEX_CREATED}, nil)
	s.mock.On("Delete", mock.AnythingOfType("string")).Return(nil)

	underTest.PruneExpired()
}

func (s *MutexSuite) TestPruneExpiredError() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}

	s.mock.On("Get", mock.AnythingOfType("string")).Return((*models.Item)(nil), errors.New("Get Error"))

	underTest.PruneExpired()
}
