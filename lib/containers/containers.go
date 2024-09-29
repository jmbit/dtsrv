package containers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	_ "github.com/joho/godotenv/autoload"
)

// PullContainer() pulls the container image
func PullContainer(dockerImage string) error {
	log.Println("pulling container image")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error pulling Docker image:", err)
		return err
	}
	defer cli.Close()

	var imageName string
	if dockerImage != "" {
		imageName = dockerImage 
	} else {

    imageName = "lscr.io/linuxserver/webtop:alpine-icewm"
	}

	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		log.Println("Error pulling Docker image:", err)
		return err
	}
  _, err = io.ReadAll(out)
  if err != nil {
    log.Println(err)
  }
	defer out.Close()

	log.Println("Done pulling container image")
	return nil
}

// ListContainers() gets a list of all docker containers with the "dtsrv-"-Prefix
func ListContainers() ([]types.Container, error) {
	var returnList []types.Container
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error creating container client,", err)
		return nil, err
	}
	defer cli.Close()

	// No idea how the filters work, they aren't really documented, so rip me
	// This is going to suck hard if you have hundreds of containers
	list, err := cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.Args{},
	})

	if err != nil {
		log.Println("Error listing containers,", err)
		return nil, err
	}

	for _, ct := range list {
		if strings.HasPrefix(ct.Names[0], "/dtsrv-") {
			returnList = append(returnList, ct)

		}
	}

	return returnList, nil
}

// GetContainer() wraps ListContainers with a filter for the container name, because for reasons ContainerList and ContainerInspect
// give different types representing the same thing
func GetContainer(ctName string) (types.Container, error) {
	containerList, err := ListContainers()
	if err != nil {
		return types.Container{}, err
	}
	for _, ct := range containerList {
		if ct.Names[0] == fmt.Sprintf("/%s", ctName) {
			return ct, nil
		}
	}
	log.Println("Could not find container", ctName)
	return types.Container{}, fmt.Errorf("Could not find container %s", ctName)
}

// StopContainer() tries to stop a docker container
func StopContainer(ctName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error creating container client,", err)
		return err
	}
	err = cli.ContainerStop(ctx, ctName, container.StopOptions{})
	if err != nil {
		log.Println("Error stopping container, ", err)
		return err
	}
	log.Println("Stopped Container", ctName)
	return nil
}

// DeleteContainer() deletes a docker container
func DeleteContainer(ctName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Error creating container client,", err)
		return err
	}
	err = cli.ContainerStop(ctx, ctName, container.StopOptions{})
	if err != nil {
		log.Println("Error stopping container, ", err)
		return err
	}
	err = cli.ContainerRemove(ctx, ctName, container.RemoveOptions{})
	if err != nil {
		log.Println("Error removing container, ", err)
		return err
	}
	log.Println("Deleted Container", ctName)
	return nil
}

// containerNameToDockerID() looks up the Docker container ID from its name (idk why the API has no relevant info)
func containerNameToDockerID(ctName string) (string, error) {
	ctList, err := ListContainers()
	if err != nil {
		return "", err
	}
	for _, ct := range ctList {
		if ct.Names[0] == fmt.Sprintf("/%s", ctName) {
			return ct.ID, nil
		}
	}
	return "", nil
}

// GetImageName() gets the name for the container image currently in use
func GetImageName() string {
	// Get default image
	var imageName string
	if envImage, exists := os.LookupEnv("IMAGE_NAME"); exists == true {
		imageName = envImage
	} else {
		imageName = "lscr.io/linuxserver/webtop"
	}

	return imageName
}
