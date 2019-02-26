package service

import (
	"fmt"
	"log"
	"net"

	dao "entry_task/DAO"
	data "entry_task/Data"

	"github.com/golang/protobuf/proto"
)

func UploadHandle(conn net.Conn, username string, avatar interface{}, token string) {
	exists, errtoken := dao.CheckToken(username, token)
	// data.FailSafeCheckErr("upload checktoken cache err", errtoken)

	//token not exists or not correct
	if !exists || errtoken != nil {
		log.Println("upload checktoken cache err", errtoken)
		// gob.Register(new(data.ResponseFromServer))
		returnValue := &data.ResponseFromServer{Success: proto.Bool(false), TcpData: nil}
		returnValueData, rErr := proto.Marshal(returnValue)
		if rErr != nil {
			panic(rErr)
		}
		_, wErr := conn.Write(returnValueData)
		if wErr != nil {
			panic(wErr)
		}

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
	}
	//mysql update not success
	tohttp := &data.ResponseFromServer{Success: proto.Bool(success), TcpData: nil}
	tohttpData, tErr := proto.Marshal(tohttp)
	if tErr != nil {
		panic(tErr)
	}
	_, wErr := conn.Write(tohttpData)
	if wErr != nil {
		panic(wErr)
	}
}
