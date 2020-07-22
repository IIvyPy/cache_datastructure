package fwcache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIPTrie(t *testing.T) {
	ipTree := newCacheIPPrefixTrie()
	ipTree.Create("source")
	ipTree.InsertIPs("source", "", []string{"123.2.3.4", "234.2.3.2/24", "1.1.1.1"})

	isExist := ipTree.IsIPExist("source", "1", "127.0.0.1")
	assert.Equal(t, false, isExist)

	isExist = ipTree.IsIPExist("test", "1", "127.0.0.1")
	assert.Equal(t, false, isExist)

	isExist = ipTree.IsIPExist("source", "", "123.2.3.4")
	assert.Equal(t, true, isExist)

	isExist = ipTree.IsIPExist("source", "", "1.2.3.4")
	assert.Equal(t, false, isExist)
}
