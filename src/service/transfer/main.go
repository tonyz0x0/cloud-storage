package main

import (
	"bufio"
	"cloud-storage/config"
	dblayer "cloud-storage/db"
	"cloud-storage/src/mq"
	"cloud-storage/src/store/oss"
	"encoding/json"
	"log"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	// Parse message
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// get fd from temp local storage location
	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Upload file to OSS
	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin),
	)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Update Database
	_ = dblayer.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	return true
}

func main() {
	if !config.AsyncTransferEnable {
		log.Println("Asynchronous File Transfer Service is not activated, Please check configuration")
		return
	}
	log.Println("File Transfer Service is running, start to listen on transference job queue")
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer,
	)
}
