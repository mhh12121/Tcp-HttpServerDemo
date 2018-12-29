package service

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	dao "entry_task/DAO"
	"entry_task/Util"
)

func UploadHandle(conn net.Conn, username string, avatar interface{}, token string) {
	exists, errtoken := dao.CheckToken(username, token)
	// Util.FailSafeCheckErr("upload checktoken cache err", errtoken)

	//token not exists or not correct
	if !exists || errtoken != nil {
		log.Println("upload checktoken cache err", errtoken)
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
	if errreturn != nil {
		log.Println("nickname encode err", errreturn)
	}
}
