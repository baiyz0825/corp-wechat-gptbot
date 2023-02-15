// Package config 配置包
package config

import (
	"fmt"

	"github.com/spf13/viper"
	constx "person-bot/const"
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
	Proxy string `json:"proxy,omitempty" yaml:"proxy"`
	Port  string `json:"port" yaml:"port"`
	Log   string `json:"log" yaml:"log"`
}

// GptConfig chatGpt api key
type GptConfig struct {
	Apikey string `json:"apikey" yaml:"apikey"`
}

// WeChatConfig 微信配置文件
type WeChatConfig struct {
	//  Corpid 企业ID
	Corpid string `json:"corpid" yaml:"corpid"`
	//  CorpSecret 企业应用Secret
	CorpSecret string `json:"corpSecret" yaml:"corpSecret"`
	//  AgentId 应用ID
	AgentId string `json:"agentId" yaml:"agentId"`
	//  WeApiRCallToken 企业微信消息Token
	WeApiRCallToken string `json:"weApiRCallToken" yaml:"weApiRCallToken"`
	//  WeApiEncodingKey 企业微信消息Key
	WeApiEncodingKey string `json:"weApiEncodingKey" yaml:"weApiEncodingKey"`
	// 企业微信API地址
	WeChatApiAddr string `json:"weChatApiAddr" yaml:"weChatApiAddr"`
}

// LoadConf  加载配置文件
func LoadConf() error {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf(constx.InitS+"load config file failure , please check you config file:%w", err)
	}
	if err := v.Unmarshal(&globalConf); err != nil {
		return fmt.Errorf(constx.InitS+"load config file failure , please check you config file:%w", err)
	}
	return nil
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
