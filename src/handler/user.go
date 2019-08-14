package handler

import (
	"cloud-storage/src/common"
	"fmt"
	"net/http"
	"time"

	dblayer "cloud-storage/db"

	"cloud-storage/src/util"

	"github.com/gin-gonic/gin"
)

const (
	pwdSalt = "#739"
)

// SignUpHandler: Handle User Sign Up Get Request
func SignUpHandler(c *gin.Context) {
	// TODO: Imporove the front end
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignupHandler: Post Reuqest
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Invalid Parameters",
			"code": common.StatusParamInvalid,
		})
		return
	}

	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Sign Up Success",
			"code": common.StatusOK,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Sign Up Fail",
			"code": common.StatusRegisterFailed,
		})
	}
}

// SignInHandler: Handle User Sign In Get Request
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func DoSignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	encPasswd := util.Sha1([]byte(password + pwdSalt))

	// 1. Validate Password
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Sign In Fail",
			"code": common.StatusLoginFailed,
		})
		return
	}

	// 2. Generate Token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Sign In Fail",
			"code": common.StatusLoginFailed,
		})
		return
	}

	// 3. Redirect to Home Page
	resp := util.RespMsg{
		Code: int(common.StatusOK),
		Msg:  "Sign In Success",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// UserInfoHandler: User Info
func UserInfoHandler(c *gin.Context) {
	// Parse the Request Parameters
	username := c.Request.FormValue("username")

	// Query User Info
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{})
		return
	}

	// Response with User Info
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// UserExistsHandler: Check UserName is existed or not
func UserExistsHandler(c *gin.Context) {
	// Parse the Request Parameters
	username := c.Request.FormValue("username")

	// Get Information
	exists, err := dblayer.UserExist(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"code": common.StatusServerError,
				"msg":  "server error",
			})
	} else {
		c.JSON(http.StatusOK,
			gin.H{
				"code":   common.StatusOK,
				"msg":    "ok",
				"exists": exists,
			})
	}
}

// GenToken
func GenToken(username string) string {
	// Token Rule: md5(username + timestamp + token_salt)+timestamp[:8] -> 40 bytes
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid: Verify the token is valid or not
func IsTokenValid(token string) bool {
	// TODO: Verify Token's validity period
	// Get Token from Database and check whether they are the same
	if len(token) != 40 {
		return false
	}
	return true
}
