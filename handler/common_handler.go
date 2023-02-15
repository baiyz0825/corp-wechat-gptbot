package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Testhandler(w http.ResponseWriter, r *http.Request) {
	log.Info("测试测试----")
	w.Write([]byte("欢迎，测试"))
}
