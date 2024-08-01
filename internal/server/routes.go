package server

import (
	"net/http"

	"dtsrv/cmd/web"
	"dtsrv/internal/reverseproxy"
)

func registerRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.IndexWebHandler)
	mux.HandleFunc("/admin", web.AdminWebHandler)
  mux.HandleFunc("POST /start", web.StartWebHandler)
  mux.HandleFunc("GET /status/{ctName}", web.StartStatusWebHandler)
  mux.HandleFunc("/view/{ctName}/", reverseproxy.HandleReverseProxy)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	return mux
}

