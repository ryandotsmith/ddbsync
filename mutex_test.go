package ddbsync

import (
	"testing"
)

func TestLock(t *testing.T) {
	m := Mutex{"mut-test", 10}
	m.Lock()
}

func TestUnlock(t *testing.T) {
	m := Mutex{"mut-test", 10}
	m.Lock()
	m.Unlock()
}
