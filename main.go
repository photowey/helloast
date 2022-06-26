package main

import (
	"math/rand"
	"time"

	App "github.com/photowey/helloast/cmd/app"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// @App
func main() {
	App.Run()
}
