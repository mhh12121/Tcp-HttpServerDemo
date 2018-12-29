package service

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	dao "entry_task/DAO"
	"entry_task/Util"
)

//tcp handle
func LoginHandle(conn net.Conn, ruser *Util.User) {
	//get remote Addr
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("tcp server connect:" + remoteAddr)
	//first go through redis cache
	//check if exists or different
	//what if login in another device?
	exists, errtoken := dao.CheckToken(ruser.Username, ruser.Token)
	if errtoken != nil {
		log.Println("login checktoken cache err", errtoken)
	}
	// Util.FailSafeCheckErr("login checktoken cache err", errtoken)
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
		if errreturn != nil {
			log.Println("login auth encode direct from cache err", errreturn)
		}
		// Util.FailSafeCheckErr("login auth encode direct from cache err", errreturn)
		return
	}

	//check from mysql
	success, errorcheck := dao.Check(ruser.Username, ruser.Password)

	//login fail
	if !success || errorcheck != nil {
		log.Println("password wrong! any error?:", errorcheck)
		gob.Register(new(Util.ResponseFromServer))
		returnValue := Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(returnValue)
		if errreturn != nil {
			log.Println("login mysql encode err", errreturn)
		}
		// Util.FailSafeCheckErr("login mysql encode err", errreturn)
		return
	}

	//if mysql check success, it will save it to redis as cache or update cache
	tokenerr := dao.SetToken(ruser.Username, ruser.Token, Util.TokenExpires)
	if tokenerr != nil {
		log.Println("login save cache err", tokenerr)
	}
	// Util.FailSafeCheckErr("login save cache err", tokenerr)
	//login success
	log.Println("login handle tcp")
	gob.Register(new(Util.ResponseFromServer))
	returnValue := Util.ResponseFromServer{Success: true, TcpData: nil}
	encoder := gob.NewEncoder(conn)
	errreturn := encoder.Encode(returnValue)
	if errreturn != nil {
		log.Println("login auth encode err", errreturn)
	}
	// Util.FailSafeCheckErr("login auth encode err", errreturn)
	return
}
