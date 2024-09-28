package admin

import (
	"net/http"

	"github.com/jmbit/dtsrv/internal/session"
)

// IsAdmin() checks if the session is logged in as Admin
func IsAdmin(w http.ResponseWriter, r *http.Request) (bool, error) {
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
