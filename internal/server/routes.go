package server

import (
	"net/http"

	"dtsrv/cmd/web"
)

func registerRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.IndexWebHandler)
  mux.HandleFunc("POST /start", web.StartWebHandler)
  mux.HandleFunc("GET /ct/{ctid}/status", handleReverseProxy)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	return mux
}

