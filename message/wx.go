package message

import (
	"sync"

	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"person-bot/cache"
)

var wxAksk *ExpireAksk

type ExpireAksk struct {
	cache *gocache.Cache
	mutex sync.Mutex
}

func init() {
	wxAksk = &ExpireAksk{
		cache: cache.GetCacheDb(),
		mutex: sync.Mutex{},
	}
}

func (e *ExpireAksk) GetAccessKey() string {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	fromCache := cache.GetDataFromCache("aksk")
	if fromCache == nil {
		// 请求获取

	}
	value, ok := fromCache.(string)
	if !ok {
		log.Error("aksk 缓存获取类型不一致")
	}
	return value

}
