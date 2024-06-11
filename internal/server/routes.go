package server

import (
	"encoding/json"
	"log"
	"net/http"

	"dtsrv/cmd/web"
	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/web", templ.Handler(web.HelloForm()))
	mux.HandleFunc("/hello", web.HelloWebHandler)

	return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
