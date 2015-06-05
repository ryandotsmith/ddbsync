package ddbsync

import (
	"sync"
)

type LockServicer interface {
	NewLock(string, int64) sync.Locker
}

type LockService struct {
	db DBer
}

var _ LockServicer = (*LockService)(nil) // Forces compile time checking of the interface

func NewLockService(tableName string, region string, endpoint string, disableSSL bool) *LockService {
	return &LockService{
		db: NewDatabase(tableName, region, endpoint, disableSSL),
	}
}

// Create a new Lock/Mutex with a particular key and timeout
func (l *LockService) NewLock(name string, ttl int64) sync.Locker {
	return NewMutex(name, ttl, l.db)
}
