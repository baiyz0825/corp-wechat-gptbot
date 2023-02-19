package test

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"corp-webot/model"
	"corp-webot/utils/gpt"
	"golang.org/x/net/http2"
)

func TestMessageModel(t *testing.T) {
	sourceXl := `<xml>
   <ToUserName><![CDATA[toUser]]></ToUserName>
   <FromUserName><![CDATA[fromUser]]></FromUserName> 
   <CreateTime>1348831860</CreateTime>
   <MsgType><![CDATA[text]]></MsgType>
   <Content><![CDATA[this is a test]]></Content>
   <MsgId>1234567890123456</MsgId>
   <AgentID>1</AgentID>
</xml>`
	message := model.RecTextMessage{}
	err := xml.Unmarshal([]byte(sourceXl), &message)
	if err != nil {
		return
	}
	fmt.Printf("读取的设置为：%v %v %v", message.AgentID, message.Content, message.MsgType)

}

func TestRuneCut(t *testing.T) {
	unicodeRune := []rune("@gpt你说我的小王子")[0:4]
	fmt.Println(string(unicodeRune))
}

func TestGPTAksk(t *testing.T) {
	gpt.GetAccessToken()
}

func TestParseUrl(t *testing.T) {
	originUrl := "https://auth0.openai.com/authorize?client_id=TdJIcbe16WoTHtN95nyywh5E4yOo6ItG&scope=openid%20email%20profile%20offline_access%20model.request%20model.read%20organization.read&response_type=code&redirect_uri=https%3A%2F%2Fexplorer.api.openai.com%2Fapi%2Fauth%2Fcallback%2Fauth0&audience=https%3A%2F%2Fapi.openai.com%2Fv1&state=I_b9zgmBonH9_nyKm3pF45sBPZeYhEIVdYNduPXP1KU&code_challenge=DjrdWVogz4KP6iQtmtbByQzqeFIO0_rckquiiEwCgxc&code_challenge_method=S256"

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	httpClient.Transport = &http2.Transport{
		// 跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", originUrl, nil)
	req.Header.Set("Authority", "auth0.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("stopped after 0 redirects")
	}
	resp, _ := httpClient.Do(req)
	location := resp.Header.Get("Location")
	state := regexp.MustCompile(`state=(.*)`).FindStringSubmatch(location)[0]
	state = strings.Split(state, "=")[1]

	defer resp.Body.Close()
	// 输出HTTP响应状态码
	fmt.Println(resp.Status)
}

func TestUrl(t *testing.T) {
	url := "https://auth0.openai.com/authorize?client_id=TdJIcbe16WoTHtN95nyywh5E4yOo6ItG&scope=openid%20email%20profile%20offline_access%20model.request%20model.read%20organization.read&response_type=code&redirect_uri=https%3A%2F%2Fexplorer.api.openai.com%2Fapi%2Fauth%2Fcallback%2Fauth0&audience=https%3A%2F%2Fapi.openai.com%2Fv1&state=I_b9zgmBonH9_nyKm3pF45sBPZeYhEIVdYNduPXP1KU&code_challenge=DjrdWVogz4KP6iQtmtbByQzqeFIO0_rckquiiEwCgxc&code_challenge_method=S256"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status code:", resp.Status)
	fmt.Println("Headers:", resp.Header)
}
