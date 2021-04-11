package gin

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// If we're in production mode, set Gin to "release" mode
	if env := os.Getenv("ENV"); env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	r.Static("/", filepath.Join(dir, "./public"))
	// custom error
	r.NoRoute(func(c *gin.Context) {
		code := http.StatusNotFound
		c.JSON(code, gin.H{"code": code, "message": http.StatusText(code)})
	})

	// custom error
	r.NoMethod(func(c *gin.Context) {
		code := http.StatusMethodNotAllowed
		c.JSON(code, gin.H{"code": code, "message": http.StatusText(code)})
	})
	return r
}
