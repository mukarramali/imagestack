package src

import (
	"fmt"

	"github.com/h2non/bimg"
)

func CompressImage(inputPath string, outputPath string, quality int, width int) error {
	buffer, err := bimg.Read(inputPath)
	if err != nil {
		return err
	}

	options := bimg.Options{
		Quality: quality,
	}

	if width > 0 {
		options.Width = width
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return fmt.Errorf("failed to process image: %v", err)
	}

	err = bimg.Write(outputPath, newImage)
	if err != nil {
		return fmt.Errorf("failed to write output image: %v", err)
	}
	return nil
}
