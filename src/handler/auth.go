package handler

import (
	"cloud-storage/src/common"
	"cloud-storage/src/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		// Verify the token is valid or not
		if len(username) < 3 || !IsTokenValid(token) {
			// Token is invalid, return error response message
			c.Abort()
			resp := util.NewRespMsg(
				int(common.StatusTokenInvalid),
				"Invalid Token",
				nil,
			)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}
