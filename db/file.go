package db

import (
	"database/sql"
	"fmt"

	mydb "cloud-storage/db/mysql"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// OnFileUploadFinished: File upload completed
func OnFileUploadFinished(
	filehash string,
	filename string,
	filesize int64,
	fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("\nFile with hash:%s has been uploaded previously", filehash)
		}
		return true
	}
	return false
}

// GetFileMeta: Query File Meta Data from Database
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1",
	)

	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}
	stmt.QueryRow(filehash).Scan(
		&tfile.FileHash,
		&tfile.FileAddr,
		&tfile.FileName,
		&tfile.FileSize,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Can not find record in DB
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	defer stmt.Close()

	return &tfile, nil
}

// IsFileUploaded: File is uploaded or not
func IsFileUploaded(filehash string) bool {
	// the SQL is just to check whether there has been a record there or not
	stmt, err := mydb.DBConn().Prepare(
		"select 1 from tbl_file where file_sha1=? and status=1 limit 1",
	)
	rows, err := stmt.Query(filehash)
	if err != nil || rows == nil || rows.Next() {
		return false
	}
	return true
}

// GetFileMetaList: Query Batch File Metas from DB
func GetFileMetaList(limit int) ([]TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size from tbl_file " +
			"where status = 1 limit ?",
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	var tfiles []TableFile
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr,
			&tfile.FileName, &tfile.FileSize)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	fmt.Println(len(tfiles))
	return tfiles, nil
}

// OnFileRemoved: Delete file in Database(only marked as deleted, not real deleted)
func OnFileRemoved(filehash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_file set status=2 where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s not uploaded", filehash)
		}
		return true
	}
	return false
}
