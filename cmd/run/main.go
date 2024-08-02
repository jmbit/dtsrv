package main

import (
	"dtsrv/internal/db"
	"dtsrv/internal/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
  var err error

  //Gracefully handle SIGINT (ctrl+c)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		fmt.Println()
		fmt.Printf("Received signal: %s\n", sig)
    db.Close()
		fmt.Println("Cleanup done. Exiting...")
    os.Exit(0)
	}()


	server := server.NewServer()
  db.Connect()

  if os.Getenv("USE_TLS") == "true" {
    err = server.ListenAndServeTLS(os.Getenv("TLS_CRT"), os.Getenv("TLS_KEY"))

  } else {
	  err = server.ListenAndServe()
  }
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
