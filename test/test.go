package main

import (
	"github.com/rweinzierl/dockergc/lib"
)

func main() {
	lib.Run([]string{"", "gc"})
}
