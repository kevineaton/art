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
	"github.com/kevineaton/art/progressbar"
	"github.com/spf13/cobra"
)

const totalCycleCount = 1000

type TransformerUserParams struct {
	DestWidth                int
	DestHeight               int
	StrokeRatio              float64
	StrokeReduction          float64
	StrokeJitter             int
	StrokeJitterRatio        float64
	StrokeInversionThreshold float64
	InitialAlpha             float64
	AlphaIncrease            float64
	MinEdgeCount             int
	MaxEdgeCount             int
	OutputFileType           string
	TotalCycles              int
}

type TransformerSketch struct {
	*TransformerUserParams
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
}

func GetCommand() *cobra.Command {
	params := &TransformerUserParams{}
	cmd := &cobra.Command{
		Use:   "transform",
		Short: "Transform the images in input to output",
		Run: func(cmd *cobra.Command, args []string) {
			Run(params)
			fmt.Printf("Done!\n")
		},
	}
	cmd.Flags().IntVar(&params.DestHeight, "dest-height", 1000, "Height of the destination target")
	cmd.Flags().IntVar(&params.DestWidth, "dest-width", 1000, "Width of the destination target")
	cmd.Flags().Float64Var(&params.StrokeJitterRatio, "stroke-jitter-ratio", .001, "How much jitter or deviation we add for targets")
	cmd.Flags().Float64Var(&params.StrokeRatio, "stroke-ratio", .75, "Size of the stroke compared to the final result")
	cmd.Flags().Float64Var(&params.StrokeReduction, "stroke-reduction", .002, "Reduce the stroke by this amount on each iteration")
	cmd.Flags().Float64Var(&params.StrokeInversionThreshold, "stroke-inversion-threshold", .05, "Once crossed, we add borders for visibility")
	cmd.Flags().Float64Var(&params.InitialAlpha, "initial-alpha", .1, "The initial transparency and we build up on each iteration")
	cmd.Flags().Float64Var(&params.AlphaIncrease, "alpha-increase", .02, "How much alpha to increase by on each iteration")
	cmd.Flags().IntVar(&params.MinEdgeCount, "min-edges", 3, "The minimum number of edges for each shape")
	cmd.Flags().IntVar(&params.MaxEdgeCount, "max-edges", 4, "The maximum number of edges for each shape")
	cmd.Flags().StringVar(&params.OutputFileType, "output-type", "png", "The desired output, either png or jpg; if set incorrectly, will be set to png")
	cmd.Flags().IntVar(&params.TotalCycles, "cycles", 10000, "The number of iterations to apply the transformation")
	return cmd
}

// Run is the entry point and where config options will be passed when implemented
func Run(params *TransformerUserParams) {
	rand.Seed(time.Now().Unix())

	files, err := ioutil.ReadDir("./input")
	if err != nil {
		log.Panicln(err)
	}

	// set some calculated values
	params.StrokeJitter = int(params.StrokeJitterRatio * float64(params.DestWidth))
	format, err := imageutils.GetImageFormatFromString(params.OutputFileType)
	if err != nil {
		format = imageutils.ImageFormatPNG
	}

	now := time.Now().Format("2006-01-02T15:04:05")

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
		outputName := fmt.Sprintf("%s_%s_transformed.%s", strings.TrimSuffix(fileName, filepath.Ext(fileName)), now, params.OutputFileType)

		// now handle the file
		img, err := imageutils.LoadImage("./input/" + fileName)
		if err != nil {
			log.Panicln(err)
		}

		sketch := newTransformerSketch(img, params)

		bar := progressbar.GetProgressBar(&progressbar.BarOptions{
			Max:          totalCycleCount,
			Width:        100,
			EnableColors: true,
			Description:  fmt.Sprintf("[%d of %d] Transforming %s to %s", i, len(files)-1, fileName, outputName),
		})

		for i := 0; i < totalCycleCount; i++ {
			bar.Add(1)
			sketch.update()
		}

		err = imageutils.SaveImage(sketch.output(), format, "./output/"+outputName)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}

}

// newTransformerSketch creates a new transforming sketch to generate art based upon a source image
func newTransformerSketch(source image.Image, userParams *TransformerUserParams) *TransformerSketch {
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

// update draws on each cycle of the algorithm
func (s *TransformerSketch) update() {
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

// output generates the output of the transformation
func (s *TransformerSketch) output() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}
