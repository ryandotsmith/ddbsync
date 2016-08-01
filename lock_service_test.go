package ddbsync

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const LOCK_SERVICE_VALID_MUTEX_NAME string = "mut-test"
const LOCK_SERVICE_VALID_MUTEX_TTL int64 = 4

type LockServiceSuite struct {
	suite.Suite
}

func TestLockServiceSuite(t *testing.T) {
	suite.Run(t, new(LockServiceSuite))
}

func (s *LockServiceSuite) TestNewLock() {
	ls := &LockService{}
	m := ls.NewLock(LOCK_SERVICE_VALID_MUTEX_NAME, LOCK_SERVICE_VALID_MUTEX_TTL)

	assert.NotNil(s.T(), ls)
	assert.NotNil(s.T(), m)
	assert.IsType(s.T(), &LockService{}, ls)
	assert.IsType(s.T(), &Mutex{}, m)
	assert.Equal(s.T(), &Mutex{Name: LOCK_SERVICE_VALID_MUTEX_NAME, TTL: LOCK_SERVICE_VALID_MUTEX_TTL}, m)
}
