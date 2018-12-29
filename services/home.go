package service

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	dao "entry_task/DAO"
	"entry_task/Util"
)

func HomeHandle(conn net.Conn, username string, token interface{}) {
	//checktoken first
	exists, errtoken := dao.CheckToken(username, token.(string))
	// Util.FailSafeCheckErr("home checktoken cache err", errtoken)
	//1. cookie still exists but token expires
	//---solution: clear cookie first then redirect to login
	//2. cookie expires but token exists
	//---solution: login and refresh the token

	//token not exists or not correct
	if !exists || errtoken != nil {
		log.Println("home checktoken cache err", errtoken)
		gob.Register(new(Util.ResponseFromServer))
		returnValue := Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(returnValue)
		if errreturn != nil {
			log.Println("home auth encode direct from cache err", errreturn)
		}
		// Util.FailSafeCheckErr("home auth encode direct from cache err", errreturn)
		return
	}

	//First go through the Redis get cache
	user, ok, err := dao.GetCacheInfo(username)
	if err != nil {
		log.Println("redis get cache fail err", err)
	}
	// Util.FailSafeCheckErr("redis get cache fail:", err)
	//cache still valid
	//
	if ok && err == nil {
		log.Println("tcp home handle cache get info okay", user)
		gob.Register(new(Util.RealUser))
		tohttp := &Util.ResponseFromServer{Success: true, TcpData: user}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(tohttp)
		if errreturn != nil {
			panic(errreturn)
		}
		// Util.FailSafeCheckErr("no this ", errreturn)
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
			//do nothing
		}

		//save cache success
		//here how
		// if successCache {
		gob.Register(new(Util.RealUser))
		tohttp := &Util.ResponseFromServer{Success: true, TcpData: userdb}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(tohttp)
		if errreturn != nil {
			panic(errreturn)
		}
		// Util.FailSafeCheckErr("home handle encode err", errreturn)
		return
		// }

	}

	return
}
