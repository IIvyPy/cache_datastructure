package fwcache

import "errors"

type DynamicArray struct {
	logicalSize uint16
	capacity    uint16
	container   []interface{}
}

func (s *DynamicArray) SetSize(length uint16) {
	newContainer := make([]interface{}, length)
	copy(newContainer, s.container)
	s.capacity = length
	s.container = newContainer
}

func (s *DynamicArray) expandDouble() {
	if s.capacity == 0 {
		s.SetSize(1)
	} else {
		s.SetSize(s.capacity * 2)
	}
}

func (s *DynamicArray) expandCapacity(size uint16) {
	s.SetSize(size)
}

// shrink method resizes the dynamic array to half
// it's size. Leaves the blanks filled with nil values
func (s *DynamicArray) shrink() {
	if s.capacity >= 2 {
		s.SetSize(s.capacity / 2)
	} else if s.capacity == 1 {
		s.SetSize(1)
	}
}

// resize automatically expands the dynamic array if an
// insertion is required and it is already full. Or shrinks
// it whenever it is only a quarter full
func (s *DynamicArray) resize(newLogicalSize uint16) {
	if newLogicalSize >= s.capacity {
		if newLogicalSize <= s.capacity*2{
			s.expandDouble()
		}else{
			s.expandCapacity(newLogicalSize)
		}
	} else if newLogicalSize > 0 && newLogicalSize == (s.capacity/4) {
		s.shrink()
	}
}

// Get retrieves the value store at the given index
// in dynamic array. Throws an error if index out of range.
func (s *DynamicArray) Get(index uint16) (interface{}, error) {
	if index < s.logicalSize {
		return s.container[index], nil
	}
	return nil, errors.New("index out of range")
}

// Delete will delete the value stored at the given index
// in dynamic array. Throws an error if index out of range.
func (s *DynamicArray) Delete(index uint16) error{
	if index < s.logicalSize {
		s.container[index] = nil
		return nil
	}
	return errors.New("index out of range")
}

// Range retrieves a dynamic store slice with the provided
// edges. Throws an error if index out of range.
func (s *DynamicArray) Range(start uint16, stop uint16) ([]interface{}, error) {
	if start >= 0 && stop < s.logicalSize {
		return s.container[start:stop], nil
	}
	return nil, errors.New("index out of range")
}

// Insert inserts a value in the array at the index pos
// (0 indexed). Resizes the array when needed.
func (s *DynamicArray) Set(index uint16, elem interface{}) {
	if index >= 0 && index < s.logicalSize{
		s.container[index] = elem
		return
	}

	s.logicalSize = index + 1
	newLogicalSize := s.logicalSize
	s.resize(newLogicalSize)
	s.container[index] = elem
	return
}

func (s *DynamicArray) Length() uint16{
	return s.logicalSize
}

func (s *DynamicArray) Cap() uint16{
	return s.capacity
}