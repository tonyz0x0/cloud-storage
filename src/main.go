package main

import (
	"fmt"
	"net/http"

	"./handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("Failed to start server, err: %s", err.Error())
	}
}
