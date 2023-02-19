package gpt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"corp-webot/xconst"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
)

var aksk string
var httpClient *http.Client

var (
	// userName = config.GetGptConf().UserName
	// password = config.GetGptConf().Passwd
	userName = "yizhuo0825@gmail.com"
	password = "adgh267adgh267"
)

// 初始化客户端连接环境
func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logrus.Fatal("初始化cookie客户端失败")
	}
	httpClient = &http.Client{
		Timeout: time.Second * 10,
		// 需要保证同一会话，自动管理cookie
		Jar: jar,
	}

	// 防止多次重定向
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return &xconst.GPTREDIRECTERR
	}
	httpClient.Transport = &http2.Transport{
		// 跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

}

// checkV1RespCorrect 检查数据流中数据
func checkV1RespCorrect(respData *messageGPTV1Resp, traceId string) error {
	if respData == nil {
		return errors.New("检查GPT接口传递的数据响应为空，traceID：" + traceId)
	}
	if respData.Error != nil {
		return errors.New("检查GPT接口传递的数据,远程处理错误，traceID：" + traceId + "原始错误是：" + respData.Error.(string))
	}
	return nil
}

// GetAccessToken 获取aksk密钥
func GetAccessToken() string {
	// 检查连通性 https://explorer.api.openai.com/auth/login
	link := "https://explorer.api.openai.com/"
	logrus.WithField("Link", link).Debugf("阶段1:开始获取AkskToken....")
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		logrus.WithError(err).Errorf("阶段1:创建Http请求失败")
		return ""
	}

	// Prepare Set Header
	req.Header.Set("Host", "explorer.api.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	response, err := httpClient.Do(req)
	defer response.Body.Close()
	if err != nil {
		logrus.WithError(err).Errorf("阶段1:创建请求失败！请检查网络连接！")
		return ""
	}

	// 正确响应
	if response.StatusCode != 200 {
		logrus.WithFields(
			logrus.Fields{
				"状态码":   response.StatusCode,
				"响应数据": response,
			},
		).Errorf("阶段1:解析响应数据失败")
		return ""
	}
	return partTwo()

}

// 获取Cf盾
func partTwo() string {
	link := "https://explorer.api.openai.com/api/auth/csrf"
	logrus.WithField("Link", link).Infof("阶段2:准备获取CF盾....")
	req, err := http.NewRequest("GET", link, nil)
	// Prepare Header
	req.Header.Set("Host", "ask.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://explorer.api.openai.com/auth/login")
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("阶段2:创建Http请求失败")
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && !strings.Contains(resp.Header.Get("Content-Type"), "json") {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段2:Htpp未能正确响应")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Errorf("阶段2:解析响应体失败")
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.WithError(err).Errorf("阶段2:发序列化body -> map失败")
		return ""
	}
	csrfToken, ok := data["csrfToken"].(string)
	if !ok {
		logrus.Errorf("cant find __Host-next-auth.csrf-token in headers")
		return ""
	}
	return partThree(csrfToken)
}

// https://explorer.api.openai.com/api/auth/signin/auth0?
func partThree(token string) string {
	link := "https://explorer.api.openai.com/api/auth/signin/auth0?"
	logrus.WithField("Link", link).Info("阶段3:开始进行模拟登陆....")
	payload := strings.NewReader(fmt.Sprintf("callbackUrl=%%2F&csrfToken=%s&json=true", token))
	req, err := http.NewRequest("POST", link, payload)
	if err != nil {
		logrus.WithError(err).Errorf("阶段1:创建Http请求失败")
		return ""
	}

	// Prepare Header
	req.Header.Add("Authority", "explorer.api.openai.com")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Add("Origin", "https://explorer.api.openai.com")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://explorer.api.openai.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		logrus.WithError(err).Errorf("阶段3:创建请求失败！请检查网络连接！")
		return ""
	}

	// 校验获取信息
	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode != 200 && !strings.Contains(contentType, "json") {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段2:Htpp可能未能正确响应")
		return ""
	}
	var jsonResponse struct {
		URL string `json:"url"`
	}
	data, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(data, &jsonResponse)
	if err != nil {
		logrus.WithError(err).Errorf("阶段3:解析响应体失败")
		return ""
	}

	// 判断是否被限制登陆
	if jsonResponse.URL == "https://explorer.api.openai.com/api/auth/error?error=OAuthSignin" || strings.Contains(jsonResponse.URL, "error") {
		logrus.WithField("地址", jsonResponse.URL).Errorf("阶段3:您已经被限制登陆")
		return ""
	}
	return partFore(jsonResponse.URL)
}

// https://auth0.openai.com/authorize?client_id
func partFore(link string) string {
	logrus.WithField("Link", link).Info("阶段4:获取登陆表单....")
	req, err := http.NewRequest("GET", link, nil)
	// Prepare Header
	req.Header.Set("Authority", "auth0.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段4:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}

	// 判断是否为成功重定向
	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段4:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""
	}
	// 匹配状态信息
	location := resp.Header.Get("Location")
	state := regexp.MustCompile(`state=(.*)`).FindStringSubmatch(location)[0]
	state = strings.Split(state, "=")[1]
	return partFive(state)
}

// https://auth0.openai.com/u/login/identifier?state
func partFive(state string) string {
	link := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	logrus.WithField("Link", link).Info("阶段5:登陆表单获取前回调....")

	req, err := http.NewRequest("GET", link, nil)
	req.Header.Set("Host", "auth0.openai.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://explorer.api.openai.com/")
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logrus.WithError(err).Errorf("阶段5:发送获取实际表单信息请求出错！")
		return ""
	}
	// 成功响应
	if resp.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段5:Http响应非200，请检查网络")
		return ""
	}
	return partSix(state)
}

// https://auth0.openai.com/u/login/identifier?state
func partSix(state string) string {
	link := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	logrus.WithField("Link", link).Info("阶段6:开始进行登陆..输入用户名中....")
	// 设置登陆参数
	usernameEncode := url.QueryEscape(userName)
	payload := fmt.Sprintf("state=%s&username=%s&action=default&js-available=false&webauthn-available=true&is-brave=false&webauthn-platform-available=true&action=default", state, usernameEncode)
	// 构建表单post请求 - 请求密码框
	req, err := http.NewRequest("POST", link, bytes.NewBufferString(payload))
	req.Header.Add("Authority", "explorer.api.openai.com")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段6:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}

	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段6:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""
	}
	return partSeven(state)
}

// https://auth0.openai.com/u/login/password?state
func partSeven(state string) string {
	link := fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state)
	logrus.WithField("Link", link).Info("阶段7:开始进行登陆..输入密码中....")
	// 设置请求参数
	usernameEncode := url.QueryEscape(userName)
	passwd := url.QueryEscape(password)

	payload := fmt.Sprintf("state=%s&username=%s&password=%s&action=default", state, usernameEncode, passwd)
	// 发起登陆请求
	req, err := http.NewRequest("POST", link, bytes.NewBufferString(payload))
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段7:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}
	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段4:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""

	}
	location := resp.Header.Get("Location")
	newState := regexp.MustCompile(`state=(.*)`).FindStringSubmatch(location)[0]
	newState = strings.Split(newState, "=")[1]
	return partEight(state, newState)
}

// https://auth0.openai.com/authorize/resume?state=%s
func partEight(oldState, newState string) string {
	link := fmt.Sprintf("https://auth0.openai.com/authorize/resume?state=%s", newState)
	logrus.WithField("Link", link).Info("阶段8:获取登陆表单....")
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", oldState))

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	// 不出错 且不是302错误
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段8:出现非自定义重定向一次限制错误，可能为http请求失败！")
		return ""
	}
	if resp.StatusCode != 302 {
		logrus.WithFields(logrus.Fields{
			"状态码":   resp.StatusCode,
			"响应信息": resp,
		}).Errorf("阶段8:Htpp可能不是302重定向，无法获取下一阶段登陆信息")
		return ""
	}
	// 获取回调地址
	return partNight(resp.Header.Get("Location"))
}

func partNight(urlCallBack string) string {
	logrus.WithField("Link", urlCallBack).Info("阶段9:进行获取token回调....")

	req, _ := http.NewRequest("GET", urlCallBack, nil)
	// Prepare Header
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Authority", "explorer.api.openai.com")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	// 最后一个请求有错误但是不是期待的重定向错误
	if err != nil && !errors.Is(err, &xconst.GPTREDIRECTERR) {
		logrus.WithError(err).Errorf("阶段9:发生非重定向错误")
		return ""
	}
	cookies := resp.Cookies()
	foundSessionToken := ""
	for _, cookie := range cookies {
		if cookie.Name == "__Secure-next-auth.session-token" {
			foundSessionToken = cookie.Value
			break
		}
	}
	if foundSessionToken == "" {
		logrus.WithField("Cookie", cookies).Errorf("阶段9:无法获取对应的Token")
		return ""
	}
	return getAkskFromReq(foundSessionToken)

}

func getAkskFromReq(sessionToken string) string {
	link := "https://explorer.api.openai.com/api/auth/session"
	logrus.WithField("Link", link).Info("最后阶段10:开始使用获取到的session-Token进行检查,并换取Token....")

	req, err := http.NewRequest("GET", link, nil)
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		logrus.WithError(err).WithField("上一步获取的Token", sessionToken).Errorf("阶段10:http请求失败！")
		return ""
	}

	// 检查是否成功
	if resp.StatusCode != http.StatusOK {
		logrus.WithField("response", resp).Errorf("阶段10:http请求响应状态码失败！")
		return ""
	}

	// 反序列化保存Token
	var data struct {
		AccessToken string `json:"accessToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logrus.WithError(err).Errorf("阶段10:反序列化解析AccessToken失败")
		return ""
	}

	logrus.WithField("LoginToken", data.AccessToken).Debugf("成功获取到Token")
	return data.AccessToken
}
