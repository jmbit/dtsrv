package containers

import (
	"context"
	"fmt"

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
  ctIP := container.NetworkSettings.IPAddress
  ctPorts := container.NetworkSettings.Ports
  // Hack to get the first port from the map. 
  var port int
  for key := range(ctPorts) {
    port = key.Int()
    break
  }
  return fmt.Sprintf("http:%s:%d", ctIP, port, ctName), nil

}
