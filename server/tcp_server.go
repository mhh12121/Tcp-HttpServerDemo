package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"path"
	"runtime"

	"../Conf"
	"../DAO"
	"../Util"
	service "../services"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	_, filepath, _, _ := runtime.Caller(0)
	p := path.Dir(filepath)
	p = path.Dir(p)

	log.Println("log path", p)

	Conf.LoadConf(p + "/Conf/config.json")
	// log.Println("dafas", Conf.Config)
}
func main() {
	dao.InitDB()
	// dao.RedisInit()
	fmt.Println("tcp start ", Conf.Config.Connect.Tcphost)
	fmt.Println("tcp start ", Conf.Config.Connect.Tcpport)
	ln, err := net.Listen("tcp", ":"+Conf.Config.Connect.Tcpport)

	if err != nil {
		fmt.Println("tcp listen failed:", err)
	}
	// defer ln.Close()
	//need keep connection

	//keep listening for multiple connections(clients)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("tcp server connection failed:", err)
			continue
		}
		log.Println("tcp listen loop", conn)
		go handleAll(conn)

	}

}

func handleAll(conn net.Conn) {
	//this for loop is for one connection with multiple requests!!!
	for {
		gob.Register(new(Util.RealUser))
		gob.Register(new(Util.User))
		gob.Register(new(Util.ToServerData))
		gob.Register(new(Util.InfoWithUsername))
		// gob.Register(new(Util.Avatar))
		//Decoder blocks here???
		decoder := gob.NewDecoder(conn)
		var data Util.ToServerData
		err := decoder.Decode(&data)
		if err != nil {
			log.Println("tcp handle all decode err", err)
		}
		// Util.FailSafeCheckErr("tcp decode err", err)
		log.Println("tcp decode", data)

		//according to Ctype to diffentiate the response
		switch data.Ctype {
		case "login":
			tmpdata := data.HttpData.(*Util.User)
			log.Println("login tcp decode data", tmpdata)
			loginHandle(conn, tmpdata)
		case "home":
			tmpdata := data.HttpData.(*Util.InfoWithUsername)
			log.Println("home tcp decode data", tmpdata)
			homeHandle(conn, tmpdata.Username, tmpdata.Info)
		case "uploadAvatar":
			tmpdata := data.HttpData.(*Util.InfoWithUsername)
			fmt.Println("tcp upload file decode data", tmpdata)
			uploadHandle(conn, tmpdata.Username, tmpdata.Info, tmpdata.Token)

		case "changeNickName":
			tmpdata := data.HttpData.(*Util.InfoWithUsername)
			fmt.Println("tcp change nickname decode data ", tmpdata)
			changeNickNameHandle(conn, tmpdata.Username, tmpdata.Info, tmpdata.Token)

		case "logout":
			tmpdata := data.HttpData.(*Util.InfoWithUsername)
			fmt.Println("tcp change logout decode data ", tmpdata)
			logoutHandle(conn, tmpdata.Username, tmpdata.Info)
		}

	}
}

func loginHandle(conn net.Conn, ruser *Util.User) {
	service.LoginHandle(conn, ruser)
}

func logoutHandle(conn net.Conn, username string, token interface{}) {
	service.LogoutHandle(conn, username, token)
	// Util.FailSafeCheckErr("logout encode err", errreturn)

}
func homeHandle(conn net.Conn, username string, token interface{}) {
	service.HomeHandle(conn, username, token)

}

//
func uploadHandle(conn net.Conn, username string, avatar interface{}, token string) {
	service.UploadHandle(conn, username, avatar, token)
	// Util.FailSafeCheckErr("uploadfile encode err", errreturn)
	// return success
}

//
func changeNickNameHandle(conn net.Conn, username string, nickname interface{}, token string) {
	service.ChangeNickNameHandle(conn, username, nickname, token)

}
