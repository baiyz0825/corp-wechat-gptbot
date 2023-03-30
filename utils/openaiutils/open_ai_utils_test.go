package openaiutils

import (
	"fmt"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestSendReqAndGetResp(t *testing.T) {
	str := "curl是什么"
	msg := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: str,
		}}
	resp := SendReqAndGetResp(msg)
	fmt.Println(resp)
}
