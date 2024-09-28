package server

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/jmbit/dtsrv/cmd/api"
	"github.com/jmbit/dtsrv/frontend"
	"github.com/jmbit/dtsrv/lib/reverseproxy"
)

func registerRoutes() http.Handler {

	mux := http.NewServeMux()
  //API
  mux.HandleFunc("/api/v0/ct/list", api.ListContainers)
  mux.HandleFunc("/api/v0/ct/new", api.NewContainer)
  mux.HandleFunc("GET /api/v0/ct/{ctName}/ready", api.ContainerReady)
  mux.HandleFunc("PUT /api/v0/ct/{ctName}/stop", api.StopContainer)
  mux.HandleFunc("POST /api/v0/admin/login", api.AdminLogin)
  mux.HandleFunc("POST /api/v0/admin/logout", api.AdminLogout)
	mux.HandleFunc("/view/{ctName}/", reverseproxy.HandleReverseProxy)
	if _, ok := os.LookupEnv("BLOCK_FILEBROWSER"); ok == true {
		mux.HandleFunc("/view/{ctName}/files", reverseproxy.HandleUnauthorized)
	}

  dist, err := fs.Sub(frontend.Files, "dist")
  if err != nil {
    panic(err)
  }
	fileServer := http.FileServer(http.FS(dist))
	mux.Handle("/", fileServer)

	return mux
}
