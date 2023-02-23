package wecom

import (
	"corp-webot/config"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
)

var WeComApp *work.Work

func LoadWeComAppConf() {
	WeComApp, _ = work.NewWork(&work.UserConfig{
		CorpID:  config.GetWechatConf().Corpid,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: config.GetWechatConf().AgentId,    // 内部应用的app id
		Secret:  config.GetWechatConf().CorpSecret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: "", //
			Scopes:   nil,
		},
		HttpDebug: true,
		// 可选，不传默认走程序内存
		// Cache: kernel.NewRedisClient(&kernel.RedisOptions{
		// 	Addr:     "127.0.0.1:6379",
		// 	Password: "",
		// 	DB:       0,
		// }),
	})
}
