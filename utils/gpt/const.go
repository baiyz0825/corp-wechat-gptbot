package gpt

const (
	// V1GPTBaseUrl V1版本GPT url
	V1GPTBaseUrl = "https://chatgpt.duti.tech"
	// OfficialGPTBaseUrl 官方GPT url
	OfficialGPTBaseUrl = "https://chatgpt.duti.tech"
)

const (
	// OfficialApiAiModel 学习模型
	OfficialApiAiModel = "text-davinci-003"
	// OfficialApiPresencePenalty 话题开放程度 -2.0 ~ 2.0 谈论新问题可能性
	OfficialApiPresencePenalty = 0.6
	// OfficialApiFrequencyPenalty0 字符重复程度
	OfficialApiFrequencyPenalty0 = 0.0
	// OfficialApiTopP 从累计概率超过某一个阈值p的词汇中进行采样 词出现重复率
	OfficialApiTopP = 1
	// OfficialApiTemperature 幽默程度
	OfficialApiTemperature = 0.3
	// OfficialApiStream 是否开启流传输，对话形式蹦
	OfficialApiStream = false
	// OfficialApiToken 最大 问题+答案文本字符数
	OfficialApiToken = 2048
)

const (
	// V1ApiAction 默认模型请求action
	V1ApiAction = "next"
	// V1ApiRole 默认Api角色
	V1ApiRole = "user"
	// V1ApiContentType 默认api角色
	V1ApiContentType = "text"
	// V1ApiModel 学习模型
	V1ApiModel = "text-davinci-002-render-sha"
	// V1ApiPaid 模型paid
	V1ApiPaid = "text-davinci-002-render-paid"
)
