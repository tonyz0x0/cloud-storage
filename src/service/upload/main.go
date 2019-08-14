package main

import (
	"fmt"

	cfg "cloud-storage/config"
	"cloud-storage/src/route"
)

func main() {
	router := route.Router()
	fmt.Println("Upload Server is running, listening on port 8080...")
	router.Run(cfg.UploadServiceHost)
}
