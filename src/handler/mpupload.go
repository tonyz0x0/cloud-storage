package handler

import (
	"cloud-storage/src/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	rPool "cloud-storage/cache/redis"
	dblayer "cloud-storage/db"
)

// MultipartUploadInfo
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler : Multipart Upload Initialization
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Parameters
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// Get Redis Connection from Redis Pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Generate Initialization Info of Multipart Upload
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// Write Info into Redis
	// key, field, value
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// Write response back
	w.Write(util.NewRespMsg(
		0,
		"OK",
		upInfo).JSONBytes())
}

// UploadPartHandler: Upload Multi Part
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Parameters
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// Get Redis Connection from Redis Pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Get File Handler, store file conten
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(
			-1,
			"Upload part failed",
			nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil { // Exit when finished
			break
		}
	}

	// Update Redis Cache Record
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// Write response back
	w.Write(util.NewRespMsg(
		0,
		"OK",
		nil,
	).JSONBytes())
}

// CompleteUploadHandler
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Parameters
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// Get Redis Connection from Redis Pool
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// Check whether all parts have been uploaded
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(
			-1,
			"complete upload failed",
			nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}

	// TODO: Merge Multi Blocks

	// Update File Table and User File Table
	// FIXME: Empty FileAddress here
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	// Write Response Message
	w.Write(util.NewRespMsg(
		0,
		"OK",
		nil).JSONBytes())
}
