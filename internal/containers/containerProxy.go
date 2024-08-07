package containers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/docker/docker/client"
)


func GetContainerUrl(ctName string) (string, error) {
  ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
  container, err := cli.ContainerInspect(ctx, ctName)
  if err != nil {
    return "", err
  }
  var port int
  ctIP := container.NetworkSettings.IPAddress
  if os.Getenv("CONTAINER_PORT") != "" {
    port, err = strconv.Atoi(os.Getenv("CONTAINER_PORT"))
    if err != nil {
      log.Println("Invalid CONTAINER_PORT environment variable!")
      return "", err
    }

  } else {
    ctPorts := container.NetworkSettings.Ports
    // Hack to get the first port from the map. 
    for key := range(ctPorts) {
      port = key.Int()
      break
    }
  }
  return fmt.Sprintf("http://%s:%d", ctIP, port), nil
}


