package handler

import (
	cfg "cloud-storage/config"
	dblayer "cloud-storage/db"
	"cloud-storage/src/meta"
	"cloud-storage/src/store/ceph"
	"cloud-storage/src/store/oss"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// DownloadURLHandler: Generate General Download URL
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	filehash := r.Form.Get("filehash")
	// Get Record from File Table
	row, _ := dblayer.GetFileMeta(filehash)

	// Find where is the file, in local storage, Ceph or OSS
	if strings.HasPrefix(row.FileAddr.String, cfg.TempLocalRootDir) ||
		strings.HasPrefix(row.FileAddr.String, "/ceph") {
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		tmpURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			r.Host, filehash, username, token)
		w.Write([]byte(tmpURL))
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// Generate OSS Download URL
		signedURL := oss.DownloadURL(row.FileAddr.String)
		w.Write([]byte(signedURL))
	}
}

// DownloadHandler : Normal Download From Local or Ceph
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fileData []byte
	fmt.Println(fm.Location)
	if strings.HasPrefix(fm.Location, cfg.TempLocalRootDir) {
		f, err := os.Open(fm.Location)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fileData, err = ioutil.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fm.Location, "/ceph") {
		fmt.Println("Download file from ceph...")
		bucket := ceph.GetCephBucket("userfile")
		fileData, err = bucket.Get(fm.Location)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	w.Write(fileData)
}

// RangeDownloadHandler : Breakpoint Download
func RangeDownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	http.ServeFile(w, r, fm.Location)
}
