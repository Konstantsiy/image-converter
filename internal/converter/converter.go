// Package converter implements functionality for converting and compressing images.
package converter

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"time"
)

const (
	TempPath   = "temp/"
	FormatJPEG = "jpeg"
	FormatJPG  = "jpg"
	FormatPNG  = "png"
)

// Converter converts and compresses images.
type Converter struct{}

// NewConverter creates new converter.
func NewConverter() *Converter {
	return &Converter{}
}

//Convert converts and compresses the given image file according to the target format and compression ratio.
func (c *Converter) Convert(file multipart.File, targetFormat string, ratio int) (*os.File, error) {
	imageData, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("can't decode source file: %w", err)
	}

	outputFile, err := os.Create(TempPath + time.Now().String() + "." + targetFormat)
	if err != nil {
		return nil, errors.New("can't create target file")
	}

	switch targetFormat {
	case FormatPNG:
		var enc png.Encoder
		enc.CompressionLevel = png.CompressionLevel(ratio)
		err := enc.Encode(outputFile, imageData)
		if err != nil {
			return nil, fmt.Errorf("can't convert image to %s format: %w", FormatJPG, err)
		}
	case FormatJPEG, FormatJPG:
		{
			err := jpeg.Encode(outputFile, imageData, &jpeg.Options{
				Quality: ratio,
			})
			if err != nil {
				return nil, fmt.Errorf("can't convert image to %s format: %w", FormatJPG, err)
			}
		}
	}
	return outputFile, nil
}
