package service

import (
	"fmt"
	"log"

	dao "entry_task/DAO"

	data "entry_task/Data"

	"github.com/golang/protobuf/proto"
)

// type LogoutService struct {
// }

// func LogoutHandle(conn net.Conn, toServerD *data.ToServerData, wg *sync.WaitGroup) {
func LogoutHandle(toServerD *data.ToServerData) (*data.ResponseFromServer, error) {
	// defer wg.Done()
	log.Println("----------tcp logout---------------------")
	tmpdata := &data.InfoWithUsername{}
	tmpErr := proto.Unmarshal(toServerD.GetHttpdata(), tmpdata)
	if tmpErr != nil {
		fmt.Println("logout err:", tmpErr)
		panic(tmpErr)
	}
	//invalid the cache in redis first
	err := dao.InvalidCache(tmpdata.GetUsername(), tmpdata.GetToken())
	if err != nil {
		log.Println("invalid logout err", err)
	}
	// data.FailSafeCheckErr("invalid logout usr", err)
	//return to logout handle
	success := (err == nil)
	// gob.Register(new(data.ResponseFromServer))

	returnValue := &data.ResponseFromServer{Success: proto.Bool(success), TcpData: nil}
	// returnValueData, rErr := proto.Marshal(returnValue)
	// if rErr != nil {
	// 	fmt.Println("logout marshal err:", rErr)
	// 	panic(rErr)
	// }
	// log.Println("------------tcp logouthandle returnValueData-------", returnValue)
	// _, writeErr := conn.Write(returnValueData)
	// log.Println("----------tcp logout write socket---------------------")
	// if writeErr != nil {
	// 	fmt.Println("logout write conn err,", writeErr)
	// }
	return returnValue, nil
}
