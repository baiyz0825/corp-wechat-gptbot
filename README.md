# 注意
项目如有问题，**优先在Issue中提问**，提问请附上相应的问题出现情况和相关场景，否则不予查看
对于功能有优化建议也可直接在Issue中提问
# 应用介绍
开发本应用的想法很偶然，一方面学习golang以及Docker，另一方面看目前ChatGPT等AI工具非常火爆，搭建页面使用还需要开网页，对于经常办公场景使用微信来说不是很方便，因此自己结合网上一些示例以及Google大法开发了这个企业微信应用，通过配置即可快速搭建自己的Chat助手。目前企业微信自建应用需要配置可信IP（这个有固定IP的忽略即可，没有的可以考虑使用网上穿透工具或者购买相应的vps部署），项目目前不是特别完善，希望各位大佬指点。

---
**支持功能：**
- 支持上下文对话，记住对话次数可配置（配置的为自己发出的消息个数）
- 支持生成图片
- 支持设置自己的Prompt
- 支持导出聊天记录（目前为导出最新一批聊天记录，后续支持导出选择某一次｜全部）

> ***快捷命令可以直接使用微信菜单触发，详细参考下文***
![image](https://user-images.githubusercontent.com/81071870/231260124-c3af3b1e-299e-41ab-8ccf-b66404171dae.png)

手动使用命令：
1. @help：帮助菜单 -> 例子：@help
获取系统指令菜单
2. @clear：清除聊天上下文 -> 例子：@clear   
清除当前会话的角色设置，以及当前聊天上下文信息
3. @image: 根据你的描述生成图片 -> 例子：@image 生成一只黑色的猫
4. @prompt-set：设置默认角色描述 -> 例子：@prompt-set 你是一个资深的程序员
设置系统提示词，充当角色
5. @export：导出你的本次对话内容 -> 例子：@export
> **注意：每次只会导出最新的一条对话记录，暂不支持删除全部服务端历史记录**
导出对话内容为pdf
---
**效果**：
![image](https://user-images.githubusercontent.com/81071870/230319030-8a6d6b98-3b36-4a5e-8c24-762e6c24f410.png)
另外如果觉得有需要改进的地方，可以提Issu，空闲时间会看
# 创建合理提示词语
- [ChatGPT Shortcut - 简单易用的 ChatGPT 快捷指令表，让生产力倍增！标签筛选、关键词搜索和一键复制 Prompts | Tag filtering, keyword search, and one-click copy prompts (newzone.top)](https://ai.newzone.top/)
![image](https://user-images.githubusercontent.com/81071870/230319572-d2311d44-4786-4be2-a87b-d53355d0f49f.png)

- [ChatGPT Prompt Generator - a Hugging Face Space by merve](https://huggingface.co/spaces/merve/ChatGPT-prompt-generator)
![image](https://user-images.githubusercontent.com/81071870/230319636-7e1c32f3-b1d5-495e-81b0-07065b71aae8.png)

- [设计师灵感助手 (aigenprompt.com)](https://www.aigenprompt.com/zh-CN)
![image](https://user-images.githubusercontent.com/81071870/230319851-c4d8c613-c1ce-471f-b8e1-ba73a53dc598.png)

> 更多GPT技巧可以参考：[gpt中文调教指南](https://github.com/yzfly/awesome-chatgpt-zh)
# 部署方式 && 配置
## 创建配置文件(通用)
### 系统配置
可以下载仓库路径下的`config/config.yaml.example`
```yaml
systemConf:
  proxy: xxxx# http代理地址
  port: 50008 # 允许端口
  log: info # panic, fatal, error, warn, info, debug, trace
  callBackUrl: http://127.0.0.1 # 回调地址
  msgMode: markdown # markdown 或者 text
  logConf:
    logLevel: "debug"
    logOutPutMode: console #console file both
    logOutPutPath: "./log/crop-bot-run.log"
    logFileDateFmt: "2006-01-02 15:04:05"
    logFileRotationCount: 10
    logFormatter: text # text json
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
  weApiEncodingKey: xxxxx # key
  weChatApiAddr:  https://qyapi.weixin.qq.com # 企业微信推送api地址
```
> 支持配置模型类型（[sashabaranov/go-openai: OpenAI API](https://github.com/sashabaranov/go-openai)）
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

### 企业微信侧
1. 进入：![image](https://user-images.githubusercontent.com/81071870/230323457-1a37ea0a-0b5f-41e8-ac3e-fc0a499d2ce6.png)
2. 获取回掉地址 && 回掉Token
![image](https://user-images.githubusercontent.com/81071870/230323318-05c199bf-bd9b-4449-aa3c-08c4ff1c7389.png)
**这里回掉地址请求路径是`你的ip或者域名:端口/gpt` 服务部署之后可以浏览器访问 `你的ip或者域名:端口/test`响应为pong则正常**
3. 获取corpSecret && 应用id
![image](https://user-images.githubusercontent.com/81071870/230324473-c16c5065-2080-49ec-8321-42d26065ec3b.png)
4. **后置操作**（服务正常启动之后再，点击保存API回掉）
![image](https://user-images.githubusercontent.com/81071870/230324760-b55e40aa-8c08-4c6a-a212-b66bcbb735fb.png)

#### 微信菜单制作
微信官方接口文档：
1. 获取应用Token：https://developer.work.weixin.qq.com/document/path/91039
```shell
# 示例curl 
curl --request POST \
  --url 'https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=你的企业id&corpsecret=你的应用secret'
```
企业id：
![image](https://user-images.githubusercontent.com/81071870/231260737-b9ad4eb9-871b-4574-b63e-92f374b9e533.png)

应用secret:
![image](https://user-images.githubusercontent.com/81071870/231260620-d6601e65-9aef-497a-990a-5bf54b5647b0.png)

2. 配置菜单：https://developer.work.weixin.qq.com/document/path/90231
```shell
# 示例curl
curl --request POST \
  --url 'https://qyapi.weixin.qq.com/cgi-bin/menu/create?access_token=你的accessToken&agentid=你的应用id' \
  --header 'content-type: application/json' \
  --data '{
    "button": [
        {
            "name": "快捷指令",
            "sub_button":[
                {
                    "type": "click",
                    "name": "清除上下文",
                    "key": "@clear"
                },
                {
                    "type": "click",
                    "name": "导出",
                    "key": "@export"
                },
                {
                    "type": "click",
                    "name": "帮助",
                    "key": "@help"
                }
            ]
        }
    ]
}'
```
应用id：
![应用id](https://user-images.githubusercontent.com/81071870/231260435-739bcbf4-c1de-47c0-bf95-6e6fca1c1d71.png)
## 使用Docker
```sh
docker run -d \
  --name=gpt-webot \
  --net=host \
  -p 8989:50008 \
  -v YouPath:/apps/config \
  -v YouPath:/apps/db \
  -v YouPath:/apps/logs \
  -e GIN_MODE=release \
  --restart=always \
  ghcr.io/baiyz0825/corp-wechat-gptbot:main
```
示例：
```shell
docker run -d \
--name=gpt-webot \
--net=host \
-p 50008:50008 \
-v /home/byz/gpt/config:/apps/config \
-v /home/byz/gpt/db:/apps/db \
-v /home/byz/gpt/logs:/apps/logs \
-e GIN_MODE=release \
--restart=always \
ghcr.io/baiyz0825/corp-wechat-gptbot:main
```
## 使用Docker Compose
1. 下载仓库中的docker-composer.yaml
```yaml
version: '3.8'
services: 
  crop-gpt-bot:
   container_name: gpt-bot
   image: ghcr.io/baiyz0825/corp-wechat-gptbot:main
   ports:
    - "本机端口":"50008"
   volumes:
    - 你的配置文件路径:/apps/config:rw 
    - 你的数据库存储路径:/apps/db:rw 
    - 你的日志路径:/apps/logs:rw     
   environment:
    - GIN_MODE=release
   deploy:         
    restart_policy:
      condition: on-failure
      delay: 5s
      max_attempts: 3
      window: 120s     
   restart: always
   networks:
     - host # bridge/host
```
2. 执行 `docker compose -f 你的路径/docker-compose.yaml up -d`
3. 查看运行状态 `docker compose ps`
4. 查看日志
- 日志文件：`tail -f 你的日志文件`
- Docker控制台：`docker logs -f 容器id`

# 采用方案
## 数据存储
数据存储使用sqlite3，在使用apline打包出现cgo问题（由于sqlite3依赖系统动态库，最后使用ubuntu镜像打包），请自行映射存储数据db
## 缓存
### 用户多次请求缓存
微信侧，在回调过程中，如果未及时收到消息会，进行多次发送消息，防止由于网络原因以及相对应的错误和限制用户重复发送消息，设置一个临时用户缓存
![image](https://user-images.githubusercontent.com/81071870/230303669-9643e141-f84e-4a55-8fa1-9605422ed8dd.png)
缓存设置规则：计算当前请求消息内容的HASH
```go
hashInt64 := xstring.HashDataConcurrently([]byte(data.Content))
cacheKey := data.FromUsername + ":" + strconv.FormatInt(hashInt64, 10)
```
### 用户消息上下文缓存
每次用户请求获取Openai消息时，先检查消息上下文是否存在缓存，存在的话，请求API成功之后加入当前用户msg以及Openai响应，保持上下文，当用户清楚缓存时@clear进行数据持久化（删除用户设置的第一次提示词语），方便下次导出。
缓存Key规则：`gpt_chat/" + keyFactor + "/" + "context`
缓存内容：“msgContext model.MessageContext”:[代码位置](https://github.com/baiyz0825/corp-webot/blob/605430c9b6/services/impl/gpt_chat_service.go#L50)
缓存逻辑：[代码位置](https://github.com/baiyz0825/corp-webot/blob/605430c9b6/services/impl/gpt_chat_service.go#L48-L108)
## 使用  wkhtmltopdf 转化pdf
> pdf导出空白问题：需要安装字体 `sudo cp ./assert/simsun.ttc /usr/share/fonts`（dockerfile中已经打包）
## 数据表结构
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
