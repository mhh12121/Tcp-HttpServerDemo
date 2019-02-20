package main

import (
	"fmt"
	"log"
	"net"
	"path"
	"runtime"

	"entry_task/Conf"
	dao "entry_task/DAO"
	data "entry_task/Data"
	service "entry_task/services"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
)

type TcpServer struct {
	Con net.Conn
}

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
	buff := make([]byte, 1024)
	// c := bufio.NewReader(conn)
	defer conn.Close()
	//this for loop is for one connection with multiple requests!!!
	for {

		size, cerr := conn.Read(buff)
		if cerr != nil {
			fmt.Println("buferr", cerr)
			panic(cerr)
		}

		// _, ioerr := io.ReadFull(c, buff[:int(size)])
		// if ioerr != nil {
		// 	fmt.Println(ioerr)
		// 	panic(ioerr)
		// }

		toServerD := &data.ToServerData{}
		dataErr := proto.Unmarshal(buff[:int(size)], toServerD)
		if dataErr != nil {
			fmt.Println("proto", dataErr)
			panic(dataErr)
		}

		//according to Ctype to set the response
		switch *toServerD.Ctype {
		case "login":
			tmpdata := &data.User{}
			tmperr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
			if tmperr != nil {
				fmt.Println("login err:", tmperr)
				panic(tmperr)
			}
			// tmpdata := data.Httpdata
			service.LoginHandle(conn, *tmpdata)
		//Home tcp
		case "home":
			tmpdata := &data.InfoWithUsername{}
			tmperr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
			if tmperr != nil {
				fmt.Println("login err:", tmperr)
				panic(tmperr)
			}

			log.Println("home tcp decode data", tmpdata)
			service.HomeHandle(conn, tmpdata.GetUsername(), tmpdata.GetToken())

		// case "uploadAvatar":
		// 	tmpdata := data.HttpData.(*Util.InfoWithUsername)
		// 	fmt.Println("tcp upload file decode data", tmpdata)
		// 	uploadHandle(conn, tmpdata.Username, tmpdata.Info, tmpdata.Token)

		// case "changeNickName":
		// 	tmpdata := data.HttpData.(*Util.InfoWithUsername)
		// 	fmt.Println("tcp change nickname decode data ", tmpdata)
		// 	changeNickNameHandle(conn, tmpdata.Username, tmpdata.Info, tmpdata.Token)

		case "logout":
			tmpdata := &data.InfoWithUsername{}
			tmperr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
			if tmperr != nil {
				fmt.Println("logout err:", tmperr)
				panic(tmperr)
			}
			fmt.Println("tcp change logout decode data ", tmpdata)
			// logoutHandle(conn, tmpdata.Username, tmpdata.Info)
			service.LogoutHandle(conn, tmpdata.GetUsername(), tmpdata.GetToken())
		}

	}
}

// func loginHandle(conn net.Conn, ruser data.User) {
// 	service.LoginHandle(conn, ruser)
// }

// func logoutHandle(conn net.Conn, username string, token interface{}) {
// 	service.LogoutHandle(conn, username, token)
//

// }
// func homeHandle(conn net.Conn, username string, token interface{}) {
// 	service.HomeHandle(conn, username, token)

// }

// //
// func uploadHandle(conn net.Conn, username string, avatar interface{}, token string) {
// 	service.UploadHandle(conn, username, avatar, token)
// 	// Util.FailSafeCheckErr("uploadfile encode err", errreturn)
// 	// return success
// }

// //
// func changeNickNameHandle(conn net.Conn, username string, nickname interface{}, token string) {
// 	service.ChangeNickNameHandle(conn, username, nickname, token)

// }
