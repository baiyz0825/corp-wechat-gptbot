package doc

import (
	"bytes"
	"html/template"
	"os"
)

// Content 待渲染到HTML模板上的数据内容实例
type Content struct {
	Css   template.CSS
	Title string
	Body  template.HTML
}

// Render 基于模板结合HTML实例数据，渲染HTML
func (tpl *Tmpl) Render(c *Content) (html string, err error) {
	t := template.Must(template.New("md").Parse(tpl.content))

	var buf = &bytes.Buffer{}
	err = t.Execute(buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MarkdownCSS 返回渲染markdown用的默认css样式
func MarkdownCSS(f string) (style string, err error) {
	// 默认css
	if f == "" {
		WorkPath, _ := os.Getwd()
		file, err := os.ReadFile(WorkPath + "/assert/github_markdown.css")
		if err != nil {
			return "", err
		}
		return string(file), nil
	}
	file, err := os.ReadFile(f)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
