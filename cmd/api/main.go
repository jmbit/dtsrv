package main

import (
	"dtsrv/internal/containers"
	"dtsrv/internal/server"
	"fmt"
	"log"
)

func main() {
  log.Println(containers.PullContainer()) 

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
