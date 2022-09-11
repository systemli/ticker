package storage

import (
	"testing"
)

func TestOpenDBPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	OpenDB("/x/y/z")
}
