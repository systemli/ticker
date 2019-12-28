package util_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/util"
)

func TestDetectContentTypeImage(t *testing.T) {
	file, err := os.Open("../../testdata/gopher.jpg")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, "image/jpeg", util.DetectContentType(file))
}

func TestDetectContentTypeOther(t *testing.T) {
	r := strings.NewReader("content")

	assert.Equal(t, "application/octet-stream", util.DetectContentType(r))
}
