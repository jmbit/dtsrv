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

// NewContainer() creates and starts a new container
func NewContainer(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionStore.Get(r, "session")
	if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
		return
	}
	ctName, err := containers.CreateContainer("")
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
	ctUrl, err := containers.GetContainerUrl(ctName, nil)
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
  returnJson(w, r, ctData)
}

// ContainerReady() checks if a container is ready to accept the incoming session
func ContainerReady(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctName")
	ready, err := containers.TestConnectionToContainer(ctName, nil)
  if err != nil {
    log.Println("Error in ContainerReady API function: ", err)
    JsonError(w, r, err, http.StatusInternalServerError)
  }
  ctData := ContainerData{Name: ctName, Ready: ready}
  returnJson(w, r, ctData)
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
  } else {
	  sess, err := session.SessionStore.Get(r, "session")
	  if err != nil {
      log.Println("Error in ListContainers API function: ", err)
      JsonError(w, r, err, http.StatusInternalServerError)
	  	return
	  }
    for _, raw := range rawCtList {
      owned, err := session.OwnsContainer(sess, raw.Names[0])
  	  if err != nil {
        log.Println("Error in ListContainers API function: ", err)
        continue
  	  }
      if owned == true {
        ct := ContainerData{
         Name: raw.Names[0],
         Image: raw.Image,
         State: raw.State,
       }
       retCtList = append(retCtList, ct)       
      }
    }
  }
  returnJson(w, r, rawCtList)
}

// StopContainer() stops a container
func StopContainer(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctName")
  sess, err := session.SessionStore.Get(r, "session")
  if err != nil {
    log.Println("Error in StopContainer API function: ", err)
    JsonError(w, r, err, http.StatusInternalServerError)
  	return
  }
  isAdmin, err := session.IsAdmin(sess)
  if err != nil {
    log.Println("Error in StopContainer API function: ", err)
    JsonError(w, r, err, http.StatusUnauthorized)
  	return
  }
  owned, err := session.OwnsContainer(sess, ctName)
  if err != nil {
    log.Println("Error in StopContainer API function: ", err)
    JsonError(w, r, err, http.StatusUnauthorized)
  	return
  }
  if isAdmin || owned {
    err := containers.StopContainer(ctName)
    if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
  } else {
    reverseproxy.HandleUnauthorized(w, r)
  }
}

// DeleteContainer() deletes a container
func DeleteContainer(w http.ResponseWriter, r *http.Request) {
  ctName := r.PathValue("ctName")
  sess, err := session.SessionStore.Get(r, "session")
  if err != nil {
    log.Println("Error in DeleteContainer API function: ", err)
    JsonError(w, r, err, http.StatusInternalServerError)
  	return
  }
  isAdmin, err := session.IsAdmin(sess)
  if err != nil {
    log.Println("Error in DeleteContainer API function: ", err)
    JsonError(w, r, err, http.StatusUnauthorized)
  	return
  }
  owned, err := session.OwnsContainer(sess, ctName)
  if err != nil {
    log.Println("Error in DeleteContainer API function: ", err)
    JsonError(w, r, err, http.StatusUnauthorized)
  	return
  }
  if isAdmin || owned {
    err := containers.DeleteContainer(ctName)
    if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
    }
		// Delete proxy to prevent leaking memory
		reverseproxy.DeleteContainerProxy(ctName)

    w.WriteHeader(http.StatusOK)
  } else {
    reverseproxy.HandleUnauthorized(w, r)
  }
}


// returnJson() marshals the struct to json and writes it to the response
func returnJson(w http.ResponseWriter, r *http.Request, data any) error {
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
