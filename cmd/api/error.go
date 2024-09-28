package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonError(w http.ResponseWriter, r *http.Request, content error, code int) {
  errorJson, err := json.MarshalIndent(content, "", " ")
  if err != nil {
    log.Fatalf("error unmarshaling error %v: %v", content, err)
  }
  w.Header().Set("content-type", "application/json")
  http.Error(w, string(errorJson), code)
} 
