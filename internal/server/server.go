package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmbit/dtsrv/internal/middlewares"
	"github.com/spf13/viper"

	"github.com/gorilla/handlers"

	_ "github.com/joho/godotenv/autoload"
)

// NewServer() builds a http server with middlewares and config from env
func NewServer() *http.Server {
  address := fmt.Sprintf("%s:%d", viper.GetString("web.host"), viper.GetInt("web.port"))
  log.Println("Listening on", address)
  var middlewareStack middlewares.Middleware

  if viper.GetBool("web.loghttp") {
	  middlewareStack = middlewares.CreateStack(
	  	middlewares.AssetCaching,
	  	middlewares.GorillaLogging,
	  	handlers.RecoveryHandler(),
	  )
  } else {
  	middlewareStack = middlewares.CreateStack(
  		middlewares.AssetCaching,
  		handlers.RecoveryHandler(),
  	)   
  }


	// Declare Server config
	server := &http.Server{
		Addr:         address,
		Handler:      middlewareStack(registerRoutes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
