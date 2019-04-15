package service

import (
	"fmt"
	"log"

	dao "entry_task/DAO"
	data "entry_task/Data"

	"github.com/golang/protobuf/proto"
)

// func HomeHandle(conn net.Conn, username string, token interface{}) {
// func HomeHandle(conn net.Conn, toServerD *data.ToServerData, wg *sync.WaitGroup) {
func HomeHandle(toServerD *data.ToServerData) (*data.ResponseFromServer, error) {
	// defer wg.Done()
	tmpdata := &data.InfoWithUsername{}
	tmpErr := proto.Unmarshal(toServerD.GetHttpdata(), tmpdata)
	if tmpErr != nil {
		fmt.Println("login err:", tmpErr)
		panic(tmpErr)
	}
	//checktoken first
	exists, errtoken := dao.CheckToken(tmpdata.GetUsername(), tmpdata.GetToken())
	// data.FailSafeCheckErr("home checktoken cache err", errtoken)
	//1. cookie still exists but token expires
	//---solution: clear cookie first then redirect to login
	//2. cookie expires but token exists
	//---solution: redirect to login??? and refresh the token

	//token not exists or not correct
	if !exists || errtoken != nil {
		log.Println("home checktoken cache err", errtoken)
		// gob.Register(new(data.ResponseFromServer))

		returnValue := &data.ResponseFromServer{Success: proto.Bool(false), TcpData: nil}
		// returnValueData, rErr := proto.Marshal(returnValue)
		// if rErr != nil {
		// 	panic(rErr)
		// }
		// packHttp := Util.Pack(Util.PACK_CLIENT, returnValueData, false)
		// _, writeErr := conn.Write(packHttp)

		// // _, writeErr :=conn.Write(returnValueData)
		// if writeErr != nil {
		// 	panic(writeErr)
		// }
		return returnValue, nil
	}

	//First go through the Redis get cache
	user, ok, err := dao.GetCacheInfo(tmpdata.GetUsername())
	if err != nil {
		log.Println("redis get cache fail err", err)
	}
	// data.FailSafeCheckErr("redis get cache fail:", err)
	//cache still valid
	//
	if ok && err == nil {
		log.Println("tcp home handle cache get info okay", *user)
		// gob.Register(new(data.RealUser))
		userData, userErr := proto.Marshal(user)
		if userErr != nil {
			panic(userErr)
		}
		tohttp := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: userData}
		// tohttpData, toHttpErr := proto.Marshal(tohttp)
		// if toHttpErr != nil {
		// 	fmt.Println("tohttperr:", toHttpErr)
		// 	panic(toHttpErr)
		// }
		// packHttp := Util.Pack(Util.PACK_CLIENT, tohttpData, false)
		// _, writeErr := conn.Write(packHttp)
		// // _, writeErr := conn.Write(tohttpData)
		// if writeErr != nil {
		// 	fmt.Println("home tcp write err")
		// 	panic(writeErr)
		// }
		return tohttp, nil
	}
	log.Println("-----------not at redis,go to mysql---------------------")
	//cache expires or not exists then go to mysql
	userdb, okdb := dao.AllInfo(tmpdata.GetUsername())
	//retrieve from mysql success
	if okdb {
		//it will also save it to cache
		successCache := dao.SaveCacheInfo(tmpdata.GetUsername(), *userdb.Nickname, *userdb.Avatar)
		if !successCache {
			fmt.Println("update redis homne cache fail")
			//do nothing
		}

		//save cache success
		//here how
		// if successCache {

		userdbData, userdbErr := proto.Marshal(userdb)
		if userdbErr != nil {
			panic(userdbErr)
		}
		tohttp := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: userdbData}
		// tohttpData, toHttpErr := proto.Marshal(tohttp)
		// if toHttpErr != nil {
		// 	fmt.Println("tohttperr:", toHttpErr)
		// 	panic(toHttpErr)
		// }
		// packHttp := Util.Pack(Util.PACK_CLIENT, tohttpData, false)
		// _, writeErr := conn.Write(packHttp)
		// // _, writeErr := conn.Write(tohttpData)
		// if writeErr != nil {
		// 	fmt.Println("home write conn err,", writeErr)
		// }

		return tohttp, nil
		// }

	}

	return nil, nil
}
