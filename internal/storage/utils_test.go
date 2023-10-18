package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenGormDB(t *testing.T) {
	_, err := OpenGormDB("sqlite", "file::memory:?cache=shared", nil)
	assert.NoError(t, err)
}

func TestDialector(t *testing.T) {
	var tests = []struct {
		dbType     string
		dsn        string
		shouldFail bool
	}{
		{"sqlite", "file::memory:?cache=shared", false},
		{"mysql", "user:password@tcp(localhost:5555)/dbname?charset=utf8mb4&parseTime=True&loc=Local", false},
		{"postgres", "host=myhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai", false},
		{"", "", true},
	}

	for _, test := range tests {
		d := dialector(test.dbType, test.dsn)
		if test.shouldFail {
			assert.Nil(t, d)
		} else {
			assert.NotNil(t, d)
		}
	}
}
