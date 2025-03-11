package main

import (
	"math/rand"
	"serializer/ui"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ui.CreateUI()
}
