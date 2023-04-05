package doc

import (
	"fmt"
	"os"
	"testing"
)

func TestGenDoc(t *testing.T) {
	file, err := os.ReadFile("./assert/11.md")
	if err != nil {
		fmt.Println(err)
	}
	bytes := []byte("\n\n## 这是新加目录\n\n")
	file = append(file, bytes...)
	data, err := GetHtmlFromMdBytes("测试", file)
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("./test.html", []byte(data), 0666)

}

func TestConvertHtmlToPDF(t *testing.T) {
	file, err := os.ReadFile("./test.html")
	if err != nil {
		fmt.Println(err)
	}
	pdf := ConvertHtmlToPDF(file)
	os.WriteFile("./test_trans.pdf", []byte(pdf), 0666)
}

func TestWorkPath(t *testing.T) {
	WorkPath, _ := os.Getwd()
	fmt.Println(WorkPath)
}
