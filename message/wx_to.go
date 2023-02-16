package message

import (
	"encoding/xml"
)

type XML struct {
	XMLName      xml.Name `xml:"xml"`
	Text         string   `xml:",chardata"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        string   `xml:"MsgId"`
	AgentID      string   `xml:"AgentID"`
}
