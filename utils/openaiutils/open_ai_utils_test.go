package openaiutils

import (
	"fmt"
	"testing"
)

func TestSendReqAndGetResp(t *testing.T) {
	str := "curl是什么"
	resp := SendReqAndGetResp(str)
	fmt.Println(resp)
}
