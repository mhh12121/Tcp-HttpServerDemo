package dao

import (
	"database/sql"
	"fmt"
	"log"

	"../Util"
)

type dbinstance struct {
	address string
}

var db *sql.DB
var ad = "root:12345678@/entrytask?charset=utf8"

func InitDB() {
	var err error
	db, err = sql.Open("mysql", ad)
	checkErr(err)
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(300)
	// defer db.Close()
	err = db.Ping()
	checkErr(err)
}

//Check exported
//check if pwd is correct
func Check(username, password string) (bool, error) {
	fmt.Println("mysqldao check")
	c, err := db.Prepare("select password from user where username = ? ")
	if err != nil {
		return false, err
	}
	defer c.Close()

	var tmppwd string

	err = c.QueryRow(username).Scan(&tmppwd)
	if err != nil {
		return false, err
	}

	if tmppwd == password && tmppwd != "" {
		fmt.Println("done!!pwd:", tmppwd)
		return true, nil
	}
	return false, nil

}

//A user's info
func AllInfo(username string) (*Util.RealUser, bool) {

	fmt.Println("dao home get info")
	c, err := db.Prepare("select nickname,avatar from user where username = ? ")
	checkErr(err)
	defer c.Close()
	var (
		nickname string
		avatar   string
	)
	errquery := c.QueryRow(username).Scan(&nickname, &avatar)
	if errquery != nil {
		fmt.Println("mysql get allinfo failed", errquery)
		return nil, false
	}

	//then return
	return &Util.RealUser{Username: username, Nickname: nickname, Avatar: avatar}, true

	// fmt.Println("mysql get allinfo failed")
	// return nil, false
}

//update nickname
func UpdateNickname(username, nickname string) (bool, error) {
	log.Println("dao.updatenickname", username, nickname)
	c, err := db.Prepare("update user SET nickname = ? where username = ?")
	if err != nil {
		return false, err
	}
	defer c.Close()
	_, err = c.Exec(nickname, username)
	if err != nil {
		fmt.Println("update nickname fail", err)
		return false, err
	}

	return true, nil
}
func UpdateAvatar(username, avatar string) bool {
	c, err := db.Prepare("update user SET avatar = ? where username = ?")
	if err != nil {
		fmt.Println("update avatar sql db connect:", err)
	}
	defer c.Close()
	_, err = c.Exec(avatar, username)
	if err != nil {
		fmt.Println("update avatar fail", err)
		return false
	}

	return true
}
func checkErr(err error) {
	if err != nil {

		fmt.Println("sql get wrong", err)

		// panic(err)
	}
}
