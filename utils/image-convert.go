package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
)

// ResizeAndConvertToWebP resizes input image to a specified width and converts it to WebP
func ResizeAndConvertToWebP(inputPath string, outputPath string, width uint) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read WebP magic bytes
	header := make([]byte, 4)
	if _, err := file.Read(header); err != nil {
		return err
	}

	// Checking magic bytes
	if bytes.Equal(header, []byte("RIFF")) {
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}

		img, err := webp.Decode(file, &decoder.Options{})
		if err != nil {
			return err
		}

		return processImage(img, outputPath, width)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	var img image.Image

	switch format {
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
	case "jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported file type") // Unsupported file type
	}

	return processImage(img, outputPath, width)
}

// processImage resizes the image to the specified width while maintaining aspect ratio and encodes it to WebP format.
func processImage(img image.Image, outputPath string, width uint) error {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	newHeight := uint(float64(originalHeight) * (float64(width) / float64(originalWidth)))

	resizedImg := resize.Resize(width, newHeight, img, resize.Lanczos3)

	outputDir := filepath.Dir(outputPath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return err
	}

	if err := webp.Encode(output, resizedImg, options); err != nil {
		return err
	}

	return nil
}
