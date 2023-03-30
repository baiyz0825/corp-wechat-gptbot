package openaiutils

import (
	"context"

	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/utils/xhttp"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

var openaiClient *openai.Client

func init() {
	logrus.Info("初始化openai工具SDK......")
	clientConfig := openai.DefaultConfig(config.GetGptConf().Apikey)
	clientConfig.HTTPClient = &xhttp.HttpClient
	openaiClient = openai.NewClientWithConfig(clientConfig)
}

// SendReqAndGetResp 发送请求
func SendReqAndGetResp(msg []openai.ChatCompletionMessage) string {
	// 获取上下文数据
	data := openai.ChatCompletionRequest{
		Model:    config.GetGptConf().Model,
		Messages: msg,
		Stream:   false,
		User:     config.GetGptConf().UserName,
	}
	response, err := openaiClient.CreateChatCompletion(context.Background(), data)
	if err != nil {
		logrus.Errorf("CreateCompletionStream returned error: %v", err)
		return ""
	}

	logrus.WithField("data:", response).Info("获取的数据是：")

	return response.Choices[0].Message.Content
}
