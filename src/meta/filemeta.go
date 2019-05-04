package meta

import (
	"fmt"
	"sort"

	mydb "cloud-storage/db"
)

// FileMeta: File Meta Data
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: Create/Update FileMeta
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB: Upload File Meta to Database
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1,
		fmeta.FileName,
		fmeta.FileSize,
		fmeta.Location,
	)
}

// GetFileMeta: Get FileMeta by sha1
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB: Get FileMeta by sha1 from Database
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

// GetLastFileMetas: Get batch FileMeta of number K
func GetLastFileMetas(count int) []FileMeta {
	// sort so the updatest will be on the first
	fmt.Println(fileMetas)
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	// make sure the fMetaArray is sorted by UploadAt
	// Customized sort rule
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

// TODO: GetLastFileMetasDB: Get batch FileMeta of number K from DB
func GetLastFileMetasDB(count int) []FileMeta {
	return nil
}

func RemoveFileMeta(fileSha1 string) {
	// TODO: If multi thread, need to consider thread-safe action, probably introduce Lock
	delete(fileMetas, fileSha1)
}
