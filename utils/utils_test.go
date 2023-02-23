package utils

import (
	"fmt"
	"testing"

	string2 "corp-webot/utils/string"
	"github.com/sirupsen/logrus"
)

func TestTransBytesToMarkdownStr(t *testing.T) {
	raw := "以下是一个使用 Golang 编写的简单的 HTTP 服务器，它将在本地端口 8080 上监听请求，并在收到请求时返回 \\\"Hello, World!\\\"：\\n\\n```go\\npackage main\\n\\nimport (\\n    \\\"fmt\\\"\\n    \\\"net/http\\\"\\n)\\n\\nfunc main() {\\n    http.HandleFunc(\\\"/\\\", func(w http.ResponseWriter, r *http.Request) {\\n        fmt.Fprint(w, \\\"Hello, World!\\\")\\n    })\\n\\n    err := http.ListenAndServe(\\\":8080\\\", nil)\\n    if err != nil {\\n        panic(err)\\n    }\\n}\\n```\\n\\n这个程序首先使用 `http.HandleFunc` 函数来注册一个路由处理器函数，它将在根路径 (\\\"/\\\") 上响应请求。当请求到达时，该函数将向 `http.ResponseWriter` 对象写入一个 \\\"Hello, World!\\\" 的消息。\\n\\n然后，程序使用 `http.ListenAndServe` 函数来启动 HTTP 服务器并监听来自客户端的请求。如果启动服务器时发生错误，程序将会抛出异常。"

	logrus.Infof(string2.TransBytesToMarkdownStr(raw))
}
func TestCutString(t *testing.T) {
	// A string
	s := "Hello-World"

	// Get first 5 values of string
	sub := s[0:5]

	// Print substring
	fmt.Println(sub)
}
