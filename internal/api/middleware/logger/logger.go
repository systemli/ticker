package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func Logger(level string) gin.HandlerFunc {
	lvl, _ := logrus.ParseLevel(level)
	logger := logrus.New()
	logger.SetLevel(lvl)

	return ginlogrus.Logger(logger)
}
