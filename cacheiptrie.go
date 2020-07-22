package fwcache

import (
	"sync"
)

type IpPrefixTrieOperator interface {
	// return the IPRangeList-length in memory
	CacheLen() int

	// create a new IPRangeList
	Create(name string)

	// delete the IPRangeList
	Delete(name string)

	// create a new IPRangeList with IPs
	InsertIPs(name string, uid string, IPs []string)

	// search if ip exist in the IPRangeList
	IsIPExist(name string, uid string, IP string) bool

	// check if the IPRangeList with given name is exist.
	IsExist(name string) bool
}

type ipPrefixTrie struct {
	mu sync.RWMutex
	m  map[string]prefixTrieMap
}

func NewCacheIPPrefixTrie() IpPrefixTrieOperator {
	return &ipPrefixTrie{
		mu: sync.RWMutex{},
		m:  make(map[string]prefixTrieMap),
	}
}

func (ipPT *ipPrefixTrie) CacheLen() int{
	return len(ipPT.m)
}

func (ipPT *ipPrefixTrie) Create(name string){
	ipPT.mu.Lock()
	defer ipPT.mu.Unlock()

	ipPT.m[name] = newPrefixTrieMap()
}

func (ipPT *ipPrefixTrie) Delete(name string){
	ipPT.mu.Lock()
	delete(ipPT.m, name)
	ipPT.mu.Unlock()
}

func (ipPT *ipPrefixTrie) InsertIPs(name string, uid string, IPs []string){
	ipPT.mu.Lock()
	defer ipPT.mu.Unlock()

	ipPT.m[name].CreateSource(uid)
	ipPT.m[name].InsertIPSlice(uid, IPs)
}

func (ipPT *ipPrefixTrie) IsIPExist(name string, uid string, IP string) bool{
	ipPT.mu.RLock()
	defer ipPT.mu.RUnlock()

	if !ipPT.IsExist(name){
		return false
	}
	return ipPT.m[name].FindIP(uid, IP)
}

func (ipPT *ipPrefixTrie) IsExist(name string) bool{
	ipPT.mu.RLock()
	defer ipPT.mu.RUnlock()

	_, ok := ipPT.m[name]
	return ok
}


