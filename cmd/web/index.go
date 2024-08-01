package web

import (

	"dtsrv/internal/reverseproxy"
	"dtsrv/internal/containers"
	"fmt"
	"log"
	"net/http"
)

func IndexWebHandler(w http.ResponseWriter, r *http.Request) {
	component := Index()
  err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in IndexWebHandler: %e", err)
	}
}

func StartWebHandler(w http.ResponseWriter, r *http.Request) {
  ctName, err := containers.CreateContainer()
  ctUrl, err := containers.GetContainerUrl(ctName)
  if err != nil {
    log.Println("Error parsing container url,", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
  }
  go reverseproxy.NewContainerProxy(ctName, ctUrl)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
	component := Start(ctName)
  err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in IndexWebHandler: %e", err)
	}
}

func StartStatusWebHandler(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctName")
  component := StartSpinner(ctName)
  running, err := containers.TestConnectionToContainer(ctName)
  if running == true {
    w.Header().Add("HX-Redirect", fmt.Sprintf("/view/%s/", ctName))
  }
  if err != nil {
    log.Println(err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in IndexWebHandler: %e", err)
	}
}

