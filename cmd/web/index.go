package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmbit/dtsrv/internal/session"
	"github.com/jmbit/dtsrv/lib/containers"
	"github.com/jmbit/dtsrv/lib/reverseproxy"
	"github.com/spf13/viper"

	"github.com/docker/docker/api/types"
)

// IndexWebHandler() is the main http handler for the application
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

	// Info for Table of containers
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

// StartWebHandler() starts a container
func StartWebHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionStore.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctName, err := containers.CreateContainer(viper.GetString("container.image"))
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
  port := viper.GetInt("container.port")
	ctUrl, err := containers.GetContainerUrl(ctName, &port)
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

// StartStatusWebHandler() is called pereodically after starting a container
func StartStatusWebHandler(w http.ResponseWriter, r *http.Request) {
	ctName := r.PathValue("ctName")
  port := viper.GetInt("container.port")
	component := StartSpinner(ctName)
	running, err := containers.TestConnectionToContainer(ctName, &port)
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

// indexQueryHandler() handles/works on any url parameters specified
func indexQueryHandler(w http.ResponseWriter, r *http.Request) error {
	sess, err := session.SessionStore.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// This can likely be cleaned up and less nested
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
