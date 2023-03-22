package xstring

import (
	"fmt"
	"strings"

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
