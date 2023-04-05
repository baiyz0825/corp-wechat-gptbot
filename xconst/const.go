package xconst

const (
	USER_DAO_FIRST_CREATE = "用户第一次使用，创建用户中..."
	USER_DAO_SEARCH_ERR   = "查询用户数据失败"
	USER_DAO_INSERT_ERR   = "查询用户数据创建失败"
)

const (
	AI_DEFAULT_PROMPT        = ""
	AI_DEFAULT_MSG           = "稍后再试试，小助理开小差了"
	AI_API_ERROR_MSG         = "稍后再试试，小助理的老爸去出差了"
	AI_KNOWN_YOUR_ASK        = "遵命，主人"
	AI_CLEAR_CONTEXT_SUCCESS = `主人已经开启新的对话啦，你可以使用"@prompt-set 新的角色（提示）"给我啦~`
)

const (
	COMMAND_HELP              = "@help"
	COMMAN_GPT_DELETE_CONTEXT = "@clear"
	COMMAN_GPT_IMAGE          = "@image"
	COMMAN_GPT_PROMPT_SET     = "@prompt-set"
	COMMAN_GPT_EXPORT         = "@export"
)

const (
	PROMPT_DEFAULT = "请全程使用中文与我对话"
)

// GetDefaultNoticeMenu
// @Description: 默认 提示消息
// @return string
func GetDefaultNoticeMenu() string {
	// 默认 提示消息
	return `这里是帮助菜单（如下是支持的菜单，以下不存在默认不进行处理）：
@help：帮助菜单 -> 例子：@help
@clear：清除聊天上下文 -> 例子：@clear
@image: 根据你的描述生成图片 -> 例子：@image 生成一只黑色的猫
@prompt-set：设置默认角色描述 -> 例子：@prompt-set 你是一个资深的程序员
@export：导出你的本次对话内容 -> 例子：@export`
}
