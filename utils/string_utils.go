package utils

import (
	"fmt"
	"strings"
)

// TransBytesToMarkdownStr 将数据转为markdown格式
func TransBytesToMarkdownStr(raw string) string {
	output := fmt.Sprintf("```\n%s\n```", raw)
	output = strings.Replace(output, "\\", "\\\\", -1)
	output = strings.Replace(output, "\"", "\\\"", -1)
	return output
}
