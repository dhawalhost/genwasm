package routes

import (
	"github.com/dhawalhost/genwasm/services"
	"github.com/gin-gonic/gin"
)

// InitRoutes -
func InitRoutes(g *gin.Engine) {
	g.GET("/getwasm/:filename", func(c *gin.Context) {
		// services.Display(c.Writer, "upload", nil)
		services.DownloadWasmFile(c)
	})
	// g.GET("/", func(c *gin.Context) {
	// 	services.Display(c.Writer, "upload", nil)
	// })
	g.POST("/buildwasm", func(c *gin.Context) {
		services.UploadFile(c)
	})
}
