package xconst

import (
	"corp-webot/xerror"
)

var (
	TextCommandMap = map[string]bool{
		"@gpt": true,
	}
)

var (
	GPTREDIRECTERR    = xerror.NewGPTError(1000, "禁止重定向")
	GPTAUTHCOOKIEERR  = xerror.NewGPTError(1001, "获取授权cookies失败")
	GPTAUTHSESSIONERR = xerror.NewGPTError(1002, "解析SessionToken失败")
	GPTAUTHLASTURLERR = xerror.NewGPTError(1002, "解析最后一个URL不是重定向")
)
