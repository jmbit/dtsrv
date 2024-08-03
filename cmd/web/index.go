package web

import (
	"dtsrv/internal/containers"
	"dtsrv/internal/reverseproxy"
	"dtsrv/internal/session"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
)

func IndexWebHandler(w http.ResponseWriter, r *http.Request) {
  sess, err := session.SessionStore.Get(r, "session")
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = indexQueryHandler(w, r)
  if err != nil {
    return
  }

  containerNameList, err := session.GetContainers(sess)
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  log.Println(containerNameList)
  containerRawList, err := containers.ListContainers()
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  containerList := []types.Container{}
  for _, container := range containerRawList {
    for _, containerName := range containerNameList {
      if container.Names[0] == fmt.Sprintf("/%s", containerName) {
        containerList = append(containerList, container)
      }
    }
  }

	component := Index(containerList)
  err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in IndexWebHandler: %e", err)
	}
}

func StartWebHandler(w http.ResponseWriter, r *http.Request) {
  sess, err := session.SessionStore.Get(r, "session")
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  ctName, err := containers.CreateContainer()
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = session.AppendContainer(sess, ctName)
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = sess.Save(r, w)
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
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



func indexQueryHandler(w http.ResponseWriter, r *http.Request) error {
  sess, err := session.SessionStore.Get(r, "session")
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    return err
  }

  if r.URL.Query().Get("action") != "" {
    ctName := r.URL.Query().Get("ctName")
    if ctName == "" {
	  	http.Error(w, "No container specified", http.StatusBadRequest)
      return fmt.Errorf("No container specified")
    }
    if ok, err := session.OwnsContainer(sess, ctName); ok == false {
      if err == nil {
        err = fmt.Errorf("Container not owned by this session")
      }
	  	http.Error(w, err.Error(), http.StatusBadRequest)
      return err
    }
 
    if r.URL.Query().Get("action") == "stop" {
      err := containers.StopContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println("Error stopping container", ctName, err)
        return err
      }
      session.RemoveContainer(sess, ctName)
      if err != nil {
    		http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
      }

      ctInfo, err := containers.GetContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
          log.Println("Error deleting container", ctName, err)
        return err
      } 
      component := indexContainerRow(ctInfo)
      err = component.Render(r.Context(), w)
      if err != nil {
      	http.Error(w, err.Error(), http.StatusBadRequest)
      	log.Fatalf("Error rendering in AdminWebHandler: %e", err)
        return err
      }
      return nil
    }
  }
  return nil
}
