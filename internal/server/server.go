package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmbit/dtsrv/internal/middlewares"
	"github.com/spf13/viper"

	"github.com/gorilla/handlers"

	_ "github.com/joho/godotenv/autoload"
)

// NewServer() builds a http server with middlewares and config from env
func NewServer() *http.Server {


	middlewareStack := middlewares.CreateStack(
		middlewares.AssetCaching,
		middlewares.GorillaLogging,
		handlers.RecoveryHandler(),
	)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("web.host"), viper.GetInt("web.port")),
		Handler:      middlewareStack(registerRoutes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
