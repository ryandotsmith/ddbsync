// Copyright 2012 Ryan Smith. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ddbsync provides DynamoDB-backed synchronization primitives such
// as mutual exclusion locks. This package is designed to behave like pkg/sync.

package ddbsync

import (
	"log"
	"sync"
	"time"
)

// A Mutex is a mutual exclusion lock.
// Mutexes can be created as part of other structures.
type Mutex struct {
	Name string
	TTL  int64
	db   DBer
}

var _ sync.Locker = (*Mutex)(nil) // Forces compile time checking of the interface

// Mutex constructor
func NewMutex(name string, ttl int64, db DBer) *Mutex {
	return &Mutex{
		Name: name,
		TTL:  ttl,
		db:   db,
	}
}

// Lock will write an item in a DynamoDB table if the item does not exist.
// Before writing the lock, we will clear any locks that are expired.
// Calling this function will block until a lock can be acquired.
func (m *Mutex) Lock() {
	for {
		m.PruneExpired()
		err := m.db.Put(m.Name, time.Now().Unix())
		if err == nil {
			return
		}
	}
}

// Unlock will delete an item in a DynamoDB table.
func (m *Mutex) Unlock() {
	for {
		err := m.db.Delete(m.Name)
		if err == nil {
			return
		}
	}
}

// PruneExpired delete all locks that have lived past their TTL.
// This is to prevent deadlock from processes that have taken locks
// but never removed them after execution. This commonly happens when a
// processor experiences network failure.
func (m *Mutex) PruneExpired() {
	item, err := m.db.Get(m.Name)
	if err != nil {
		log.Printf("PruneExpired. error = %v", err)
		return
	}
	if item != nil {
		if item.Created < (time.Now().Unix() - m.TTL) {
			m.Unlock()
		}
	}
	return
}
