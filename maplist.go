package fwcache

import (
	"container/list"
	"sync"
)

// MapList is the data struct used to save policy-info, element in the list is the value in the map and itemID is the key in the map
type MapList struct {
	sync.RWMutex
	myMap      map[uint64]interface{}
	linkedList *list.List
}

type ItemInterface interface {
	GetPriority() uint16
}

func newMapList() *MapList {
	return &MapList{
		linkedList: list.New(),
		myMap:      make(map[uint64]interface{}),
	}
}

func (mapList *MapList) init() {
	mapList.Lock()
	defer mapList.Unlock()

	mapList.linkedList = list.New()
	mapList.myMap = make(map[uint64]interface{})
}

func (mapList *MapList) Get(id uint64) (interface{}, bool) {
	mapList.RLock()
	defer mapList.RUnlock()

	value, isOk := mapList.myMap[id]
	if !isOk {
		return nil, isOk
	}

	return value.(*list.Element).Value, isOk
}

func (mapList *MapList) pushBack(id uint64, item interface{}) {
	mapList.Lock()
	defer mapList.Unlock()

	element := mapList.linkedList.PushBack(item)
	mapList.myMap[id] = element
}

func (mapList *MapList) update(id uint64, item interface{}) bool {
	mapList.Lock()
	defer mapList.Unlock()

	value, isOk := mapList.myMap[id]
	if !isOk {
		return isOk
	}

	var element *list.Element
	prev := value.(*list.Element).Prev()
	next := value.(*list.Element).Next()

	// 删掉当前节点
	mapList.linkedList.Remove(value.(*list.Element))
	// 将需要更新的节点插入进去
	if prev != nil{
		element = mapList.linkedList.InsertAfter(item, prev)
	}else if next != nil{
		element = mapList.linkedList.InsertBefore(item, next)
	}else{
		element = mapList.linkedList.PushBack(item)
	}
	mapList.myMap[id] = element
	return true
}

func (mapList *MapList) insertBefore(id uint64, item interface{}, each *list.Element) {
	element := mapList.linkedList.InsertBefore(item, each)
	mapList.myMap[id] = element
}

// HasKey is used to check if the id exist.
func (mapList *MapList) hasKey(id uint64) bool {
	mapList.RLock()
	defer mapList.RUnlock()

	_, isOk := mapList.myMap[id]
	return isOk
}

func (mapList *MapList) removeByMapID(id uint64) interface{} {
	mapList.Lock()
	defer mapList.Unlock()

	value, isOk := mapList.myMap[id]
	if !isOk {
		return nil
	}

	removeValue := mapList.linkedList.Remove(value.(*list.Element))
	delete(mapList.myMap, id)

	return removeValue
}

func (mapList *MapList) add(id uint64, item interface{}) {
	for each := mapList.Front(); each != nil; each = each.Next() {
		// find first bigger number
		if item.(ItemInterface).GetPriority() > each.Value.(ItemInterface).GetPriority() {
			mapList.insertBefore(id, item, each)
			return
		}
	}
	// add to last element if not found.
	mapList.pushBack(id, item)
	return
}

func (mapList *MapList) delete(id uint64) interface{} {
	return mapList.removeByMapID(id)
}

// Traversal all the items in the policy into slice
func (mapList *MapList) traversal() []interface{} {
	mapList.RLock()
	defer mapList.RUnlock()

	itemList := make([]interface{}, 0, mapList.linkedList.Len())

	for each := mapList.Front(); each != nil; each = each.Next() {
		itemList = append(itemList, each.Value)
	}
	return itemList
}

// Front return the head of the list
func (mapList *MapList) Front() *list.Element {
	return mapList.linkedList.Front()
}

// Len return the length of the list
func (mapList *MapList) len() int {
	return mapList.linkedList.Len()
}
