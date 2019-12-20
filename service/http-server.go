package service

import (
	"fmt"
	"net/http"
)

type HTTPServer struct {
	*ServerBase
}

func CreateHTTPServer(base *ServerBase) *HTTPServer {
	server := &HTTPServer{
		base,
	}
	return server
}

func (server *HTTPServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {

}

func (server *HTTPServer) Run() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", *server.port), server)
	if err != nil {
		return err
	}

	return nil
}