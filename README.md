# 支持模型
```go
package openai
const (
GPT432K0314             = "gpt-4-32k-0314"
GPT432K                 = "gpt-4-32k"
GPT40314                = "gpt-4-0314"
GPT4                    = "gpt-4"
GPT3Dot5Turbo0301       = "gpt-3.5-turbo-0301"
GPT3Dot5Turbo           = "gpt-3.5-turbo"
GPT3TextDavinci003      = "text-davinci-003"
GPT3TextDavinci002      = "text-davinci-002"
GPT3TextCurie001        = "text-curie-001"
GPT3TextBabbage001      = "text-babbage-001"
GPT3TextAda001          = "text-ada-001"
GPT3TextDavinci001      = "text-davinci-001"
GPT3DavinciInstructBeta = "davinci-instruct-beta"
GPT3Davinci             = "davinci"
GPT3CurieInstructBeta   = "curie-instruct-beta"
GPT3Curie               = "curie"
GPT3Ada                 = "ada"
GPT3Babbage             = "babbage"
)
```

# 部署方式
1. 新建配置文件在 `YourPath`
```yaml
systemConf:
  proxy: xxxx# http代理地址
  port: 50008 # 允许端口
  log: info # panic, fatal, error, warn, info, debug, trace
  logPath: ./logs/log #日志文件位置
  callBackUrl: http://127.0.0.1 # 回调地址
  msgMode: markdown # markdown 或者 text
gptConfig:
  apikey: xxxx # openapi key
  model: gpt-3.5-turbo # 对话模型
  UserName: xxx # 用户名
  url: https://api.openai.com/v1 # 请求基地址
weConfig:
  corpid: xxxxx # 企业id
  corpSecret: xxxxx #应用密码
  agentId: xxxx # 应用Id
  weApiRCallToken: xxxxx # 调用token
  weApiEncodingKey: xxxxx # jiami
  weChatApiAddr:  https://qyapi.weixin.qq.com # 企业微信推送api地址
```

2. 运行容器
```sh
docker run -d \
  --name=gpt-webot \
  --net=host \
  -p 8989:50008 \
  -v YouPath:/apps/config \
  -v YouPath:/apps/db \
  -v YouPath:/apps/logs
  --restart=always \
  ghcr.io/baiyz0825/corp-webot:main
```
Example:
```shell
docker run -d \
--name=gpt-webot \
--net=host \
-p 50008:50008 \
-v /home/byz/gpt/config:/apps/config \
-v /home/byz/gpt/db:/apps/db \
-v /home/byz/gpt/logs:/apps/logs \
--restart=always \
ghcr.io/baiyz0825/corp-webot:main
```
3. 数据表创建
```sqlite
create table if not exists "users" (
    id          INTEGER
        primary key autoincrement,
    name        CHAR(50) not null ,
    sys_prompt  CHAR(512),
    update_time BIGINT
);
create table if not exists "context"
(
    id          INTEGER
        primary key autoincrement,
    name        CHAR(50) not null,
    context_msg CHAR(512),
    update_time BIGINT
);

```

# 使用  wkhtmltopdf 转化pdf
dockerfile中已经打包
需要安装字体 `sudo cp ./assert/simsun.ttc /usr/share/fonts`

# 命令使用方法
1. @help：帮助菜单 -> 例子：@help
获取系统指令菜单
2. @clear：清除聊天上下文 -> 例子：@clear   
清除当前会话的角色设置，以及当前聊天上下文信息
3. @image: 根据你的描述生成图片 -> 例子：@image 生成一只黑色的猫
4. @prompt-set：设置默认角色描述 -> 例子：@prompt-set 你是一个资深的程序员
设置系统提示词，充当角色
5. @export：导出你的本次对话内容 -> 例子：@export
导出对话内容为pdf

**注意：每次只会导出最新的一条对话记录，暂不支持删除全部服务端历史记录**

# 数据存储
数据db文件存储在容器中/apps/db 中
