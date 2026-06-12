package application

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CollectByLinkedin handles job collection from linkedin.
func CollectByLinkedin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info("here you go")
	}
}
