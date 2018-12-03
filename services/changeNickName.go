package service

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"../DAO"
	"../Util"
)

func ChangeNickNameHandle(conn net.Conn, username string, nickname interface{}, token string) {
	//first check token from redis frist
	exists, errtoken := dao.CheckToken(username, token)

	// Util.FailSafeCheckErr("updatenickname checktoken cache err", errtoken)
	if !exists || errtoken != nil {
		//token expires or not correct
		gob.Register(new(Util.ResponseFromServer))
		tohttp := &Util.ResponseFromServer{Success: false, TcpData: nil}
		encoder := gob.NewEncoder(conn)
		errreturn := encoder.Encode(tohttp)
		if errreturn != nil {
			log.Println("nickname encode err", errreturn)
		}
		// Util.FailSafeCheckErr("changenickname encode err", errreturn)
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
	if errreturn != nil {
		log.Println("nickname encode err", errreturn)
	}
	// Util.FailSafeCheckErr("changenickname encode err", errreturn)
}
