// Package config 配置包
package config

import (
	"time"

	"github.com/spf13/viper"
)

var globalConf *GlobalConf

// GlobalConf 全局配置文件
type GlobalConf struct {
	// SystemConf 系统配置
	SystemConf SystemConf `json:"systemConf" yaml:"systemConf"`
	// GptApiKey 密钥
	GptConfig GptConfig `json:"gptConfig" yaml:"gptConfig"`
	// WeConfig 微信配置
	WeConfig WeChatConfig `json:"weConfig" yaml:"weConfig"`
}

// SystemConf 系统配置
type SystemConf struct {
	// Proxy 代理地址
	Proxy       string     `json:"proxy,omitempty" yaml:"proxy"`
	Port        string     `json:"port" yaml:"port"`
	CallBackUrl string     `json:"callBackUrl"`
	MsgMode     string     `json:"msgMode"`
	LogConf     LoggerConf `json:"logConf" yaml:"logConf"`
}

type LoggerConf struct {
	LogLevel             string `json:"logLevel,omitempty" yaml:"logLevel"`
	LogOutPutMode        string `json:"logOutPutMode" yaml:"logOutPutMode"`
	LogOutPutPath        string `json:"logOutPutPath,omitempty" yaml:"logOutPutPath"`
	LogFileDateFmt       string `json:"logFileDateFmt,omitempty" yaml:"logFileDateFmt"`
	LogFileMaxSizeM      int64  `json:"logFileMaxSizeM,omitempty" yaml:"logFileMaxSizeM"`
	LogFileRotationCount uint   `json:"logFileRotationCount,omitempty" yaml:"logFileRotationCount"`
	LogFormatter         string `json:"logFormatter,omitempty" yaml:"logFormatter"`
}

// GptConfig chatGpt api key
type GptConfig struct {
	Apikey            string        `json:"apikey" yaml:"apikey"`
	Model             string        `json:"model" yaml:"model"`
	UserName          string        `json:"UserName" yaml:"UserName"`
	URL               string        `json:"url" yaml:"url"`
	ContextNumber     int           `json:"contextNumber" yaml:"contextNumber"`
	ContextExpireTime time.Duration `json:"contextExpireTime" yaml:"contextExpireTime"`
}

// WeChatConfig 微信配置文件
type WeChatConfig struct {
	//  Corpid 企业ID
	Corpid string `json:"corpid" yaml:"corpid"`
	//  CorpSecret 企业应用Secret
	CorpSecret string `json:"corpSecret" yaml:"corpSecret"`
	//  AgentId 应用ID
	AgentId int `json:"agentId" yaml:"agentId"`
	//  WeApiRCallToken 企业微信消息Token
	WeApiRCallToken string `json:"weApiRCallToken" yaml:"weApiRCallToken"`
	//  WeApiEncodingKey 企业微信消息Key
	WeApiEncodingKey string `json:"weApiEncodingKey" yaml:"weApiEncodingKey"`
	// 企业微信API地址
	WeChatApiAddr string `json:"weChatApiAddr" yaml:"weChatApiAddr"`
}

// LoadConf  加载配置文件
func init() {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.SetConfigName("config")
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&globalConf); err != nil {
		panic(err)
	}
}

func GetWechatConf() *WeChatConfig {
	if globalConf == nil {
		return nil
	}
	return &globalConf.WeConfig
}

func GetGptConf() *GptConfig {
	if globalConf == nil {
		return nil
	}
	return &globalConf.GptConfig
}

func GetSystemConf() *SystemConf {
	if globalConf == nil {
		return nil
	}
	return &globalConf.SystemConf
}
