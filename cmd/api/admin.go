package api

import (
	"net/http"

	"github.com/jmbit/dtsrv/internal/admin"
)


func AdminContainerList(w http.ResponseWriter, r *http.Request) {
  admin, err := admin.IsAdmin(w, r)
  if err != nil {
    JsonError(w, r, err, http.StatusInternalServerError)
  }
  if admin == false {
    return
  }
  
}

