package main

import (
	"errors"
	"sync"
)

// ErrNotFound is returned when a short code doesn't exist.
var ErrNotFound = errors.New("not found")

// Store is anything that can save and look up links.
// In Stage 7 you'll write a DynamoDB version of this interface,
// and the rest of the app won't need to change.
type Store interface {
	Save(code, url string) error
	Get(code string) (string, error)
	Exists(code string) bool
}

// MemoryStore keeps links in a map. Everything is lost when the app restarts.
type MemoryStore struct {
	mu    sync.RWMutex // protects links: many requests can run at the same time
	links map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{links: make(map[string]string)}
}

func (s *MemoryStore) Save(code, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links[code] = url
	return nil
}

func (s *MemoryStore) Get(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.links[code]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}

func (s *MemoryStore) Exists(code string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.links[code]
	return ok
}
