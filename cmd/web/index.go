package web

import (
	"dtsrv/internal/containers"
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
  ctid := containers.CreateContainer()
	component := Start(ctid)
  err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in IndexWebHandler: %e", err)
	}
}

func StartStatusWebHandler(w http.ResponseWriter, r *http.Request) {
  ctid := r.PathValue("ctid")
  component := StartSpinner(ctid)
  running, err := containers.TestConnectionToContainer(ctid)
  if running == true {
    w.Header().Add("HX-Trigger", "done")
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

func RedirectToContainer(w http.ResponseWriter, r *http.Request)
