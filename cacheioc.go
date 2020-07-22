package fwcache

import (
	"github.com/bluele/gcache"
	"time"
)

// IocOperator is used to operator ioc info
type IocOperator interface {
	// clear the ioc cache
	IocPurge()

	// get value of the key from cache, error will not be nil if not exist.
	IocGet(key string) (interface{}, error)

	// set the key without expired
	IocSet(key string, iocInfo interface{}) error

	// set the key value with expired
	SetWithExpired(key string, iocInfo interface{}, expiration time.Duration) error

	// check if the key exists
	IocHas(key string) bool

	// get the length of cache with checkExpired
	CacheLen(checkExpired bool) int
}

type ioc struct {
	gcache.Cache
}

func NewIoc(size int) IocOperator {
	return ioc{gcache.New(size).Build()}
}

func (i ioc) IocPurge() {
	i.Purge()
}

func (i ioc) IocGet(key string) (interface{}, error) {
	return i.Get(key)
}

func (i ioc) SetWithExpired(key string, iocInfo interface{}, expiration time.Duration) error {
	return i.SetWithExpire(key, iocInfo, expiration)
}

func (i ioc) IocSet(key string, iocInfo interface{}) error {
	return i.Set(key, iocInfo)
}

func (i ioc) IocHas(key string) bool {
	return i.Has(key)
}

func (i ioc) CacheLen(checkExpired bool) int {
	return i.Len(checkExpired)
}
