package server

import (
	"net/http"

	"github.com/jmbit/dtsrv/cmd/api"
	"github.com/jmbit/dtsrv/cmd/web"
	"github.com/jmbit/dtsrv/lib/reverseproxy"
	"github.com/spf13/viper"
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

  //Web UI
  mux.HandleFunc("/", web.IndexWebHandler)
	mux.HandleFunc("/admin", web.AdminWebHandler)
	mux.HandleFunc("POST /start", web.StartWebHandler)
	mux.HandleFunc("GET /status/{ctName}", web.StartStatusWebHandler)



	mux.HandleFunc("/view/{ctName}/", reverseproxy.HandleReverseProxy)
	if viper.GetBool("web.blockfilebrowser") {
		mux.HandleFunc("/view/{ctName}/files", reverseproxy.HandleUnauthorized)
	}

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	return mux
}
