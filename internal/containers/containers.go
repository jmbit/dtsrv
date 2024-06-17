package containers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

// CreateContainer() wraps creatContainer asyncronously to avoid UI lag
func CreateContainer() string {
  ctName := fmt.Sprintf("dtsrv-%s", uuid.New().String())
  go createContainer(ctName)
  return ctName
}


func createContainer(ctName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
    log.Println("Error creating container,", err)
    return
	}
	defer cli.Close()

	imageName := "lscr.io/linuxserver/webtop"

	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
    log.Println("Error creating container,", err)
    return
	}
	defer out.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
    // Env: []string{
    //   fmt.Sprintf("SUBFOLDER=/ct/%s/", ctName),
    // },
		Image: imageName,
	}, nil, nil, nil, ctName)
  	if err != nil {
    log.Println("Error creating container,", err)
    return
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
    log.Println("Error creating container,", err)
    return
	}

	fmt.Println(ctName, resp.ID)


}

func PullContainer() error {
  log.Println("pullling container image")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
    log.Println("Error pulling Docker image:", err)
    return  err
	}
	defer cli.Close()

	imageName := "lscr.io/linuxserver/webtop"

	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
    log.Println("Error pulling Docker image:", err)
    return err
	}
	defer out.Close()
  io.Copy(os.Stdout, out)

  log.Println("Done pulling container image")
  return nil
}

func CleanupContainers() {
  log.Println("Deleting previously created containers")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
    log.Println("Error creating container,", err)
    return
	}
	defer cli.Close()

  cli.ContainerList(ctx, container.ListOptions{
    All: true,
    Filters: filters.Args{},

})

}
