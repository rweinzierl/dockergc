package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client

func removeObsoleteContainersAndImages() (int, int) {
	cli, _ = client.NewEnvClient()
	runningContainerIds := make(map[string]bool)
	runningContainers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	for _, container := range runningContainers {
		runningContainerIds[container.ID] = true
	}
	pinnedContainerNames := make(map[string]bool)
	for _, name := range ReadAll(TypContainer) {
		pinnedContainerNames[NormalizeContainerName(name)] = true
	}
	containerToKeepIds := make(map[string]bool)
	imageToKeepIds := make(map[string]bool)
	containers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	for _, container := range containers {
		if _, ok := runningContainerIds[container.ID]; ok {
			containerToKeepIds[container.ID] = true
			imageToKeepIds[container.ImageID] = true
		}
		for _, name := range container.Names {
			if _, ok := pinnedContainerNames[name]; ok {
				containerToKeepIds[container.ID] = true
				imageToKeepIds[container.ImageID] = true
			}
		}
	}
	pinnedImageNames := make(map[string]bool)
	for _, name := range ReadAll(TypImage) {
		pinnedImageNames[name] = true
	}
	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if _, ok := pinnedImageNames[tag]; ok {
				imageToKeepIds[image.ID] = true
			}
		}
	}
	nContainers := 0
	for _, container := range containers {
		if _, ok := containerToKeepIds[container.ID]; !ok {
			cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
			nContainers++
		}
	}
	nImages := 0
	for _, image := range images {
		if _, ok := imageToKeepIds[image.ID]; !ok {
			cli.ImageRemove(context.Background(), image.ID, types.ImageRemoveOptions{})
			nImages++
		}
	}
	return nContainers, nImages
}
