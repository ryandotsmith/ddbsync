package ddbsync

import (
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	//"github.com/ryandotsmith/ddbsync/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

const VALID_MUTEX_NAME string = "mut-test"
const VALID_MUTEX_TTL int64 = 4

type MutexSuite struct {
	suite.Suite
	//mock *mocks.DBer
}

func TestMutexSuite(t *testing.T) {
	suite.Run(t, new(MutexSuite))
}

func (s *MutexSuite) SetupSuite() {

}

func (s *MutexSuite) SetupTest() {

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
	/*
		underTest := Mutex{
			Name: VALID_MUTEX_NAME,
			TTL:  VALID_MUTEX_TTL,
		}
		underTest.Lock()
	*/
}

func (s *MutexSuite) TestLockUnlock() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}

	underTest.Lock()
	// It should take us 4 seconds to acquire this lock.
	underTest.Lock()
	underTest.Unlock()
}

func (s *MutexSuite) TestUnlock() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}

	underTest.Lock()
	underTest.Unlock()
}

func (s *MutexSuite) TestPruneExpired() {
	underTest := Mutex{
		Name: VALID_MUTEX_NAME,
		TTL:  VALID_MUTEX_TTL,
	}
	underTest.PruneExpired()
}
