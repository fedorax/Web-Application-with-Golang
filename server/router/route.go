package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) *gin.Engine {

	// Ping test
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusText(http.StatusOK), "message": "pong"})
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := app.Group("/api", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")
		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		//user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		//var json struct {
		//	Value string `json:"value" binding:"required"`
		//}

		//if c.Bind(&json) == nil {
		//	db[user] = json.Value
		//	c.JSON(http.StatusOK, gin.H{"status": "ok"})
		//}
	})

	return app
}
