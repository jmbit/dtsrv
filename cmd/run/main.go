package main

import (
	"fmt"
	"os"

	"github.com/jmbit/dtsrv/internal/config"
	"github.com/jmbit/dtsrv/internal/server"
	"github.com/jmbit/dtsrv/lib/containers"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var err error
  config.ReadConfigFile("")

	// Start background job to fetch container image
	go containers.PullContainer("")

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
