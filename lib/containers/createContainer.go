package containers

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

// CreateContainer() wraps creatContainer asyncronously to avoid UI lag,
// wraps createContainer(). If Image "" is specified, will use default instead
func CreateContainer(dockerImage string) (string, error) {
	ctName := fmt.Sprintf("dtsrv-%s", uuid.New().String())
	err := createContainer(ctName, dockerImage)
	return ctName, err
}

// createContainer() creates a container with the given name.
func createContainer(ctName string, dockerImage string) error {
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

	// Create container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Env: []string{
			fmt.Sprintf("SUBFOLDER=/view/%s/", ctName),
		},
		Image: imageName,
	}, nil, nil, nil, ctName)
	if err != nil {
		log.Println("Error creating container,", err)
		return err
	}

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Println("Error creating container,", err)
		return err
	}
	fmt.Println(ctName, resp.ID)
	return nil
}
