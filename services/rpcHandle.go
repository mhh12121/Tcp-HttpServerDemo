package service

import (
	"encoding/gob"
	"net"

	"entry_task/Util"
)

type RpcHandle struct {
	Con   net.Conn
	Ctype string
	Data  interface{}
	// RegMap map[string]*Processor
}

// func (rpc *RpcHandle) Register(ctype string) error {
// 	rpc.RegMap["login"] = rpcLogin(rpc)
// }

//read from server
func (rpc *RpcHandle) Get() error {

	switch rpc.Ctype {
	case "login":
		{

		}

	}
	return nil
}

// func (rpc *RpcHandle)  error {
// 	gob.Register(new(Util.User))
// 	gob.Register(new(Util.RealUser))
// 	gob.Register(new(Util.ToServerData))
// 	encoder := gob.NewEncoder(rpc.Con)
// 	tmpdata := rpc.Data.(*Util.ToServerData)
// 	err := encoder.Encode(tmpdata)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// readchan <- 1
// 	return err
// }

//send to server
func (rpc *RpcHandle) Set(readchan chan int) error {
	switch rpc.Ctype {
	case "login":
		{
			gob.Register(new(Util.User))
			gob.Register(new(Util.RealUser))
			gob.Register(new(Util.ToServerData))
			encoder := gob.NewEncoder(rpc.Con)
			tmpdata := rpc.Data.(*Util.ToServerData)
			err := encoder.Encode(tmpdata)
			if err != nil {
				panic(err)
			}
			// readchan <- 1
			return err
		}
	case "logout":
		{

		}
	case "updateNickName":
		{

		}
	case "":
		{
		}
	}
	return nil
}
