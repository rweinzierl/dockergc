package main

import (
	"log"
	"strings"
)

func NormalizeContainerName(name string) string {
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}
	return name
}

func NormalizeImageName(name string) string {
	if !strings.Contains(name, ":") {
		name = name + ":latest"
	}
	return name
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
