package util

import (
	"image"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
)

func ResizeImage(file io.Reader, maxDimension int) (image.Image, error) {
	img, err := imaging.Decode(file)
	if err != nil {
		return img, err
	}
	if img.Bounds().Dx() > maxDimension {
		img = imaging.Resize(img, maxDimension, 0, imaging.Linear)
	}
	if img.Bounds().Dy() > maxDimension {
		img = imaging.Resize(img, 0, maxDimension, imaging.Linear)
	}

	return img, nil
}

func SaveImage(img image.Image, path string) error {
	opts := []imaging.EncodeOption{
		imaging.JPEGQuality(60),
		imaging.PNGCompressionLevel(png.BestCompression),
	}

	return imaging.Save(img, path, opts...)
}
