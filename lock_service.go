package ddbsync

import (
	"sync"
)

type LockServicer interface {
	NewLock(string, int64) sync.Locker
}

type LockService struct{}

var _ LockServicer = (*LockService)(nil) // Forces compile time checking of the interface

func (l *LockService) NewLock(name string, ttl int64) sync.Locker {
	return &Mutex{
		Name: name,
		TTL:  ttl,
	}
}
