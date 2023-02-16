package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"person-bot/config"
)

// checkSign  (验证企业微信回调签名)[https://developer.work.weixin.qq.com/document/path/90968#%E6%B6%88%E6%81%AF%E4%BD%93%E7%AD%BE%E5%90%8D%E6%A0%A1%E9%AA%8C]
func checkSign(token, signature, timestamp, nonce, msg_encrypt string) bool {
	sortStr := []string{token, timestamp, nonce, msg_encrypt}
	sort.Strings(sortStr)
	var buffer bytes.Buffer
	for _, value := range sortStr {
		buffer.WriteString(value)
	}
	hash := sha1.New()
	hash.Write([]byte(buffer.Bytes()))
	encodeStr := hex.EncodeToString(hash.Sum(nil))
	return strings.EqualFold(signature, encodeStr)
}

// messageDecrypt (解密消息内容)[https://developer.work.weixin.qq.com/document/path/90968#%E5%AF%86%E6%96%87%E8%A7%A3%E5%AF%86%E5%BE%97%E5%88%B0msg%E7%9A%84%E8%BF%87%E7%A8%8B]
func messageDecrypt(originStrBase64Encrypt, encodingAESKey string) ([]byte, error) {
	AESKey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return nil, err
	}
	encryptStr, err := base64.StdEncoding.DecodeString(originStrBase64Encrypt)
	if err != nil {
		return nil, err
	}
	// 根据 key 创建 AES 解密器
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return nil, err
	}
	// confirm size
	if len(encryptStr) < aes.BlockSize {
		return nil, errors.New("encrypt_msg size is not valid")
	}
	if len(encryptStr)%aes.BlockSize != 0 {
		return nil, errors.New("encrypt_msg not a multiple of the block size")
	}
	// 准备解密器
	iv := AESKey[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	// 解密
	mode.CryptBlocks(encryptStr, encryptStr)
	return encryptStr, nil
}

// pKCS7Unpadding 将pKCS7填充的字节去掉，还原出原始数据
// 加密算法使用的分组大小：blockSize 需要解填充的数据（PKCS7 填充的数据）：plaintext
func pKCS7Unpadding(plaintext []byte, blockSize int) ([]byte, error) {
	plaintextLen := len(plaintext)
	// null test
	if nil == plaintext || plaintextLen == 0 {
		return nil, errors.New("pKCS7Unpadding error nil or zero")
	}
	// if match block size
	if plaintextLen%blockSize != 0 {
		return nil, errors.New("pKCS7Unpadding text not a multiple of the block size")
	}
	// 获取填充的字节长度 paddingLen
	paddingLen := int(plaintext[plaintextLen-1])
	if paddingLen > blockSize {
		return nil, errors.New("the text not match PKCS7 rules")
	}
	// 删除末尾的 paddingLen 个字节删除 得到还原后的数据
	return plaintext[:plaintextLen-paddingLen], nil
}

// ParseWxMsgToOriginContent 解析解密后的消息内容
func ParseWxMsgToOriginContent(msgDecrypted []byte) ([]byte, uint32, []byte, []byte, error) {
	// 块大小
	const blockSize = 32
	// 还原PKCS7填充
	plaintext, err := pKCS7Unpadding(msgDecrypted, blockSize)
	if nil != err {
		return nil, 0, nil, nil, err
	}
	// 计算文本长度
	textLen := uint32(len(plaintext))
	// 16个字节的随机字符串、4个字节的msg长度、明文msg和receiveId拼接
	if textLen < 20 {
		return nil, 0, nil, nil, errors.New("plain is to small 1")
	}
	random := plaintext[:16]
	// 网络协议也都是采用big endian的方式来传输数据的 内存的低地址存放着数据高位
	msgLen := binary.BigEndian.Uint32(plaintext[16:20])
	if textLen < (20 + msgLen) {
		return nil, 0, nil, nil, errors.New("plain is to small 2")
	}

	msg := plaintext[20 : 20+msgLen]
	receiverId := plaintext[20+msgLen:]

	return random, msgLen, msg, receiverId, nil
}

// DecodeAndDecryptWxMsg 解密获取微信回掉中的url消息内容（含有编码）
func decodeAndDecryptUrl(urlParse url.URL) (error, []byte) {
	// 获取url请求参数
	param := urlParse.Query()
	signature := param.Get("msg_signature")
	timestamp := param.Get("timestamp")
	nonce := param.Get("nonce")
	msgEncrypt := param.Get("echostr")
	if !checkSign(config.GetWechatConf().WeApiRCallToken, signature, timestamp, nonce, msgEncrypt) {
		// 验证失败
		log.Error("企业微信签名验证失败")
		return errors.New("企业微信签名验证失败"), nil
	}
	// 解密消息内容
	decrypt, err := messageDecrypt(msgEncrypt, config.GetWechatConf().WeApiEncodingKey)
	if err != nil {
		// 解密失败
		log.Error("解密url回调消息内容失败：%w", err)
		return errors.New("解密url回调消息内容失败"), nil
	}
	return nil, decrypt
}

// CheckUrlFromWeChat CheckUrlAndBodyFromWeChat 检查url有效性 && 返回解析后的真实消息内容
func CheckUrlFromWeChat(urlParse url.URL) []byte {
	err, decrypt := decodeAndDecryptUrl(urlParse)
	if err != nil {
		log.Error("CheckUrl: 解码 && 解密回调消息内容失败：%w", err)
		return nil
	}
	// 解码还原消息内容
	_, _, content, receiverId, err := ParseWxMsgToOriginContent(decrypt)
	if err != nil {
		log.Error("CheckUrl: 解码回调消息内容失败：%w", err)
		return nil
	}
	if !strings.EqualFold(string(receiverId), config.GetWechatConf().Corpid) {
		// 不是发给此企业
		log.Error("CheckUrl: 回调企业错误")
		return nil
	}
	return content
}

// CheckAndParseBody CheckUrl && GetPost Data 检查URL有效性并且获取响应体数据
func CheckAndParseBody(urlParse url.URL, messageEncrypt string) []byte {
	// 检查URL
	// 获取url请求参数
	param := urlParse.Query()
	signature := param.Get("msg_signature")
	timestamp := param.Get("timestamp")
	nonce := param.Get("nonce")
	if !checkSign(config.GetWechatConf().WeApiRCallToken, signature, timestamp, nonce, messageEncrypt) {
		// 验证失败
		log.Error("解析回掉POST请求: 企业微信签名验证失败")
		return nil
	}
	// 解密消息内容
	decrypt, err := messageDecrypt(messageEncrypt, config.GetWechatConf().WeApiEncodingKey)
	if err != nil {
		// 解密失败
		log.Error("解析回掉POST请求: 解密url回调消息内容失败：%w", err)
		return nil
	}
	// 还原数据
	_, _, msg, receiverId, err := ParseWxMsgToOriginContent(decrypt)
	if err != nil {
		log.Error("解析回掉POST请求: 还原解密后的数据失败：%w", err)
		return nil
	}
	if !strings.EqualFold(string(receiverId), config.GetWechatConf().Corpid) {
		// 不是发给此企业
		log.Error("解析回掉POST请求: 消息内容校验失败，receiverId与Corpid不相符")
		return nil
	}
	return msg
}

// ParseCommandFromStr 从内容中获取用户指令
func ParseCommandFromStr(str string) string {
	return str
}
