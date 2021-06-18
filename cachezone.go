package fwcache

import (
	"github.com/infobloxopen/go-trees/domaintree"
	"github.com/orcaman/concurrent-map"
	"strings"
)

var none = struct{}{}

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

	// search if domain exist in the zone-tree
	IsDomainExist(name, domain string) (interface{}, bool)

	// IsExist if the name of zone-tree exist
	IsExist(name string) bool

	// insert an zone with value to an exist zone-tree
	InsertZoneWithValue(name string, value *DomainContent) bool

	// delete an zone with an exist zone-tree
	DeleteZone(name string, domain string) bool
}

type cacheZoneTrees struct{
	m cmap.ConcurrentMap
}

func NewZoneTrees() ZoneTreesOperator {
	return &cacheZoneTrees{m: cmap.New()}
}

func (zTree *cacheZoneTrees) CacheLen() int{
	return zTree.m.Count()
}

func (zTree *cacheZoneTrees) Create(name string){
	dValue := &domainWithValue{node: new(domaintree.Node)}
	zTree.m.Set(name, dValue)
}

func (zTree *cacheZoneTrees) Delete(name string){
	zTree.m.Remove(name)
}

func (zTree *cacheZoneTrees) CreateWithZones(name string, zones []string){
	Value := &domainWithValue{node: new(domaintree.Node)}
	for _, zone := range zones{
		Value.Insert(&DomainContent{Domain: trimPrefixAsterisk(zone), Value: none})
	}
	zTree.m.Set(name, Value)
}

func (zTree *cacheZoneTrees) UpdateZones(name string, zones []string) bool{
	oldValue, isExist := zTree.m.Get(name)
	if !isExist{
		return false
	}
	v, ok := oldValue.(*domainWithValue)
	if !ok{
		return false
	}
	v.Clear()

	Value := &domainWithValue{node: new(domaintree.Node)}
	for _, zone := range zones{
		Value.Insert(&DomainContent{Domain: trimPrefixAsterisk(zone), Value: none})
	}
	zTree.m.Set(name, Value)
	return true
}

func (zTree *cacheZoneTrees) IsDomainExist(name, domain string) (interface{}, bool){
	oldValue, isExist := zTree.m.Get(name)
	if !isExist{
		return nil, false
	}
	v, ok := oldValue.(*domainWithValue)
	if !ok{
		return nil, false
	}

	return v.Get(domain)
}

func (zTree *cacheZoneTrees) InsertZoneWithValue(name string, value *DomainContent) bool{
	oldValue, isExist := zTree.m.Get(name)
	if !isExist{
		return false
	}

	v, ok := oldValue.(*domainWithValue)
	if !ok{
		return false
	}

	v.Insert(value)
	return true
}

func (zTree *cacheZoneTrees) DeleteZone(name string, zone string) bool{
	oldValue, isExist := zTree.m.Get(name)
	if !isExist{
		return false
	}

	v, ok := oldValue.(*domainWithValue)
	if !ok{
		return false
	}

	v.Delete(zone)
	return true
}

func (zTree *cacheZoneTrees) IsExist(name string) bool{
	return zTree.m.Has(name)
}

func trimPrefixAsterisk(zone string) string{
	if zone == "*." || zone == "."{
		return zone
	}else{
		return strings.TrimPrefix(zone, "*.")
	}
}












