package containers

import (
	"context"
	"fmt"
	"log"
  "os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

// CreateContainer() wraps creatContainer asyncronously to avoid UI lag
func CreateContainer() (string, error) {
  ctName := fmt.Sprintf("dtsrv-%s", uuid.New().String())
  err :=  createContainer(ctName)
  return ctName, err
}


func createContainer(ctName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
    log.Println("Error creating container,", err)
    return err
	}
	defer cli.Close()

  var imageName string
  if envImage, exists := os.LookupEnv("IMAGE_NAME"); exists == true {
    imageName = envImage
  } else {

	imageName = "lscr.io/linuxserver/webtop"
  }

  err = PullContainer()
	if err != nil {
    log.Println("Error creating container,", err)
    return err
	}

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

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
    log.Println("Error creating container,", err)
    return err
	}
	fmt.Println(ctName, resp.ID)
  return nil
}


