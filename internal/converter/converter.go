// Package converter implements functionality for converting and compressing images.
package converter

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
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
		return nil, err
	}

	outputFile, err := os.Create("temp/file." + targetFormat)
	if err != nil {
		return nil, errors.New("can't create target file")
	}

	switch targetFormat {
	case "png":
		var enc png.Encoder
		enc.CompressionLevel = png.CompressionLevel(ratio)
		err := enc.Encode(outputFile, imageData)
		if err != nil {
			return nil, err
		}
	case "jpeg":
	case "jpg":
		{
			err := jpeg.Encode(outputFile, imageData, &jpeg.Options{
				Quality: ratio,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return outputFile, nil
}
