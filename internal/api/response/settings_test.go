package response

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/systemli/ticker/internal/storage"
)

type SettingsResponseTestSuite struct {
	suite.Suite
}

func (s *SettingsResponseTestSuite) TestInactiveSettingsResponse() {
	inactiveSettings := storage.DefaultInactiveSettings()

	setting := InactiveSettingsResponse(inactiveSettings)

	s.Equal(storage.SettingInactiveName, setting.Name)
	s.Equal(inactiveSettings, setting.Value)
}

func TestSettingsResponseTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsResponseTestSuite))
}
