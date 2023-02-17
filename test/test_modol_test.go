package test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"corp-webot/model"
)

func TestMessageModel(t *testing.T) {
	sourceXl := `<xml>
   <ToUserName><![CDATA[toUser]]></ToUserName>
   <FromUserName><![CDATA[fromUser]]></FromUserName> 
   <CreateTime>1348831860</CreateTime>
   <MsgType><![CDATA[text]]></MsgType>
   <Content><![CDATA[this is a test]]></Content>
   <MsgId>1234567890123456</MsgId>
   <AgentID>1</AgentID>
</xml>`
	message := model.RecTextMessage{}
	err := xml.Unmarshal([]byte(sourceXl), &message)
	if err != nil {
		return
	}
	fmt.Printf("读取的设置为：%v %v %v", message.AgentID, message.Content, message.MsgType)

}

func TestRuneCut(t *testing.T) {
	unicodeRune := []rune("@gpt你说我的小王子")[0:4]
	fmt.Println(string(unicodeRune))
}
