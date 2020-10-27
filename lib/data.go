package lib

import "github.com/docker/docker/api/types"

type data struct {
	allContainers           *[]types.Container
	pinnedContainerPatterns *[]string
	allImages               *[]types.ImageSummary
	pinnedImagePatterns     *[]string
}

var theData *data

func reloadData() {
	theData = new(data)
	theData.allContainers = getAllContainers()
	theData.allImages = getAllImages()
	theData.pinnedContainerPatterns = readAll(typContainer)
	theData.pinnedImagePatterns = readAll(typImage)
}

func getData() *data {
	if theData == nil {
		reloadData()
	}
	return theData
}
