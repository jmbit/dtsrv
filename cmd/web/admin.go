package web

import (
	"dtsrv/internal/containers"
	"log"
	"os"

	"net/http"
)

func AdminWebHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Query().Get("action") != "" {
    ctName := r.URL.Query().Get("ctName")
    if ctName == "" {
	  	http.Error(w, "No container specified", http.StatusBadRequest)
      return
    }
 
    if r.URL.Query().Get("action") == "stop" {
      err := containers.StopContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println("Error stopping container", ctName, err)
        return
      }

      ctInfo, err := containers.GetContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
          log.Println("Error deleting container", ctName, err)
        return
      } 
      component := containerRow(ctInfo)
      err = component.Render(r.Context(), w)
      if err != nil {
      	http.Error(w, err.Error(), http.StatusBadRequest)
      	log.Fatalf("Error rendering in AdminWebHandler: %e", err)
        return
      }
      return
    } else if r.URL.Query().Get("action") == "delete" {
      err := containers.DeleteContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      } else {
        w.WriteHeader(http.StatusOK)
        return
      }
    }
  }

  // Get default image
  var imageName string
  if envImage, exists := os.LookupEnv("IMAGE_NAME"); exists == true {
    imageName = envImage
  } else {
	imageName = "lscr.io/linuxserver/webtop"
  }
  ctList, err := containers.ListContainers()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

	component := Admin(imageName, ctList)
  err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in AdminWebHandler: %e", err)
    return
	}

}

