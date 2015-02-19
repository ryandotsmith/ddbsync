package ddbsync

import (
	//"github.com/stretchr/testify/assert"
	//"github.com/ryandotsmith/ddbsync/mocks"
	//"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DBSuite struct {
	suite.Suite
	//mock *mocks.DBer
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func (s *DBSuite) SetupSuite() {

}

func (s *DBSuite) SetupTest() {
	//s.mock = new(mocks.DBer)
	//db = s.mock
}

func (s *DBSuite) TestPut() {

}

func (s *DBSuite) TestGet() {

}

func (s *DBSuite) TestDelete() {

}
