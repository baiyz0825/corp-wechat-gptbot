package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"person-bot/config"
)

func init() {
	// load conf
	if err := config.LoadConf(); err != nil {
		panic(err)
	}
	// config log
	level, err := log.ParseLevel(config.GetSystemConf().Log)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level) // 设置输出警告级别
	log.SetOutput(os.Stdout)
	log.Info("Init config success")
}

func main() {
	// check proxy
	// if len(config.GetSystemConf().Proxy) != 0 {
	// 	parse, err := url.Parse(config.GetSystemConf().Proxy)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// start server
	log.Info("start server in : " + config.GetSystemConf().Port)
	err := http.ListenAndServe(":"+config.GetSystemConf().Port, newServer())
	if err != nil {
		log.Fatal(err)
	}

}
