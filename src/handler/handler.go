package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

/**
 * Handle upload
 */
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// return upload page
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// acept file stream and store to local
	}
}
