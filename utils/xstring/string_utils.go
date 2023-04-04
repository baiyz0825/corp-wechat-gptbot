package xstring

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
)

// TransBytesToMarkdownStr 将数据转为markdown格式
func TransBytesToMarkdownStr(raw string) string {
	output := fmt.Sprintf("```\n%s\n```", raw)
	output = strings.Replace(output, "\\", "\\\\", -1)
	output = strings.Replace(output, "\"", "\\\"", -1)
	return output
}

// GenLogoAscii 生成ascii
func GenLogoAscii(text string, color string) {
	myFigure := figure.NewColorFigure(text, "", color, true)
	myFigure.Print()
}

const keyLength = 10

// GenerateRandomStr
// @Description: 生成随机字符串
// @return string
func GenerateRandomStr() string {
	rand.Seed(time.Now().UnixNano())
	var builder strings.Builder
	for i := 0; i < keyLength; i++ {
		builder.WriteRune(rune(rand.Intn(26) + 97)) // generate random lowercase letter
	}

	return builder.String()
}
