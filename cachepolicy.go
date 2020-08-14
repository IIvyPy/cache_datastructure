package fwcache

import "sync"

type PolicyOperator interface {
	Create(name string)

	Delete(name string)

	IsExist(name string) bool

	// when the policy's priority is not known, use insert func. Time Complexity -> O(n)
	InsertPolicy(name string, id uint64, policy ItemInterface)

	// when the policy's priority is known, use push back func. Time Complexity -> O(1)
	PushBackPolicy(name string, id uint64, policy ItemInterface)

	DeletePolicy(name string, id uint64)

	UpdatePolicy(name string, id uint64, policy ItemInterface)

	GetPolicyForMatch(name string) *MapList

	CacheLen() int

	Purge(name string) bool
}

type cachePolicy struct {
	mu sync.RWMutex
	m  map[string]*MapList
}

func NewPolicyOperator() PolicyOperator {
	return &cachePolicy{
		mu: sync.RWMutex{},
		m:  make(map[string]*MapList),
	}
}

func (c *cachePolicy) Create(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[name] = newMapList()
}

func (c *cachePolicy) Delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.m, name)
}

func (c *cachePolicy) IsExist(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.m[name]
	return ok
}

func (c *cachePolicy) InsertPolicy(name string, id uint64, policy ItemInterface) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[name].add(id, policy)
}

func (c *cachePolicy) PushBackPolicy(name string, id uint64, policy ItemInterface) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[name].pushBack(id, policy)
}

func (c *cachePolicy) DeletePolicy(name string, id uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[name].delete(id)
}

func (c *cachePolicy) UpdatePolicy(name string, id uint64, policy ItemInterface) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[name].update(id, policy)
}

func (c *cachePolicy) GetPolicyForMatch(name string) *MapList {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.m[name]
}

func (c *cachePolicy) CacheLen() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.m)
}

func (c *cachePolicy) Purge(name string) bool{
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsExist(name){
		return false
	}

	c.m[name].init()
	return true
}
