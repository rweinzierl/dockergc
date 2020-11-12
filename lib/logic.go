package lib

import (
	"github.com/docker/docker/api/types"
)

func containerIsRunning(container *types.Container) bool {
	runningStates := []string{"running", "paused", "restarting"}
	for _, state := range runningStates {
		if container.State == state {
			return true
		}
	}
	return false
}

func containerMayBeDeleted(container *types.Container) bool {
	return !containerMatchesAnyPattern(container, getData().pinnedContainerPatterns) && !containerIsRunning(container)
}

func imageIsLockedByContainer(image *types.ImageSummary) bool {
	for _, container := range *getData().allContainers {
		if image.ID == container.ImageID && !containerMayBeDeleted(&container) {
			return true
		}
	}
	return false
}

func imageMayBeDeleted(image *types.ImageSummary) bool {
	return !imageMatchesAnyPattern(image, getData().pinnedImagePatterns) && image.Containers <= 0 && !imageIsLockedByContainer(image)
}

func gcContainers() {
	for _, container := range *getAllContainers() {
		if containerMayBeDeleted(&container) {
			deleteContainer(&container)
		}
	}
}

func gcImages() {
	deleteImageDanglingImages()
	for _, image := range *getAllImages() {
		if imageMayBeDeleted(&image) {
			deleteImage(&image)
		}
	}
}

func gc() {
	gcContainers()
	reloadData()
	gcImages()
}
