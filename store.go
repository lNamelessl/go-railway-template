package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// memStore is an in-memory item store guarded by a mutex. It exists so the
// template runs with zero dependencies and zero configuration — swap it for
// a database-backed implementation in your own service. State is intentionally
// not persisted: a restart starts with an empty store.
type memStore struct {
	mu    sync.RWMutex
	items map[string]item
}

var store memStore

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand failure is unrecoverable
	}
	return hex.EncodeToString(b)
}

func (s *memStore) create(name string) item {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.items == nil {
		s.items = make(map[string]item)
	}
	it := item{ID: newID(), Name: name}
	s.items[it.ID] = it
	return it
}

func (s *memStore) list() []item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	return out
}

func (s *memStore) get(id string) (item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	return it, ok
}
