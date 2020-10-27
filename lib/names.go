package lib

import (
	"strings"

	"github.com/docker/docker/api/types"
)

func normalizeContainerName(name string) string {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	return name
}

func normalizeImageName(name string) string {
	if !strings.Contains(name, ":") {
		name = name + ":latest"
	}
	return name
}

func splitImageName(name string) (repo, tag string) {
	name = normalizeImageName(name)
	pos := strings.LastIndex(name, ":")
	return name[:pos], name[pos+1:]
}

func splitImageID(name string) (algorithm, hash string) {
	name = normalizeImageName(name)
	pos := strings.Index(name, ":")
	return name[:pos], name[pos+1:]
}

func containerNameMatchesPattern(name string, pattern string) bool {
	name = normalizeContainerName(name)
	for _, name2 := range []string{name, name[1:]} {
		if stringMatchesPattern(name2, pattern) {
			return true
		}
	}
	return false
}

func imageNameMatchesPattern(name string, pattern string) bool {
	name = normalizeImageName(name)
	repo, _ := splitImageName(name)
	for _, name2 := range []string{name, repo} {
		if stringMatchesPattern(name2, pattern) {
			return true
		}
	}
	return false
}

func imageIDMatchesPattern(id string, pattern string) bool {
	_, id2 := splitImageID(id)
	for _, id3 := range []string{id, id2} {
		if stringMatchesPattern(id3, pattern) {
			return true
		}
	}
	return false
}

func containerMatchesAnyPattern(container *types.Container, patterns *[]string) bool {
	for _, pattern := range *patterns {
		if stringMatchesPattern(container.ID, pattern) {
			return true
		}
		for _, containerName := range container.Names {
			if containerNameMatchesPattern(containerName, pattern) {
				return true
			}
		}
	}
	return false
}

func imageMatchesAnyPattern(image *types.ImageSummary, patterns *[]string) bool {
	for _, pattern := range *patterns {
		if stringMatchesPattern(image.ID, pattern) {
			return true
		}
		for _, repoTag := range image.RepoTags {
			if imageNameMatchesPattern(repoTag, pattern) {
				return true
			}
		}
	}
	return false
}
