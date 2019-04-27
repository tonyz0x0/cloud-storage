package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Print("Failed to get data, err: %s", err.Error())
			return
		}
		defer file.Close()

		newFile, err := os.Create("../tmp/" + head.Filename)
		if err != nil {
			fmt.Print("Failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Print("Failed to save data into file, err: %s", err.Error())
			return
		}

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

/**
 * Upload Success Handler
 */
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}
