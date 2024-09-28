package api

import (
	"log"
	"net/http"
	"github.com/jmbit/dtsrv/internal/admin"
	"github.com/jmbit/dtsrv/internal/session"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
  authorized, err := admin.Login(w, r)
  if err != nil {
    log.Println("Error in Admin login handler: ", err)
			JsonError(w, r, err, http.StatusInternalServerError)
    return
  }
  if authorized {
    w.WriteHeader(http.StatusOK)
  } else {
    w.WriteHeader(http.StatusUnauthorized)
  }
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
			sess, err := session.SessionStore.Get(r, "session")
			if err != nil {
        log.Println("Error in Admin logout handler: ", err)
				JsonError(w, r, err, http.StatusInternalServerError)
        return
			}
			sess.Values["admin"] = false
			err = sess.Save(r, w)
			if err != nil {
        log.Println("Error in Admin logout handler: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
}
