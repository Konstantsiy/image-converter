// Package converter implements functionality for converting and compressing images.
package converter

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

const (
	// FormatJPEG represents JPEG image format.
	FormatJPEG = "jpeg"
	// FormatJPG represents JPG image format.
	FormatJPG = "jpg"
	// FormatPNG represents PNG image format.
	FormatPNG = "png"
)

// Convert converts and compresses the given image file according to the target format and compression ratio.
func Convert(reader io.Reader, targetFormat string, ratio int) (io.ReadSeeker, error) {
	imageData, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("can't decode source file: %w", err)
	}

	buf := new(bytes.Buffer)

	switch targetFormat {
	case FormatPNG:
		var enc png.Encoder
		enc.CompressionLevel = png.CompressionLevel(ratio)
		err := enc.Encode(buf, imageData)
		if err != nil {
			return nil, fmt.Errorf("can't convert image to %s format: %w", FormatJPG, err)
		}
	case FormatJPEG, FormatJPG:
		err := jpeg.Encode(buf, imageData, &jpeg.Options{
			Quality: ratio,
		})
		if err != nil {
			return nil, fmt.Errorf("can't convert image to %s format: %w", FormatJPG, err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", targetFormat)
	}

	return bytes.NewReader(buf.Bytes()), nil
}
