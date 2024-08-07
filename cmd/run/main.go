package main

import (
	"dtsrv/internal/containers"
	"dtsrv/internal/server"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/securecookie"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
  var err error

  if os.Getenv("ADMIN_PW") == "" {
    pw := string(securecookie.GenerateRandomKey(32))
    os.Setenv("ADMIN_PW", pw)
    log.Printf("Set Admin PW to %s temporarily, please set a proper admin PW.\n", pw)
  }
  go containers.PullContainer()


	server := server.NewServer()
//  db.Connect()

  if os.Getenv("USE_TLS") == "true" {
    err = server.ListenAndServeTLS(os.Getenv("TLS_CRT"), os.Getenv("TLS_KEY"))

  } else {
	  err = server.ListenAndServe()
  }
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
