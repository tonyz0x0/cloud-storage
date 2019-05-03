package db

import (
	"fmt"

	mydb "./mysql"
)

// UserSignup: User Sign Up
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`) values (?,?)",
	)

	if err != nil {
		fmt.Println("Failed to insert, err: " + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert err: " + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignin: Check User's Password is valid
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from tbl_user where user_name=? limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken: Refresh user sign in token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`, `user_token`) values(?,?)",
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
