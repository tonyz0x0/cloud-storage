package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	dblayer "cloud-storage/db"

	"cloud-storage/src/util"
)

const (
	pwdSalt = "#739"
)

// SignUpHandler: Handle User Sign Up Request
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// TODO: Imporove the front end
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()

	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SignInHandler: Sign In Handler
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		// http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(passwd + pwdSalt))

	// Validate Password
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	// Generate Token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	// Redirect to Home Page
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the Request Parameters
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")

	// Verify Token
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		w.WriteHeader(http.StatusForbidden)
	}

	// Query User Info
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Response with User Info
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
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
