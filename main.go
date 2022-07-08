package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/kevineaton/art/sketch"
)

// todo: pass in command line
var (
	sourceImageName = "source.jpg"
	outputImageName = "out.png"
	totalCycleCount = 20000
)

func main() {
	rand.Seed(time.Now().Unix())

	img, err := loadImage(sourceImageName)
	if err != nil {
		log.Panicln(err)
	}

	destW := 4000
	params := sketch.UserParams{
		DestWidth:                destW,
		DestHeight:               destW,
		StrokeRatio:              0.9,
		StrokeReduction:          0.002,
		StrokeInversionThreshold: 0.05,
		StrokeJitter:             int(0.001 * float64(destW)),
		InitialAlpha:             0.1,
		AlphaIncrease:            0.02,
		MinEdgeCount:             3,
		MaxEdgeCount:             4,
	}

	sketch := sketch.NewSketch(img, params)

	fmt.Printf("\nSketch Params\n%+v\n", params)

	for i := 0; i < totalCycleCount; i++ {
		sketch.Update()
	}

	saveImage(sketch.Output(), outputImageName)
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not load image: %w", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}
	return img, nil
}

func saveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create that file: %w", err)
	}
	err = png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("could not encode that image: %w", err)
	}
	return nil
}
