package main

import (
	"net/http"

	"person-bot/handler"
	"person-bot/logic"
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
	// 注册handler && strategy
	var routes = map[string]http.Handler{
		"/wx": &handler.WxHandler{
			WxStrategy: map[string]handler.WxBusStrategy{
				"chat": &logic.WxGptLogic{},
				"cmd":  &logic.WxJokeLogic{},
			},
		},
	}
	for path, routers := range routes {
		s.mux.Handle(path, routers)
	}
	return s
}
