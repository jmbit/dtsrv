package web

import (
	"dtsrv/internal/containers"
	"dtsrv/internal/session"
	"fmt"
	"log"
  "os"

	"net/http"
)

func AdminWebHandler(w http.ResponseWriter, r *http.Request) {

  admin, err := isAdmin(w, r) 
  if err != nil {
    return 
  }
  if admin == false {
    authorized, err := adminLoginHandler(w, r)
    if err != nil || authorized == false {
      return
    }
  }
  

  if r.URL.Query().Get("action") != "" {
    adminQueryHandler(w, r)
    return
  }

  ctList, err := containers.ListContainers()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  imageName := containers.GetImageName()

	component := Admin(imageName, ctList)
  err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in AdminWebHandler: %e", err)
    return
	}

}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) (bool,error) {
if r.Method == http.MethodPost {
    sess, err := session.SessionStore.Get(r, "session")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return false, err
    }
 
    err = r.ParseForm()
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return false, err
      }
    if r.PostFormValue("password") == os.Getenv("ADMIN_PW") {
        sess.Values["admin"] = true
        err := sess.Save(r, w)
        if err != nil {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return false, err
        }
        return true, nil
      }
    } else {
    component := AdminLogin()
    err := component.Render(r.Context(), w)
	  if err != nil {
	  	http.Error(w, err.Error(), http.StatusBadRequest)
	  	log.Fatalf("Error rendering in AdminWebHandler: %e", err)
      return false, err
	  }
    return false, nil
   }
    component := AdminLogin()
    err := component.Render(r.Context(), w)
	  if err != nil {
	  	http.Error(w, err.Error(), http.StatusBadRequest)
	  	log.Fatalf("Error rendering in AdminWebHandler: %e", err)
      return false, err
	  }
    return false, nil
}

func adminQueryHandler(w http.ResponseWriter, r *http.Request) error {
  if r.URL.Query().Get("action") != "" {
    if r.URL.Query().Get("action") == "logout" {
      sess, err := session.SessionStore.Get(r, "session")
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
      }
      sess.Values["admin"] = false
      err = sess.Save(r, w)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
      }
      http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)

      return nil
    }

    ctName := r.URL.Query().Get("ctName")
    if ctName == "" {
	  	http.Error(w, "No container specified", http.StatusBadRequest)
      return fmt.Errorf("No container specified")
    }
 
    if r.URL.Query().Get("action") == "stop" {
      err := containers.StopContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Println("Error stopping container", ctName, err)
        return err
      }

      ctInfo, err := containers.GetContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
          log.Println("Error deleting container", ctName, err)
        return err
      } 
      component := adminContainerRow(ctInfo)
      ctInfo.State = "exited" 
      err = component.Render(r.Context(), w)
      if err != nil {
      	http.Error(w, err.Error(), http.StatusBadRequest)
      	log.Fatalf("Error rendering in AdminWebHandler: %e", err)
        return err
      }
      return nil
    } else if r.URL.Query().Get("action") == "delete" {
      err := containers.DeleteContainer(ctName)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
      } else {
        w.WriteHeader(http.StatusOK)
        return nil
      }
    }
  }
  return nil
}

func isAdmin(w http.ResponseWriter, r *http.Request) (bool, error) {
  sess, err := session.SessionStore.Get(r, "session")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return false, err
  }
  if isAdmin, ok := sess.Values["admin"].(bool); ok && isAdmin {
    return true, nil
  }

  return false, nil
}
