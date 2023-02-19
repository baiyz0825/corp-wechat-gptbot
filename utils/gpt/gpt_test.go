package gpt

import (
	"fmt"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	accessToken := GetAccessToken()
	fmt.Println(accessToken)
}
