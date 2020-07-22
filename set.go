package fwcache

// Set is one of the primary interface provided by the cachedata package.
// the set only support set-level operation.
// the element-level in the set is not supported.
type Set interface {
	// check whether the input item is in the set
	Contains(interface{}) bool

	// clear all the element in the set
	Clear()

	// return the number of the element in the set
	Len() int

	// return the slice type of the set
	ToSlice() []interface{}

	// add to set
	Add(interface{}) bool

	// remove the value
	Remove(i interface{})

	// add slice to set
	AddWithSlice([]interface{})
}

// NewThreadUnsafeSet return a blank thread unsafe set
func NewThreadUnsafeSet() Set {
	set := newThreadUnsafeSet()
	return &set
}

// NewThreadUnsafeSetFromSlice return a thread unsafe set from slice
func NewThreadUnsafeSetFromSlice(src []interface{}) Set {
	a := newThreadUnsafeSet()
	for _, item := range src {
		a.Add(item)
	}
	return &a
}
