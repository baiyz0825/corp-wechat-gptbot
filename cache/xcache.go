package xcache

import (
	"time"

	"github.com/baiyz0825/corp-webot/utils/xlog"
	gocache "github.com/patrickmn/go-cache"
)

var cacheDb *gocache.Cache

func init() {
	xlog.Log.Info("初始化缓存中....")
	cacheDb = gocache.New(time.Hour, 2*time.Hour)
	xlog.Log.Info("初始化缓存成功，默认2h清理全局缓存")
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
		xlog.Log.Error("SetCache failure: %v", err)
		return false
	}
	return true
}

// GetUserCacheKey
// @Description: 生成上下文key
// @param keyFactor
// @return string
func GetUserCacheKey(keyFactor string) string {
	return "gpt_chat/" + keyFactor + "/" + "context"
}
