package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var cacheDb *gocache.Cache

func init() {
	log.Info("初始化缓存中....")
	cacheDb = gocache.New(time.Hour, 2*time.Hour)
	log.Info("初始化缓存成功，默认缓存时间1h，2h清理缓存")
}

func GetCacheDb() *gocache.Cache {
	return cacheDb
}

// GetDataFromCache 从缓存中获取值
func GetDataFromCache(key string) interface{} {
	if data, _, b := cacheDb.GetWithExpiration(key); b && data != nil {
		return data
	} else {
		return nil
	}
}

// SetDataToCache 设置缓存（不存在 || 已过期设置成功）
func SetDataToCache(key string, data interface{}, duration time.Duration) bool {
	err := cacheDb.Add(key, data, duration)
	if err != nil {
		log.Error("SetCache failure: %v", err)
		return false
	}
	return true
}
