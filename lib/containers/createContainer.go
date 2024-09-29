package containers

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

// CreateContainer() wraps creatContainer asyncronously to avoid UI lag,
// wraps createContainer(). If Image "" is specified, will use default instead
func CreateContainer(dockerImage string, isolate bool) (string, error) {
	ctName := fmt.Sprintf("dtsrv-%s", uuid.New().String())
	err := createContainer(ctName, dockerImage, isolate)
	return ctName, err
}

// createContainer() creates a container with the given name.
func createContainer(ctName string, dockerImage string, isolate bool) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error creating container,", err)
		return err
	}
	defer cli.Close()

	var imageName string
	// Get image from environment variable, fallback to default if required
	if dockerImage != "" {
		imageName = dockerImage
	} else {
		imageName = "lscr.io/linuxserver/webtop"
	}
  containerConfig := container.Config{
		Env: []string{
			fmt.Sprintf("SUBFOLDER=/view/%s/", ctName),
		},
		Image: imageName,
  }


	// Create container
	resp, err := cli.ContainerCreate(ctx, &containerConfig, nil, nil, nil, ctName)
	if err != nil {
		log.Println("Error creating container,", err)
		return err
	}

  if isolate == true {
    resp, err := cli.NetworkCreate(ctx, ctName, types.NetworkCreate{})
	  if err != nil {
	  	log.Println("Error creating container network,", err)
	  	return err
	  }
    err = cli.NetworkConnect(ctx, resp.ID, ctName, &network.EndpointSettings{})
	  if err != nil {
	  	log.Println("Error connecting container to network,", err)
      DeleteContainer(ctName)
	  	return err
	  }
  }
	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Println("Error creating container,", err)
		return err
	}
	fmt.Println(ctName, resp.ID)
	return nil
}

