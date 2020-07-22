package fwcache

import (
	"github.com/infobloxopen/go-trees/iptree"
	"net"
	"strings"
)
type prefixTrieMap map[string]*iptree.Tree

func newPrefixTrieMap() prefixTrieMap{
	return make(map[string]*iptree.Tree)
}

func (ptMap prefixTrieMap) InsertIPSlice(mark string, ips []string){
	for _, ip := range ips{
		ipSplit := strings.Split(ip, "/")
		if len(ipSplit) == 0{
			continue
		}
 		if len(ipSplit) == 2{
			_, ipNet, err := net.ParseCIDR(ip)
			if err != nil{
				continue
			}
			ptMap[mark] = ptMap[mark].InsertNet(ipNet, struct{}{})
		}else{
			ipNet := net.ParseIP(ip)
			ptMap[mark] = ptMap[mark].InsertIP(ipNet, struct{}{})
		}
	}
}

func (ptMap prefixTrieMap) CreateSource(mark string) {
	var treeNode = new(iptree.Tree)
	ptMap[mark] = treeNode
}

func (ptMap prefixTrieMap) FindIP(mark string, ip string) bool{
	_, ok := ptMap[mark]
	if ok{
		ipNet := net.ParseIP(ip)
		_, isExist := ptMap[mark].GetByIP(ipNet)
		return isExist
	}

	return false
}