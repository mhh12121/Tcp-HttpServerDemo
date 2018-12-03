package service

import (
	"encoding/gob"
	"log"
	"net"

	"../DAO"
	"../Util"
)

// type LogoutService struct {
// }

func LogoutHandle(conn net.Conn, username string, token interface{}) {
	//invalid the cache in redis first
	err := dao.InvalidCache(username, token.(string))
	if err != nil {
		log.Println("invalid logout err", err)
	}
	// Util.FailSafeCheckErr("invalid logout usr", err)
	//return to logout handle
	success := (err == nil)
	gob.Register(new(Util.ResponseFromServer))
	returnValue := Util.ResponseFromServer{Success: success, TcpData: nil}
	encoder := gob.NewEncoder(conn)
	errreturn := encoder.Encode(returnValue)
	if errreturn != nil {
		log.Println("logout encode err", errreturn)
	}
}
