package gin

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// If we're in production mode, set Gin to "release" mode
	if env := os.Getenv("ENV"); env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	return r
}
