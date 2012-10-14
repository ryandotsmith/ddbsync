package ddbsync

import (
	"testing"
)

func TestLockUnlock(t *testing.T) {
	m := Mutex{"mut-test", 4}
	m.Lock()
	// It should take us 4 seconds to acquire this lock.
	m.Lock()
	m.Unlock()
}

func TestUnlock(t *testing.T) {
	m := Mutex{"mut-test", 4}
	m.Lock()
	m.Unlock()
}
