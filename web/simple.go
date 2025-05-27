package web

import (
	"net/http"

	"github.com/jmbit/dtsrv/internal/session"
	"github.com/jmbit/dtsrv/lib/containers"
	"github.com/jmbit/dtsrv/lib/reverseproxy"
	"github.com/spf13/viper"
)

// HandleSimple() is used for "simple Mode", a very stripped down
// mode of operation
// TODO
func HandleSimple(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionStore.Get(r, "dtsrv-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ready, _ := sess.Values["simpleready"].(bool)
	if ready {
		reverseproxy.HandleReverseProxy(w, r)
	} else {
		if session.GetSimpleContainer(sess) == "" {
			containers.CreateContainer(viper.GetString("container.image"), viper.GetBool("container.isolated"))
		}
	}
}
