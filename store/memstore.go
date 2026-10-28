package store

import "sync"

// MemStore is an in-process KVStore backed by a map. It is the default store
// for tests and for `borlabd start` in simulation mode.
type MemStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewMemStore returns an empty MemStore.
func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string][]byte)}
}

// Get returns a copy of the value stored at key, or nil when absent.
func (s *MemStore) Get(key []byte) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[string(key)]
	if !ok {
		return nil
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out
}

// Has reports whether key is present.
func (s *MemStore) Has(key []byte) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[string(key)]
	return ok
}

// Set stores value at key.
func (s *MemStore) Set(key, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]byte, len(value))
	copy(cp, value)
	s.data[string(key)] = cp
}

// Delete removes key.
func (s *MemStore) Delete(key []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, string(key))
}
