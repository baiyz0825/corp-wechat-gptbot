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
  -p 50008:50008 \
  -v YouPath:/apps/config \
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