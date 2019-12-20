package util

import (
	"io"
	"net/http"
)

//DetectContentType detects the ContentType from the first 512 bytes of the given io.Reader.
func DetectContentType(r io.Reader) string {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := r.Read(buffer)
	if err != nil {
		return "application/octet-stream"
	}

	return http.DetectContentType(buffer)
}
