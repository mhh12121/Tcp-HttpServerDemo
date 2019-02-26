package service

import (
	"fmt"
	"net"

	dao "entry_task/DAO"
	data "entry_task/Data"

	"github.com/golang/protobuf/proto"
)

func ChangeNickNameHandle(conn net.Conn, username string, nickname interface{}, token string) {
	//first check token from redis frist
	exists, errtoken := dao.CheckToken(username, token)
	if !exists || errtoken != nil {
		//token expires or not correct
		fmt.Println("tcp token expires")
		tohttp := &data.ResponseFromServer{Success: proto.Bool(false), TcpData: nil}
		tohttpData, tErr := proto.Marshal(tohttp)
		if tErr != nil {
			panic(tErr)
		}
		_, wErr := conn.Write(tohttpData)
		if wErr != nil {
			panic(wErr)
		}
		return
	}
	//update mysql first
	success, errorupdate := dao.UpdateNickname(username, nickname.(string))
	//then update cache
	if success && errorupdate == nil {
		//if successfully change data in mysql
		err := dao.UpdateCacheNickname(username, nickname.(string))
		//update cache fail
		if err != nil {
			fmt.Println("update nickname fail", err)
			//todo
			//do nothing
			// return
		}
	}
	// gob.Register(new(util.ResponseFromServer))
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
