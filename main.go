package main

import (
	"math/rand"
	"time"

	"github.com/kevineaton/art/transformer"
)

var (
	totalCycleCount = 1000
)

func main() {
	rand.Seed(time.Now().Unix())

	transformer.Run()
}
