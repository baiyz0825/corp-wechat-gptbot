package openaiutils

import (
	"context"

	"github.com/baiyz0825/corp-webot/config"
	"github.com/baiyz0825/corp-webot/utils/xhttp"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
	"github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

func init() {
	xlog.Log.Info("初始化openai工具SDK......")
	clientConfig := openai.DefaultConfig(config.GetGptConf().Apikey)
	if len(config.GetGptConf().URL) > 0 {
		clientConfig.BaseURL = config.GetGptConf().URL
	}
	clientConfig.HTTPClient = xhttp.HttpClient
	openaiClient = openai.NewClientWithConfig(clientConfig)
}

// SendReqAndGetTextResp 发送请求
func SendReqAndGetTextResp(msg []openai.ChatCompletionMessage) string {
	// 获取上下文数据
	data := openai.ChatCompletionRequest{
		Model:    config.GetGptConf().Model,
		Messages: msg,
		Stream:   false,
		User:     config.GetGptConf().UserName,
	}
	response, err := openaiClient.CreateChatCompletion(context.Background(), data)
	if err != nil {
		xlog.Log.Errorf("CreateCompletionStream returned error: %v", err)
		return ""
	}

	xlog.Log.WithField("data:", response).Debug("获取的数据是：")

	return response.Choices[0].Message.Content
}

// SendReqAndGetImageResp
// @Description: openai 生成图片
// @param promptMsg
// @return string
func SendReqAndGetImageResp(promptMsg string) string {
	request := openai.ImageRequest{
		Prompt:         promptMsg,
		Size:           xconst.IMAGE_TYPE_SIZE_MIDDLE,
		ResponseFormat: xconst.IMAGE_TYPE_URL,
		User:           config.GetGptConf().UserName,
	}
	image, err := openaiClient.CreateImage(context.Background(), request)
	if err != nil {
		xlog.Log.Errorf("CreateImageCompletion returned error: %v", err)
		return ""
	}
	xlog.Log.WithField("data:", image.Data).Debug("获取的数据是：")
	return image.Data[0].URL
}
