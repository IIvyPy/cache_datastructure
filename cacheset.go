package fwcache

import (
	"sync"
)

// SetOperator return an interface to operator set-type data in memory
type SetOperator interface {
	// create a new set
	Create(name string)

	// add a slice to set
	CreateWithSlice(name string, val []interface{})

	// return all the members in set
	SMembers(name string) []interface{}

	// return the number of the members in the given name set
	SCard(name string) int

	// delete set by name
	Delete(name string)

	// remove a value from the set
	SRemove(name string, val interface{}) bool

	// return if the key exists in the set given name
	SIsMember(name string, member interface{}) bool

	// update the given name set with slice
	SUpdateWithSlice(name string, val []interface{})

	// return the set-length in memory
	CacheLen() int

	// return if the set given name is in the memory
	IsExist(name string) bool

	// return a val to the set given name
	SAdd(name string, val interface{})

	// return values to the set given name
	SAddWithSlice(name string, values []interface{})
}

type cacheSet struct {
	mu sync.RWMutex
	m  map[string]Set
}

func NewSet() SetOperator {
	return &cacheSet{
		mu: sync.RWMutex{},
		m:  make(map[string]Set),
	}
}

func (s *cacheSet) Create(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[name] = NewThreadUnsafeSet()
}

func (s *cacheSet) CreateWithSlice(name string, val []interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[name] = NewThreadUnsafeSetFromSlice(val)
}

func (s *cacheSet) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, name)
}

func (s *cacheSet) SRemove(name string, val interface{}) bool {
	if !s.IsExist(name) {
		return false
	}

	s.mu.Lock()
	s.m[name].Remove(val)
	s.mu.Unlock()
	return true
}

func (s *cacheSet) SUpdateWithSlice(name string, val []interface{}) {
	s.CreateWithSlice(name, val)
}

func (s *cacheSet) SIsMember(name string, member interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cache, ok := s.m[name]; ok {
		return cache.Contains(member)
	}
	return false
}

func (s *cacheSet) IsExist(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.m[name]
	return ok
}

func (s *cacheSet) SMembers(name string) []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cache, ok := s.m[name]; ok {
		return cache.ToSlice()
	}

	return []interface{}{}
}

func (s *cacheSet) SCard(name string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.m[name]; ok {
		return s.m[name].Len()
	}

	return 0
}

func (s *cacheSet) SAdd(name string, val interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[name].Add(val)
}

func (s *cacheSet) SAddWithSlice(name string, values []interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[name].AddWithSlice(values)
}

func (s *cacheSet) CacheLen() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.m)
}
