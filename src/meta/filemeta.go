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

/**
 * In-Memory Operations
 */

// UpdateFileMeta: Create/Update FileMeta
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMeta: Get FileMeta by sha1
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
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

func RemoveFileMeta(fileSha1 string) {
	// TODO: If multi thread, need to consider thread-safe action, probably introduce Lock
	delete(fileMetas, fileSha1)
}

/**
 * Database Operations
 */

// UpdateFileMetaDB: Upload File Meta to Database
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1,
		fmeta.FileName,
		fmeta.FileSize,
		fmeta.Location,
	)
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

// GetLastFileMetasDB: Get batch FileMeta of number K from DB
func GetLastFileMetasDB(limit int) ([]FileMeta, error) {
	tfiles, err := mydb.GetFileMetaList(limit)
	if err != nil {
		return make([]FileMeta, 0), err
	}

	tfilesm := make([]FileMeta, len(tfiles))
	for i := 0; i < len(tfilesm); i++ {
		tfilesm[i] = FileMeta{
			FileSha1: tfiles[i].FileHash,
			FileName: tfiles[i].FileName.String,
			FileSize: tfiles[i].FileSize.Int64,
			Location: tfiles[i].FileAddr.String,
		}
	}
	return tfilesm, nil
}

// OnFileRemoveDB: Delete File from Database
func OnFileRemoveDB(filehash string) bool {
	return mydb.OnFileRemoved(filehash)
}
