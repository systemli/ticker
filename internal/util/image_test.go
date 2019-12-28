package util_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/util"
)

func TestResizeImage(t *testing.T) {
	file, err := os.Open("../../testdata/gopher.jpg")
	if err != nil {
		t.Fail()
		return
	}

	img, err := util.ResizeImage(file, 100)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 63, img.Bounds().Dy())
	assert.Equal(t, 100, img.Bounds().Dx())

	r := bytes.NewReader([]byte{})
	img, err = util.ResizeImage(r, 100)
	if err == nil {
		t.Fail()
	}
}

func TestSaveImage(t *testing.T) {
	img, err := imaging.Open("../../testdata/gopher.jpg")
	if err != nil {
		t.Fail()
		return
	}

	err = util.SaveImage(img, fmt.Sprintf("%s/%d.jpg", os.TempDir(), time.Now().Nanosecond()))
	if err != nil {
		t.Fail()
	}
}
