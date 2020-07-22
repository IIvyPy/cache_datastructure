package fwcache

type threadUnsafeSet map[interface{}]struct{}

func newThreadUnsafeSet() threadUnsafeSet {
	return make(threadUnsafeSet)
}

func (set *threadUnsafeSet) Add(name interface{}) bool {
	_, found := (*set)[name]

	if found {
		return false
	}

	(*set)[name] = struct{}{}
	return true
}

func (set *threadUnsafeSet) AddWithSlice(names []interface{}) {
	for _, name := range names {
		(*set)[name] = struct{}{}
	}
}

func (set *threadUnsafeSet) Contains(i interface{}) bool {
	_, found := (*set)[i]
	return found
}

func (set *threadUnsafeSet) Clear() {
	*set = newThreadUnsafeSet()
}

func (set *threadUnsafeSet) Remove(i interface{}) {
	delete(*set, i)
}

func (set *threadUnsafeSet) Len() int {
	return len(*set)
}

func (set *threadUnsafeSet) ToSlice() []interface{} {
	keys := make([]interface{}, 0, set.Len())

	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}
