package util

import (
	"image"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
)

func ResizeImage(file io.Reader, maxDimension int) (image.Image, error) {
	img, err := imaging.Decode(file)
	if err != nil {
		log.WithError(err).Error("can't decode uploaded file")
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
	return imaging.Save(img, path, imaging.JPEGQuality(60), imaging.PNGCompressionLevel(png.BestCompression))
}
