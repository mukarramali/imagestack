package compress

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func CompressImage(inputPath string, outputPath string, quality int) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)

	if err != nil {
		fmt.Printf("image compression failed for format %s with error:%s", format, err)
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var opts jpeg.Options
	opts.Quality = quality
	return jpeg.Encode(outFile, img, &opts)
}
