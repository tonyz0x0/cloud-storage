package route

import (
	"filestore-server/service/upload/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router Configuration
func Router() *gin.Engine {
	// gin framework, including Logger, Recovery
	router := gin.Default()

	// Static Resources
	router.Static("/static/", "./static")

	// support CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://localhost:8080"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))

	// File Normal Upload
	router.POST("/file/upload", api.DoUploadHandler)
	// File Fast Upload
	router.POST("/file/fastupload", api.TryFastUploadHandler)

	// File Multi-part Upload
	router.POST("/file/mpupload/init", api.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart", api.UploadPartHandler)
	router.POST("/file/mpupload/complete", api.CompleteUploadHandler)

	return router
}
