package route

import (
	"cloud-storage/src/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")

	// Sign Up/In
	router.GET("/user/signup", handler.SignUpHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	router.GET("/user/signin", handler.SignInHandler)
	router.POST("/user/signin", handler.DoSignInHandler)

	router.Use(handler.HTTPInterceptor())

	// File Interface
	// router.GET("/file/upload", handler.UploadHandler)
	// router.POST("/file/upload", handler.DoUploadHandler)
	// router.GET("/file/upload/suc", handler.UploadSucHandler)
	// router.POST("/file/meta", handler.GetFileMetaHandler)
	// router.POST("/file/query", handler.FileQueryHandler)
	// router.GET("/file/download", handler.DownloadHandler)
	// router.POST("/file/update", handler.FileMetaUpdateHandler)
	// router.POST("/file/delete", handler.FileDeleteHandler)
	// router.POST("/file/downloadurl", handler.DownloadURLHandler)

	// XSpeed Upload
	// router.POST("/file/fastupload", handler.TryFastUploadHandler)

	// Multipart Upload
	// router.POST("/file/mpupload/init", handler.InitialMultipartUploadHandler)
	// router.POST("/file/mpupload/uppart", handler.UploadPartHandler)
	// router.POST("/file/mpupload/complete", handler.CompleteUploadHandler)

	// User Interface
	// router.POST("/user/info", handler.UserInfoHandler)
	return router
}
