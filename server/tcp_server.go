package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"../DAO"
	"../Util"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

}

func main() {
	dao.InitDB()
	// dao.RedisInit()
	fmt.Println("tcp start ", Util.Tcpaddress)

	ln, err := net.Listen("tcp", ":"+Util.Tcpport)

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
		Util.FailSafeCheckErr("tcp decode err", err)
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
	//get remote Addr
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("tcp server connect:" + remoteAddr)
	//first go through redis cache
	//check if exists or different
	exists, errtoken := dao.CheckToken(ruser.Username, ruser.Token)
	Util.FailSafeCheckErr("login checktoken cache err", errtoken)
	//todo
	//some problems here(consistency)
	//1.checktoken in redis success then return success msg to http
	//2.http redirect to home
	//3.in the same time, the token in redis expires
	if exists {
		//if exists just take info from cache
		gob.Register(new(Util.ResponseFromServer))
		returnValue := Util.ResponseFromServer{Success: true, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(returnValue)
		Util.FailSafeCheckErr("login auth encode direct from cache err", errreturn)
		return
	}

	//check from mysql
	success, errorcheck := dao.Check(ruser.Username, ruser.Password)

	//login fail
	if !success || errorcheck != nil {
		fmt.Println("password wrong!")
		gob.Register(new(Util.ResponseFromServer))
		returnValue := Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(returnValue)
		Util.FailSafeCheckErr("login mysql encode err", errreturn)
		return
	}

	//if mysql check success, it will save it to redis as cache or update cache
	tokenerr := dao.SetToken(ruser.Username, ruser.Token, Util.TokenExpires)
	Util.FailSafeCheckErr("login save cache err", tokenerr)
	//login success
	log.Println("login handle tcp")
	gob.Register(new(Util.ResponseFromServer))
	returnValue := Util.ResponseFromServer{Success: true, TcpData: nil}
	encoder := gob.NewEncoder(conn)
	errreturn := encoder.Encode(returnValue)
	Util.FailSafeCheckErr("login auth encode err", errreturn)
	return
}

func logoutHandle(conn net.Conn, username string, token interface{}) {
	//invalid the cache in redis first
	err := dao.InvalidCache(username, token.(string))
	Util.FailSafeCheckErr("invalid logout usr", err)
	//return to logout handle
	success := (err == nil)
	gob.Register(new(Util.ResponseFromServer))
	returnValue := Util.ResponseFromServer{Success: success, TcpData: nil}
	encoder := gob.NewEncoder(conn)
	errreturn := encoder.Encode(returnValue)
	Util.FailSafeCheckErr("logout encode err", errreturn)

}
func homeHandle(conn net.Conn, username string, token interface{}) {
	//checktoken first
	exists, errtoken := dao.CheckToken(username, token.(string))
	Util.FailSafeCheckErr("home checktoken cache err", errtoken)
	//1. cookie still exists but token expires
	//---solution: clear cookie first then redirect to login
	//2. cookie expires but token exists
	//---solution: login and refresh the token

	//token not exists or not correct
	if !exists || errtoken != nil {
		gob.Register(new(Util.ResponseFromServer))
		returnValue := Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(returnValue)
		Util.FailSafeCheckErr("home auth encode direct from cache err", errreturn)
		return
	}

	//First go through the Redis get cache
	user, ok, err := dao.GetCacheInfo(username)
	Util.FailSafeCheckErr("redis get cache fail:", err)
	//cache still valid
	if ok {
		log.Println("tcp home handle cache get info okay", user)
		gob.Register(new(Util.RealUser))
		tohttp := &Util.ResponseFromServer{Success: true, TcpData: user}

		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(tohttp)
		Util.FailSafeCheckErr("home handle encode fail", errreturn)
		return
	}

	//cache expires or not exists then go to mysql
	userdb, okdb := dao.AllInfo(username)
	//retrieve from mysql success
	if okdb {
		//it will also save it to cache
		successCache := dao.SaveCacheInfo(username, userdb.Nickname, userdb.Avatar)
		if !successCache {
			fmt.Println("update redis homne cache fail")
		}

		//save cache success
		//here how
		// if successCache {
		gob.Register(new(Util.RealUser))
		tohttp := &Util.ResponseFromServer{Success: true, TcpData: userdb}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(tohttp)
		Util.FailSafeCheckErr("home handle encode err", errreturn)
		return
		// }
		//todo
		// //save cache fail,but still need to send back
		// gob.Register(new(Util.RealUser))
		// tohttp := &Util.ResponseFromServer{Success: false, TcpData: userdb}
		// encoder := gob.NewEncoder(conn)
		// errreturn := encoder.Encode(tohttp)
		// Util.FailSafeCheckErr("home handle encode err", errreturn)

	}

	return

}

//
func uploadHandle(conn net.Conn, username string, avatar interface{}, token string) {
	exists, errtoken := dao.CheckToken(username, token)
	Util.FailSafeCheckErr("upload checktoken cache err", errtoken)

	//token not exists or not correct
	if !exists || errtoken != nil {
		gob.Register(new(Util.ResponseFromServer))
		returnValue := Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(returnValue)
		Util.FailSafeCheckErr("home auth encode direct from cache err", errreturn)
		return
	}
	success := dao.UpdateAvatar(username, "/"+avatar.(string))
	//update mysql success
	if success {
		//todo
		//update cache
		err := dao.UpdateCacheAvatar(username, avatar.(string))
		if err != nil {
			//update cache fail
			fmt.Println("update avatar redis cache fail", err)
			//do nothing
			// return
		}

		// gob.Register(new(Util.ResponseFromServer))
		// tohttp := &Util.ResponseFromServer{Success: success, TcpData: nil}
		// encoder := gob.NewEncoder(conn)
		// errreturn := encoder.Encode(tohttp)
		// Util.FailSafeCheckErr("uploadfile encode err", errreturn)
	}
	//mysql update not success
	gob.Register(new(Util.ResponseFromServer))
	tohttp := &Util.ResponseFromServer{Success: success, TcpData: nil}
	encoder := gob.NewEncoder(conn)
	errreturn := encoder.Encode(tohttp)
	Util.FailSafeCheckErr("uploadfile encode err", errreturn)
	// return success
}

//
func changeNickNameHandle(conn net.Conn, username string, nickname interface{}, token string) {
	//first check token from redis frist
	exists, errtoken := dao.CheckToken(username, token)
	Util.FailSafeCheckErr("updatenickname checktoken cache err", errtoken)
	if !exists || errtoken != nil {
		//token expires or not correct
		gob.Register(new(Util.ResponseFromServer))
		tohttp := &Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(tohttp)
		Util.FailSafeCheckErr("changenickname encode err", errreturn)
		return
	}
	//update mysql first
	success, errorupdate := dao.UpdateNickname(username, nickname.(string))
	if success && errorupdate == nil {
		//update cache
		//if successfully change data in mysql
		err := dao.UpdateCacheNickname(username, nickname.(string))
		if err != nil {
			fmt.Println("update nickname fail", err)
			//update cache fail
			//todo
			//do nothing
			// return
		}

	}
	gob.Register(new(Util.ResponseFromServer))
	tohttp := &Util.ResponseFromServer{Success: success, TcpData: nil}
	encoder := gob.NewEncoder(conn)
	errreturn := encoder.Encode(tohttp)
	Util.FailSafeCheckErr("changenickname encode err", errreturn)
}
