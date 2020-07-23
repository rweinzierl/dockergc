package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch os := args[0]; os {
		case "gc":
			nContainers, nImages := removeObsoleteContainersAndImages()
			fmt.Printf("%d containers and %d images deleted.\n", nContainers, nImages)
			return
		case "pin-container", "pc":
			Add(TypContainer, args[1])
			return
		case "pin-image", "pi":
			Add(TypImage, NormalizeImageName(args[1]))
			return
		case "unpin-container", "uc":
			Remove(TypContainer, args[1])
			return
		case "unpin-image", "ui":
			Remove(TypImage, args[1])
			return
		case "list-pinned-containers", "lc":
			containers := ReadAll(TypContainer)
			for _, container := range containers {
				fmt.Println(container)
			}
			return
		case "list-pinned-images", "li":
			images := ReadAll(TypImage)
			for _, image := range images {
				fmt.Println(image)
			}
			return
		}
	}
	fmt.Printf(`Docker garbage collector (dockergc)

Storage location of dockergc database is determined by environment variable %s
Default storage location is %s

gc: remove unpinned containers and images
pc <name>: pin container
pi <name>: pin image
uc <name>: unpin container
ui <name>: unpin image
lc <name>: list pinned containers
li <name>: list pinned images
`, VarNameDbPath, DbPathDefault())
}
