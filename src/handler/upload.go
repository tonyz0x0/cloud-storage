package handler

import (
	"bytes"
	"cloud-storage/src/mq"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	cfg "cloud-storage/config"
	dblayer "cloud-storage/db"
	cmn "cloud-storage/src/common"
	"cloud-storage/src/meta"
	"cloud-storage/src/store/ceph"
	"cloud-storage/src/store/oss"
	"cloud-storage/src/util"

	"github.com/gin-gonic/gin"
)

func init() {
	// Directory is existed
	if _, err := os.Stat(cfg.TempLocalRootDir); err == nil {
		return
	}

	// Try to create directory
	err := os.MkdirAll(cfg.TempLocalRootDir, 0744)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

// UploadHandler
func UploadHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./static/view/index.html")
	if err != nil {
		c.String(404, `Page Not Found`)
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

// DoUploadHandler
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "Upload failed",
			})
		}
	}()
	// acept file stream and store to local
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err: %s", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	// Convert File Content to byte array
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Printf("Failed to get file data, err:%s\n", err.Error())
		errCode = -2
		return
	}

	// Construct File Meta Info
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		FileSha1: util.Sha1(buf.Bytes()),
		FileSize: int64(len(buf.Bytes())),
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Write File Content to temp local storage location
	fileMeta.Location = cfg.TempLocalRootDir + fileMeta.FileSha1
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		errCode = -3
		return
	}
	defer newFile.Close()

	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != fileMeta.FileSize || err != nil {
		fmt.Printf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error())
		errCode = -4
		return
	}

	// Store into Ceph/OSS
	newFile.Seek(0, 0)
	if cfg.CurrentStoreType == cmn.StoreCeph {
		// Write File into Ceph
		data, _ := ioutil.ReadAll(newFile)
		cephPath := "/ceph/" + fileMeta.FileSha1
		_ = ceph.PutObject("userfile", cephPath, data)
		fileMeta.Location = cephPath
	} else if cfg.CurrentStoreType == cmn.StoreOSS {
		ossPath := "oss/" + fileMeta.FileSha1
		if !cfg.AsyncTransferEnable {
			// Write File into OSS Synchronously
			err = oss.Bucket().PutObject(ossPath, newFile)
			if err != nil {
				fmt.Println(err.Error())
				errCode = -5
				return
			}
			fileMeta.Location = ossPath
		} else {
			data := mq.TransferData{
				FileHash:      fileMeta.FileSha1,
				CurLocation:   fileMeta.Location,
				DestLocation:  ossPath,
				DestStoreType: cmn.StoreOSS,
			}
			pubData, _ := json.Marshal(data)
			pubSuc := mq.Publish(
				cfg.TransExchangeName,
				cfg.TransOSSRoutingKey,
				pubData,
			)
			if !pubSuc {
				// TODO: Fail to Publish, Retry later
			}
		}
	}

	fmt.Printf("\nNew File %s Created in %s", fileMeta.FileSha1, fileMeta.Location)
	// meta.UpdateFileMeta(fileMeta)
	_ = meta.UpdateFileMetaDB(fileMeta)

	//Update User File Table in Database
	username := c.Request.FormValue("username")
	suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
		fileMeta.FileName, fileMeta.FileSize)
	if suc {
		c.Redirect(http.StatusFound, "/static/view/home.html")
	} else {
		errCode = -6
	}
}

// UploadSuccessHandler
func UploadSuccessHandler(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "Upload Finish!",
		})
}

// GetFileMetaHandler: get file meta data
func GetFileMetaHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"code": -2,
				"msg":  "Upload failed!",
			})
		return
	}

	if err != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"code": -3,
					"msg":  "Upload failed!",
				})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	} else {
		c.JSON(http.StatusOK,
			gin.H{
				"code": -4,
				"msg":  "No such file",
			})
	}
}

// FileQueryHandler: Query Batch File Metas
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))

	// fileMetas := meta.GetLastFileMetas(limitCnt)
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileMetaUpdateHandler: Update FileName
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")
	newFileName := r.Form.Get("filename")

	if opType != "0" || len(newFileName) < 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// curFileMeta := meta.GetFileMeta(fileSha1)
	// curFileMeta.FileName = newFileName
	// meta.UpdateFileMeta(curFileMeta)

	// Update file name in User File Table, File Table does not need to change
	_ = dblayer.RenameFileName(username, fileSha1, newFileName)

	// Get newest File Meta
	userFile, err := dblayer.QueryUserFileMeta(username, fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	fileSha1 := r.Form.Get("filehash")

	// fMeta := meta.GetFileMeta(fileSha1)
	fMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Remove(fMeta.Location) // Omit the failure of deleting file

	// meta.RemoveFileMeta(fileSha1)
	suc := dblayer.DeleteUserFile(username, fileSha1)
	if !suc {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TryFastUploadHandler: Try XSpeed Upload
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//Parse Parameters
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	//Query File Records with Same Hash Value From File Table
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//If records not found, failed to XSpeed Upload and go to normal upload
	if fileMeta.FileSha1 == "" {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "Failed to XSpeed Upload, redirect to Normal Upload",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// If records found, XSpeed Upload is activated, update records in User File Table
	// return Success
	suc := dblayer.OnUserFileUploadFinished(
		username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "XSpeed Upload SUCCESS",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "Failed to XSpeed Upload，pleas try later",
	}
	w.Write(resp.JSONBytes())
	return
}
