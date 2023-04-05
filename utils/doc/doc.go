package doc

import (
	"bytes"
	"html/template"
	"io"

	pdfC "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// GetHtmlFromMdBytes
//  @Description: 从md bytes渲染html
//  @param title
//  @param mdBytes
//  @return html
//  @return err
//
func GetHtmlFromMdBytes(title string, mdBytes []byte) (html string, err error) {
	// body parser
	var p Parser
	p = NewBlackFriday()
	body, err := p.Markdown2HTML(mdBytes)
	if err != nil {
		return "", err
	}

	// content builder
	css, err := MarkdownCSS("")
	if err != nil {
		return "", err
	}
	c := &Content{
		Css:   template.CSS(css),
		Title: title,
		Body:  template.HTML(body),
	}

	// html render
	tmpl, err := NewTmpl("")
	if err != nil {
		return "", err
	}
	return tmpl.Render(c)
}

// ConvertHtmlToPDF
//  @Description: 转化html到pdf 需要安装 sudo apt install wkhtmltopdf 并且安装中午字体
//  sudo cp ./assert/simsun.ttc /usr/share/fonts
//  @param html
//  @return []byte
//
func ConvertHtmlToPDF(html []byte) []byte {
	// Create new PDF generator
	pdfg, err := pdfC.NewPDFGenerator()
	if err != nil {
		return nil
	}

	// dpi设置
	pdfg.Dpi.Set(300)
	// 竖页面
	pdfg.Orientation.Set(pdfC.OrientationPortrait)
	pdfg.Grayscale.Set(true)

	// data
	page := pdfC.NewPageReader(io.Reader(bytes.NewBuffer(html)))

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil
	}
	// 返回转换字节流
	return pdfg.Bytes()
}
