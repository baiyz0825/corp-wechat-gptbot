package doc

import (
	"os"
)

var indexTpl = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.Title}}</title>
		<style type="text/css">
		{{.CSS}}
		</style>
</head>
<body>
<article class="markdown-body">
{{.Body}}
</article>
</body>
</html>
`
var defaultTmpl = &Tmpl{indexTpl}

// Tmpl 渲染模板
type Tmpl struct {
	content string
}

// NewTmpl 初始化Markdown渲染模板
func NewTmpl(f string) (*Tmpl, error) {
	// 默认css
	if f == "" {
		WorkPath, _ := os.Getwd()
		file, err := os.ReadFile(WorkPath + "/assert/index.tpl")
		if err != nil {
			return nil, err
		}
		return &Tmpl{string(file)}, nil
	}
	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Tmpl{string(file)}, nil
}
