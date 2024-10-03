package admin

import (
	"log"
	"net/http"

	"github.com/jmbit/dtsrv/internal/session"
	"github.com/spf13/viper"
)

func Login(w http.ResponseWriter, r *http.Request) (bool, error){
	if r.Method == http.MethodPost {
		sess, err := session.SessionStore.Get(r, "session")
		if err != nil {
			return false, err
		}

		err = r.ParseForm()
		if err != nil {
			return false, err
		}
		if r.PostFormValue("password") == viper.GetString("web.adminpw") {
			sess.Values["admin"] = true
			err := sess.Save(r, w)
			if err != nil {
				return false, err
			}
      log.Printf("successful admin login from %s\n", r.RemoteAddr)
			return true, nil
		}
	} 
  log.Printf("failed admin login from %s\n", r.RemoteAddr)
  return false, nil
}
