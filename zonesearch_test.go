package fwcache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZoneSearch(t *testing.T) {
	zone := newZTree()
	zone.InsertWithSlice([]string{"www.sina.com", "*.zte.com", "baidu.com"})
	isMatch := zone.Search("wwww.baidu.com")
	assert.Equal(t, false, isMatch)

	isMatch = zone.Search("baidu.com")
	assert.Equal(t, true, isMatch)

	isMatch = zone.Search("zte.com")
	assert.Equal(t, true, isMatch)

	isMatch = zone.Search("www.zte.com")
	assert.Equal(t, true, isMatch)

	isMatch = zone.Search("abc.www.zte.com")
	assert.Equal(t, true, isMatch)

	isMatch = zone.Search("abc.www.sina.com")
	assert.Equal(t, false, isMatch)

	isMatch = zone.Search("com")
	assert.Equal(t, false, isMatch)
}
