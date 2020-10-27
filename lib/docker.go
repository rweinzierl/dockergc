package lib

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var theDockerClient *client.Client

func getDocker() *client.Client {
	if theDockerClient == nil {
		newClient, err := client.NewEnvClient()
		exitOnError(err)
		theDockerClient = newClient
	}
	return theDockerClient
}

func getAllImages() *[]types.ImageSummary {
	images, err := getDocker().ImageList(context.Background(), types.ImageListOptions{All: true})
	exitOnError(err)
	return &images
}

func getAllContainers() *[]types.Container {
	containers, err := getDocker().ContainerList(context.Background(), types.ContainerListOptions{All: true})
	exitOnError(err)
	return &containers
}

func deleteImage(image *types.ImageSummary) {
	getDocker().ImageRemove(context.Background(), image.ID, types.ImageRemoveOptions{})
}

func deleteContainer(container *types.Container) {
	getDocker().ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
}
