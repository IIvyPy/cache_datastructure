package fwcache

import (
	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/go-trees/domaintree"
)
/**
 * @Description: insert, delete and get domain, should notice the root zone.
 */
type DomainContent struct {
	Domain string
	Value  interface{}
}

type domainWithValue struct {
	node      *domaintree.Node
	rootMark  bool
	rootValue interface{}
}

func (dv *domainWithValue) Insert(content *DomainContent) {
	if content.Domain == "." || content.Domain == "*." {
		dv.rootMark = true
		dv.rootValue = content.Value
		return
	}
	name, err := domain.MakeNameFromString(content.Domain)
	if err != nil {
		dv.node.InplaceInsert(name, content.Value)
	}
}

func (dv *domainWithValue) Delete(contentDomain string) {
	if contentDomain == "." || contentDomain == "*." {
		dv.rootMark = false
		dv.rootValue = nil
		return
	}
	name, err := domain.MakeNameFromString(contentDomain)
	if err != nil {
		newNode, success := dv.node.DeleteSubdomains(name)
		if success {
			dv.node = newNode
		}
	}
}

func (dv *domainWithValue) Get(content string) (interface{}, bool) {
	if dv.rootMark {
		return dv.rootValue, true
	}
	name, err := domain.MakeNameFromString(content)
	if err != nil {
		return dv.node.Get(name)
	}

	return nil, false
}

func (dv *domainWithValue) Clear() {
	dv.node = nil
	dv.rootMark = false
	dv.rootValue = nil
}
