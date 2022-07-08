package transformer

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	_ "image/jpeg"

	"github.com/fogleman/gg"
	"github.com/kevineaton/art/imageutils"
)

const totalCycleCount = 10

type TransformerUserParams struct {
	DestWidth                int
	DestHeight               int
	StrokeRatio              float64
	StrokeReduction          float64
	StrokeJitter             int
	StrokeInversionThreshold float64
	InitialAlpha             float64
	AlphaIncrease            float64
	MinEdgeCount             int
	MaxEdgeCount             int
}

type TransformerSketch struct {
	TransformerUserParams
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
}

// NewTransformerSketch creates a new transforming sketch to generate art based upon a source image
func NewTransformerSketch(source image.Image, userParams TransformerUserParams) *TransformerSketch {
	s := &TransformerSketch{TransformerUserParams: userParams}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.StrokeRatio * float64(s.DestWidth)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()

	s.source = source
	s.dc = canvas
	return s
}

// Update draws on each cycle of the algorithm
func (s *TransformerSketch) Update() {
	// get the color info
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := rgb255(s.source.At(int(rndX), int(rndY)))

	// determine the output
	destX := rndX * float64(s.DestWidth) / float64(s.sourceWidth)
	destX += float64(randRange(s.StrokeJitter))
	destY := rndY * float64(s.DestHeight) / float64(s.sourceHeight)
	destY += float64(randRange(s.StrokeJitter))

	// draw the stroke
	edges := s.MinEdgeCount + rand.Intn(s.MaxEdgeCount-s.MinEdgeCount+1)

	s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
	s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rand.ExpFloat64())
	s.dc.FillPreserve()

	if s.strokeSize <= s.StrokeInversionThreshold*s.initialStrokeSize {
		if (r+g+b)/3 < 128 {
			s.dc.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
		} else {
			s.dc.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
		}
	}
	s.dc.Stroke()

	s.strokeSize -= s.StrokeReduction * s.strokeSize
	s.InitialAlpha += s.AlphaIncrease

}

// Run is the entry point and where config options will be passed when implemented
func Run() {
	rand.Seed(time.Now().Unix())

	files, err := ioutil.ReadDir("./input")
	if err != nil {
		log.Panicln(err)
	}

	for i := range files {
		fileName := files[i].Name()
		// split on the name to identify the file type
		parts := strings.Split(fileName, ".")
		if len(parts) < 2 {
			continue // not a valid file type
		}
		extension := parts[len(parts)-1]
		if extension != "jpg" && extension != "jpeg" && extension != "png" {
			continue
		}
		outputName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".png"

		fmt.Printf("Working on %s, outputting to %s\n", fileName, outputName)

		// now handle the file
		img, err := imageutils.LoadImage("./input/" + fileName)
		if err != nil {
			log.Panicln(err)
		}

		destW := 4000
		params := TransformerUserParams{
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

		sketch := NewTransformerSketch(img, params)

		for i := 0; i < totalCycleCount; i++ {
			sketch.Update()
		}

		err = imageutils.SaveImage(sketch.Output(), imageutils.ImageFormatPNG, "./output/"+outputName)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}

}

// Output generates the output of the transformation
func (s *TransformerSketch) Output() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}
