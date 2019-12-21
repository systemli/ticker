package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/systemli/ticker/internal/model"
	"github.com/systemli/ticker/internal/storage"
)

func TestFindTicker(t *testing.T) {
	setup()

	ticker := model.NewTicker()
	ticker.Domain = "localhost"
	_ = storage.DB.Save(ticker)

	ticker, err := storage.FindTicker("localhost")
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 1, ticker.ID)

	ticker, err = storage.FindTicker("example.com")
	if err == nil {
		t.Fail()
		return
	}
}
