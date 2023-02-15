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
	"sort"
	"strings"
)

// CheckSign  (验证企业微信回调签名)[https://developer.work.weixin.qq.com/document/path/90968#%E6%B6%88%E6%81%AF%E4%BD%93%E7%AD%BE%E5%90%8D%E6%A0%A1%E9%AA%8C]
func CheckSign(token, signature, timestamp, nonce, msg_encrypt string) bool {
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

// MessageDecrypt (解密消息内容)[https://developer.work.weixin.qq.com/document/path/90968#%E5%AF%86%E6%96%87%E8%A7%A3%E5%AF%86%E5%BE%97%E5%88%B0msg%E7%9A%84%E8%BF%87%E7%A8%8B]
func MessageDecrypt(originStrBase64Encrypt, encodingAESKey string) ([]byte, error) {
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

// ParseContent 解析消息内容
func ParseContent(msgDecrypted []byte) ([]byte, uint32, []byte, []byte, error) {
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
