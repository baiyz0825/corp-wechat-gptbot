package doc

import (
	"strings"

	"github.com/russross/blackfriday/v2"
)

// Parser markdown解析器
type Parser interface {
	// Parse 通过解析器解析MD，
	Parse(md []byte) (html []byte)

	Markdown2HTML(mdContent []byte) (html string, err error)
}

// BlackFriday Markdown解析器 https://github.com/russross/blackfriday
type BlackFriday struct {
	Option blackfriday.Extensions
}

// NewBlackFriday 初始化一个BlackFriday解析器
func NewBlackFriday() *BlackFriday {
	return &BlackFriday{}
}

// Parse 通过BlackFriday解析器解析MD
func (p *BlackFriday) Parse(input []byte) (html []byte) {
	return blackfriday.Run(input)
}

// Markdown2HTML 将指定目录位置的markdown文件，解析生成html文档
func (p *BlackFriday) Markdown2HTML(mdContent []byte) (html string, err error) { // 解析md文件
	var s = strings.Builder{}
	s.Write(p.Parse(mdContent))
	return s.String(), nil
}
