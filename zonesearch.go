package fwcache

import "strings"

type zoneTree struct {
	next   map[string]*zoneTree
	isMark bool
}

func newZTree() *zoneTree {
	root := new(zoneTree)
	root.isMark = false
	root.next = make(map[string]*zoneTree)
	return root
}

func (z *zoneTree) Insert(zone string) {
	if zone == "." || zone == "*." {
		z.isMark = true
		return
	}
	zoneList := strings.Split(zone, ".")
	for i := len(zoneList) - 1; i >= 0; i-- {
		// zTree 的 末尾是"*"时，这里需要将其标志为true
		if zoneList[i] == "*" {
			z.isMark = true
			z.next["*"] = new(zoneTree)
			break
		}else{
			if z.next[zoneList[i]] == nil {
				node := new(zoneTree)
				node.next = make(map[string]*zoneTree)
				node.isMark = false
				z.next[zoneList[i]] = node
			}
			z = z.next[zoneList[i]]
		}
	}
	z.isMark = true
}

func (z *zoneTree) InsertWithSlice(zones []string) {
	for _, zone := range zones {
		z.Insert(zone)
	}
}

func (z *zoneTree) Search(domain string) bool {
	if z.isMark {
		return true
	}

	zoneList := strings.Split(domain, ".")
	for i := len(zoneList) - 1; i >= 0; i-- {
		if z == nil {
			return false
		} else {
			// 如果一个domain走到一个next是"*"时, 则认为其是该zone下的子域, 直接返回true, 可以认为此处是一个泛域的匹配。
			if z.next["*"] != nil{
				return true
			}else{
				// 当命中一个并且刚好跟原来的是一个值时,则返回true,可以认为此处的匹配是一个域名的匹配。
				z = z.next[zoneList[i]]
				if z != nil && z.isMark && i == 0{
					return true
				}
			}
		}
	}
	return false
}
