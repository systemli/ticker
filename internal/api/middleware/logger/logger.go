package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return ginlogrus.Logger(logger)
}
