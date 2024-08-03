package containers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/docker/docker/client"
)

var containerTimeouts sync.Map
var maxTimeout = 30

// testConnectionToContainer() takes the Id of a container and checks if it's reachable on Port 3000
func TestConnectionToContainer(ctName string) (bool, error) {
  ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
  container, err := cli.ContainerInspect(ctx, ctName)
  if err != nil {
    return false, err
  }
  timeoutCounter := getTimeoutCount(ctName)
  ctIP := container.NetworkSettings.IPAddress
  ctPorts := container.NetworkSettings.Ports
  // Hack to get the first port from the map. 
  var port int
  for key := range(ctPorts) {
    port = key.Int()
    break
  }
  log.Printf("Trying to connect to %s:%v", ctIP, port)
  resp, err := http.Get(fmt.Sprintf("http://%s:%d", ctIP, port))
  log.Println("Code", resp.Status, err)
  if resp.StatusCode == 200 {
      return true, nil
  }
  log.Printf("Trying to connect to http://%s:%d/view/%s/", ctIP, port, ctName)
  resp, err = http.Get(fmt.Sprintf("http://%s:%d/view/%s/", ctIP, port, ctName))
  log.Println("Code", resp.Status, err)
  if resp.StatusCode == 200 {
      return true, nil
  }
  if timeoutCounter > maxTimeout {
    return false, fmt.Errorf("Could not connect to %s on IP %s", ctName, ctIP)
  }
  return false, nil

}

func getTimeoutCount(ctName string) int {
  counter, ok := containerTimeouts.Load(ctName)
  if !ok {
    containerTimeouts.Store(ctName, 0)
    return 0
  } else {
    return counter.(int)
  }

}

