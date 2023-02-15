package main

import (
	"net/http"

	"person-bot/handler"
)

type server struct {
	mux *http.ServeMux
}

func (s server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func newServer() *server {
	s := &server{
		mux: http.NewServeMux(),
	}
	var routes = map[string]func(http.ResponseWriter, *http.Request){
		"/test": handler.Testhandler,
		"/":     handler.WeChatConfirmHandler,
	}
	for path, routers := range routes {
		s.mux.HandleFunc(path, routers)
	}
	return s
}
