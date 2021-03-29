package fwcache

import (
	"fmt"
	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
	"testing"
)

func TestZoneSearch(t *testing.T) {
	zone := new(domaintree.Node)
	name, _ := domain.MakeNameFromString(".")
	zone.InplaceInsert(name, true)
	name1, _ := domain.MakeNameFromString("www.baidu.com")
	_, isMatch := zone.Get(name1)
	fmt.Println("isMatch: ", isMatch)
}
