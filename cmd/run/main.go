package main

import (
	"fmt"
	"os"

	"github.com/jmbit/dtsrv/internal/config"
	"github.com/jmbit/dtsrv/internal/server"
	"github.com/jmbit/dtsrv/internal/session"
	"github.com/jmbit/dtsrv/lib/containers"
	"github.com/spf13/viper"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var err error

  // Read configuration
  config.ReadConfigFile("")

	// Start background job to fetch container image
	go containers.PullContainer(viper.GetString("container.image"))

  go containers.StartCleanup(viper.GetInt64("container.maxage"), 600)

  //Setup Session store etc
  session.InitSessions()

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
