package fwcache

import "strings"

type ZoneTree struct {
	next   map[string]*ZoneTree
	isMark bool
}

func NewZTree() *ZoneTree {
	root := new(ZoneTree)
	root.isMark = false
	root.next = make(map[string]*ZoneTree)
	return root
}

func (z *ZoneTree) Insert(zone string) {
	if zone == "." || zone == "*." {
		z.isMark = true
		return
	}
	zoneList := strings.Split(zone, ".")
	for i := len(zoneList) - 1; i >= 0; i-- {
		// zTree 的 末尾是"*"时，这里需要将其标志为true
		if zoneList[i] == "*" {
			z.isMark = true
			z.next["*"] = new(ZoneTree)
			break
		}else{
			if z.next[zoneList[i]] == nil {
				node := new(ZoneTree)
				node.next = make(map[string]*ZoneTree)
				node.isMark = false
				z.next[zoneList[i]] = node
			}
			z = z.next[zoneList[i]]
		}
	}
	z.isMark = true
}

func (z *ZoneTree) InsertWithSlice(zones []string) {
	for _, zone := range zones {
		z.Insert(zone)
	}
}

func (z *ZoneTree) Search(domain string) bool {
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
