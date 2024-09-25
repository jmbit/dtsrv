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

	// Get or generate Admin PW
	if os.Getenv("ADMIN_PW") == "" {
		pw := string(securecookie.GenerateRandomKey(32))
		os.Setenv("ADMIN_PW", pw)
		log.Printf("Set Admin PW to %s temporarily, please set a proper admin PW.\n", pw)
	}
	// Start background job to fetch container image
	go containers.PullContainer()

	server := server.NewServer()

	if os.Getenv("USE_TLS") == "true" {
		err = server.ListenAndServeTLS(os.Getenv("TLS_CRT"), os.Getenv("TLS_KEY"))

	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
