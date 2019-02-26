package service

import (
	"fmt"
	"log"
	"net"

	dao "entry_task/DAO"
	data "entry_task/Data"
	Util "entry_task/Util"

	"github.com/golang/protobuf/proto"
)

//tcp handle
func LoginHandle(conn net.Conn, ruser data.User) {
	// log.Println("login tcp decode data", tmpdata)
	//get remote Addr
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("tcp server connect:" + remoteAddr)
	//first go through redis cache
	//check if exists or different
	//what if login in another device?
	exists, errtoken := dao.CheckToken(*ruser.Username, *ruser.Token)
	if errtoken != nil {
		log.Println("login checktoken cache err", errtoken)
	}

	//todo
	//some problems here(consistency)
	//1.checktoken in redis success then return success msg to http
	//2.http redirect to home
	//3.in the same time, the token in redis expires
	if exists {
		//if exists just take info from cache
		// gob.Register(new(data.ResponseFromServer))

		returnValue := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: nil}
		returnValueData, errReturn := proto.Marshal(returnValue)
		if errReturn != nil {
			fmt.Println("proto login marshal:", errReturn)
			panic(errReturn)
		}
		// writer := bufio.NewWriter(conn)

		_, writeErr := conn.Write(returnValueData)
		if writeErr != nil {
			fmt.Println("write login:", writeErr)
			panic(writeErr)
		}
		//-------------old ---------------------
		// encoder := gob.NewEncoder(conn)
		// errreturn := encoder.Encode(returnValue)
		// if errreturn != nil {
		// 	log.Println("login auth encode direct from cache err", errreturn)
		// }
		//-----------------------------------

		return
	}

	//check from mysql
	success, errorcheck := dao.Check(*ruser.Username, *ruser.Password)

	//login fail
	if !success || errorcheck != nil {
		log.Println("password wrong! any error?:", errorcheck)

		returnValue := &data.ResponseFromServer{Success: proto.Bool(false), TcpData: nil}
		returnValueData, errReturn := proto.Marshal(returnValue)
		if errReturn != nil {
			panic(errReturn)
		}
		_, writeErr := conn.Write(returnValueData)
		if writeErr != nil {
			panic(writeErr)
		}
		//-------------old ---------------------
		// encoder := gob.NewEncoder(conn)
		// errreturn := encoder.Encode(returnValue)
		// if errreturn != nil {
		// 	log.Println("login mysql encode err", errreturn)
		// }
		// -------------------------------------
		return
	}

	//if mysql check success, it will save it to redis as cache or update cache
	tokenerr := dao.SetToken(*ruser.Username, *ruser.Token, Util.TokenExpires)
	if tokenerr != nil {
		log.Println("login save cache err", tokenerr)
	}
	// data.FailSafeCheckErr("login save cache err", tokenerr)
	//login success
	log.Println("login handle tcp")

	returnValue := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: nil}
	returnValueData, errReturn := proto.Marshal(returnValue)
	if errReturn != nil {
		fmt.Println("errReturn:", errReturn)
		panic(errReturn)
	}

	log.Println("login handle tcp write next")
	_, writeErr := conn.Write(returnValueData)
	if writeErr != nil {
		fmt.Println("login writeErr", writeErr)
		panic(writeErr)
	}

	return
}
