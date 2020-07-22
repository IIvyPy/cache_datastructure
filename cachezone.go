package fwcache

import "sync"

// ZoneTreesOperator is a type for zone/qname match, which is a map inner zTree-type.
type ZoneTreesOperator interface {

	// return the zoneTrees-length in memory
	CacheLen() int

	// create a new zone-tree
	Create(name string)

	// delete the new zone-tree
	Delete(name string)

	// create a new zone-tree with zones
	CreateWithZones(name string, zones []string)

	// update an exist zone-tree
	UpdateZones(name string, zones []string) bool

	// insert slice of zones to an exist zone-tree
	InsertZones(name string, zones []string) bool

	// search if domain exist in the zone-tree
	IsDomainExist(name string, domain string) bool

	// check if the zone-tree with given name is exist.
	IsExist(name string) bool
}

type cacheZoneTrees struct {
	mu sync.RWMutex
	m  map[string]*zoneTree
}

func NewZoneTree() ZoneTreesOperator {
	return &cacheZoneTrees{
		mu: sync.RWMutex{},
		m:  make(map[string]*zoneTree),
	}
}

func (s *cacheZoneTrees) Create(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[name] = newZTree()
}

func (s *cacheZoneTrees) Delete(name string) {
	s.mu.Lock()
	delete(s.m, name)
	s.mu.Unlock()
	return
}

func (s *cacheZoneTrees) CreateWithZones(name string, zones []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[name] = newZTree()
	s.m[name].InsertWithSlice(zones)
}

func (s *cacheZoneTrees) UpdateZones(name string, zones []string) bool {
	if !s.IsExist(name) {
		return false
	}

	newTree := newZTree()
	newTree.InsertWithSlice(zones)

	s.mu.Lock()
	s.m[name] = newTree
	s.mu.Unlock()

	return true
}

func (s *cacheZoneTrees) InsertZones(name string, zones []string) bool {
	if !s.IsExist(name) {
		return false
	}

	s.mu.Lock()
	s.m[name].InsertWithSlice(zones)
	s.mu.Unlock()

	return true
}

func (s *cacheZoneTrees) IsDomainExist(name string, domain string) bool {
	if !s.IsExist(name) {
		return false
	}

	return s.m[name].Search(domain)
}

func (s *cacheZoneTrees) CacheLen() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.m)
}

func (s *cacheZoneTrees) IsExist(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.m[name]
	return ok
}
