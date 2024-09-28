package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmbit/dtsrv/internal/admin"
	"github.com/jmbit/dtsrv/internal/session"
	"github.com/jmbit/dtsrv/lib/containers"
	"github.com/jmbit/dtsrv/lib/reverseproxy"
)

// ContainerData gets populated as needed
type ContainerData struct {
  Name string `json:",omitempty"`
  Image string `json:",omitempty"`
  State string `json:",omitempty"`
  Ready bool `json:",omitempty"`
  Created time.Time `json:",omitempty"`
}

// ContainerStart() starts a new container
func ContainerStart(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionStore.Get(r, "session")
	if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
		return
	}
	ctName, err := containers.CreateContainer()
	if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
		return
	}
	err = session.AppendContainer(sess, ctName)
	if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
		return
	}
	err = sess.Save(r, w)
	if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
		return
	}
	ctUrl, err := containers.GetContainerUrl(ctName)
	if err != nil {
		log.Println("Error parsing container url,", err)
    JsonError(w, r, err, http.StatusBadRequest)
	}
	go reverseproxy.NewContainerProxy(ctName, ctUrl)

	if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
    return
	}
  ctData := ContainerData{Name: ctName}
  returnContainerData(w, r, ctData)
}

// ContainerReady() checks if a container is ready to accept the incoming session
func ContainerReady(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctName")
	ready, err := containers.TestConnectionToContainer(ctName)
  if err != nil {
    log.Println("Error in ContainerReady API function: ", err)
    JsonError(w, r, err, http.StatusInternalServerError)
  }
  ctData := ContainerData{Name: ctName, Ready: ready}
  returnContainerData(w, r, ctData)
}

// ListContainers() returns a list of containers, for admins it's all containers,
// for non-admin users only the users containers
func ListContainers(w http.ResponseWriter, r *http.Request) {
  var retCtList []ContainerData
  isAdmin, err := admin.IsAdmin(w, r)
  if err != nil {
    log.Println("Error in ListContainers API function: ", err)
    JsonError(w, r, err, http.StatusInternalServerError)
    return
  }
  rawCtList, err := containers.ListContainers()
  if err != nil {
    log.Println("Error in ListContainers API function: ", err)
    JsonError(w, r, err, http.StatusInternalServerError)
    return
  }
  if isAdmin {
    for _, raw := range rawCtList {
      ct := ContainerData{
        Name: raw.Names[0],
        Image: raw.Image,
        State: raw.State,

      }
      retCtList = append(retCtList, ct)
    }
    
  }

}


// returnJson() marshals the struct to json and writes it to the response
func returnContainerData(w http.ResponseWriter, r *http.Request, data any) error {
  w.Header().Add("content-type", "application/json")
  jsonData, err := json.MarshalIndent(data, "", " ") 
  if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
    return err
  }
  _, err = w.Write(jsonData)
  if err != nil {
    log.Println(err)
    panic(1)
  }

  return nil
}
