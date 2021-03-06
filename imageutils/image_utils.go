package imageutils

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type ImageFormat string

const (
	ImageFormatJPG     ImageFormat = "jpg"
	ImageFormatPNG     ImageFormat = "png"
	ImageFormatInvalid ImageFormat = ""
)

// GetImageFormatFromString is a helper to get the ImageFormat from a string
func GetImageFormatFromString(input string) (ImageFormat, error) {
	switch strings.ToLower(input) {
	case "jpg":
		return ImageFormatJPG, nil
	case "jpeg":
		return ImageFormatJPG, nil
	case "png":
		return ImageFormatPNG, nil
	default:
		return ImageFormatInvalid, errors.New("invalid format")
	}
}

// LoadImage loads an image from the filesystem and attempts to decode it
func LoadImage(path string) (image.Image, error) {
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

// SaveImage takes an image and saves it to the target path in the desired format
func SaveImage(img image.Image, format ImageFormat, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create that file: %w", err)
	}
	defer f.Close()

	switch format {
	case ImageFormatJPG:
		err = jpeg.Encode(f, img, nil)
	case ImageFormatPNG:
		err = png.Encode(f, img)
	default:
		return fmt.Errorf("invalid output format passed: %v; only png and jpg are supported at this time", format)
	}

	if err != nil {
		return fmt.Errorf("could not encode that image: %w", err)
	}
	return nil
}
