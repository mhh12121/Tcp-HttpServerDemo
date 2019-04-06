package dao

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"runtime"

	"entry_task/Conf"
	data "entry_task/Data"

	"github.com/golang/protobuf/proto"
)

var db *sql.DB
var ad string

func init() {
	_, filepath, _, _ := runtime.Caller(0)
	p := path.Dir(filepath)
	p = path.Dir(p)

	log.Println("log path", p)

	Conf.LoadConf(p + "/Conf/config.json")

}

//

func InitDB() {
	ad = Conf.Config.Mysql.Username + ":" + Conf.Config.Mysql.Password + "@" +
		"tcp(" + Conf.Config.Mysql.Host + ":" + Conf.Config.Mysql.Port + ")/" +
		Conf.Config.Mysql.Db + "?charset=utf8"
	var err error
	log.Println("das", ad)
	db, err = sql.Open("mysql", ad)
	checkErr(err)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
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
func AllInfo(username string) (*data.RealUser, bool) {

	fmt.Println("------------------dao home get info--------------------------------")
	c, err := db.Prepare("select nickname,avatar from user where username = ? ")
	if err != nil {
		log.Println("mysql db get allinfo err", err)
		return nil, false
	}
	// checkErr(err)
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
	return &data.RealUser{Username: proto.String(username), Nickname: proto.String(nickname), Avatar: proto.String(avatar)}, true

	// fmt.Println("mysql get allinfo failed")
	// return nil, false
}

//update nickname
func UpdateNickname(username, nickname string) (bool, error) {
	log.Println("dao.updatenickname", username, nickname)
	c, err := db.Prepare("update user SET nickname = ? where username = ?")
	if err != nil {
		log.Println("mysql db updatenickname err", err)
		return false, err
	}
	defer c.Close()
	_, err = c.Exec(nickname, username)
	if err != nil {
		fmt.Println("mysql update nickname fail", err)
		return false, err
	}

	return true, nil
}
func UpdateAvatar(username, avatar string) bool {
	c, err := db.Prepare("update user SET avatar = ? where username = ?")
	if err != nil {
		fmt.Println("mysql db update avatar err:", err)
	}
	defer c.Close()
	_, err = c.Exec(avatar, username)
	if err != nil {
		fmt.Println("mysql update avatar fail", err)
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
