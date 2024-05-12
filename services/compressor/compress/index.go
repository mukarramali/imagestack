package compress

import (
	"image"
	"image/jpeg"
	"os"
)

func CompressImage(inputPath string, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var opts jpeg.Options
	opts.Quality = 10
	return jpeg.Encode(outFile, img, &opts)
}
