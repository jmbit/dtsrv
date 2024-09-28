package server

import (
	"github.com/jmbit/dtsrv/internal/middlewares"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"

	_ "github.com/joho/godotenv/autoload"
)

// NewServer() builds a http server with middlewares and config from env
func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	middlewareStack := middlewares.CreateStack(
		middlewares.AssetCaching,
		middlewares.GorillaLogging,
		handlers.RecoveryHandler(),
	)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      middlewareStack(registerRoutes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
