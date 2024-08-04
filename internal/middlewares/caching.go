package middlewares

import (
	"net/http"
  "strings"

)

func AssetCaching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.EscapedPath(), "/assets/") {
      w.Header().Set("Cache-Control", "max-age=3600")
    }

		next.ServeHTTP(w, r)

	})
}
