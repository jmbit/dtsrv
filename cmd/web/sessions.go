package web

import (
	"os"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	_ "github.com/joho/godotenv/autoload"
)

var SessionStore = sessionStore()

func sessionStore() sessions.CookieStore {
  var key []byte
  if keystring, ok := os.LookupEnv("SESSION_KEY"); ok == true {
    key = []byte(keystring)
  } else {
    key = securecookie.GenerateRandomKey(32)
  }
  return *sessions.NewCookieStore(key)

}
