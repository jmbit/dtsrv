package containers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/docker/docker/client"
)

// GetContainerUrl finds the IP and port a container listens on
func GetContainerUrl(ctName string) (string, error) {
	var port int
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
	// attempt to get container listen port from environment variable, falling back to first port configured
	if os.Getenv("CONTAINER_PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("CONTAINER_PORT"))
		if err != nil {
			log.Println("Invalid CONTAINER_PORT environment variable!")
			return "", err
		}
	} else {
		ctPorts := container.NetworkSettings.Ports
		// Hack to get the first port from the map.
		for key := range ctPorts {
			port = key.Int()
			break
		}
	}
	return fmt.Sprintf("http://%s:%d", ctIP, port), nil
}
