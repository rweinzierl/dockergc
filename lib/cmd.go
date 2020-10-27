package lib

import (
	"fmt"
	"strings"

	"github.com/docker/docker/client"
)

func listImagePatterns(cli *client.Client) {
	w := createTabWriter()
	addTabRow(w, "IMAGE")
	for _, pattern := range *readAll(typImage) {
		addTabRow(w, pattern)
	}
}

func listImages(cli *client.Client, imageNamePatterns []string) {
	imageNamePatterns = replaceEmptyPatternListByWildcard(imageNamePatterns)
	w := createTabWriter()
	addTabRow(w, "REPOSITORY", "TAG", "IMAGE ID", "GC")
	for _, image := range *getData().allImages {
		if imageMatchesAnyPattern(&image, &imageNamePatterns) {
			for _, repoTag := range image.RepoTags {
				repo, tag := splitImageName(repoTag)
				_, id := splitImageID(image.ID)
				addTabRow(w, repo, tag, id[:12], imageMayBeDeleted(&image))
			}
		}
	}
	w.Flush()
}

func listContainerPatterns(cli *client.Client) {
	w := createTabWriter()
	addTabRow(w, "CONTAINER")
	for _, pattern := range *readAll(typContainer) {
		addTabRow(w, pattern)
	}
}

func listContainers(cli *client.Client, containerNamePatterns []string) {
	containerNamePatterns = replaceEmptyPatternListByWildcard(containerNamePatterns)
	w := createTabWriter()
	addTabRow(w, "CONTAINER ID", "IMAGE", "NAMES", "GC")
	for _, container := range *getData().allContainers {
		if containerMatchesAnyPattern(&container, &containerNamePatterns) {
			addTabRow(w, container.ID[:12], container.Image, strings.Join(container.Names, ", "), containerMayBeDeleted(&container))
		}
	}
	w.Flush()
}

func Run(args []string) {
	cli, err := client.NewEnvClient()
	exitOnError(err)
	if len(args) > 1 {
		args = args[1:]
		switch subCommand := args[0]; subCommand {
		case "gc":
			gc()
			return
		case "pi":
			for _, pattern := range args[1:] {
				add(typImage, pattern)
			}
			return
		case "ui":
			for _, pattern := range args[1:] {
				remove(typImage, pattern)
			}
			return
		case "lpi":
			listImagePatterns(cli)
			return
		case "li":
			listImages(cli, args[1:])
			return
		case "pc":
			for _, pattern := range args[1:] {
				add(typContainer, pattern)
			}
			return
		case "uc":
			for _, pattern := range args[1:] {
				remove(typContainer, pattern)
			}
			return
		case "lpc":
			listContainerPatterns(cli)
			return
		case "lc":
			listContainers(cli, args[1:])
			return
		}
	}
	fmt.Printf(`Docker garbage collector (dockergc)

Storage location of dockergc database is determined by environment variable %s
Default storage location is %s

gc				: remove unpinned containers and images
li [PATTERN]...	: list images and if they are due for garbage collection
lc [PATTERN]...	: list containers and if they are due for garbage collection
lpi				: list pinned images
lpc				: list pinned containers
pi [PATTERN]...	: pin image
pc [PATTERN]...	: pin container
ui [PATTERN]...	: unpin image
uc [PATTERN]...	: unpin container
`, VarNameDbPath, dbPathDefault())
}
