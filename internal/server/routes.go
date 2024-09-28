package server

import (
	"net/http"
	"os"

	"github.com/jmbit/dtsrv/cmd/web"
	"github.com/jmbit/dtsrv/lib/reverseproxy"
)

func registerRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.IndexWebHandler)
	mux.HandleFunc("/admin", web.AdminWebHandler)
	mux.HandleFunc("POST /start", web.StartWebHandler)
	mux.HandleFunc("GET /status/{ctName}", web.StartStatusWebHandler)
	mux.HandleFunc("/view/{ctName}/", reverseproxy.HandleReverseProxy)
	if _, ok := os.LookupEnv("BLOCK_FILEBROWSER"); ok == true {
		mux.HandleFunc("/view/{ctName}/files", reverseproxy.HandleUnauthorized)
	}

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	return mux
}
