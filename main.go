package main

import (
	"math/rand"
	"serializer/ui"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	ui.CreateUI()
}
