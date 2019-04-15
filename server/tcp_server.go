package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"entry_task/Conf"
	dao "entry_task/DAO"
	data "entry_task/Data"
	Util "entry_task/Util"
	service "entry_task/services"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
)

var FunctionCode map[int32]func(net.Conn, *data.ToServerData) = map[int32]func(net.Conn, *data.ToServerData){
	Util.LOGINCODE:  service.LoginHandle,
	Util.LOGOUTCODE: service.LogoutHandle,
	Util.HOMECODE:   service.HomeHandle,

	// 8003: service.LogoutHandle,
}

type TCPServer struct {
	Proto    string
	Addr     string
	TcpMutex sync.Mutex
	// ServiceMap map[string]func(net.Conn, *data.ToServerData)
	// handler func(c *net.Conn)
}

func (tserver *TCPServer) Run() {
	addr := Conf.Config.Connect.Tcphost + ":" + Conf.Config.Connect.Tcpport
	tcpAddr, addErr := net.ResolveTCPAddr("tcp4", addr)
	if addErr != nil {
		panic(addErr)
	}
	ln, err := net.ListenTCP("tcp", tcpAddr)
	// ln, err := net.Listen(tserver.Proto, ":"+tserver.Addr)

	if err != nil {
		fmt.Println("tcp listen failed:", err)
	}

	for {
		fmt.Println("block accept")
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("tcp server connection failed:", err)
			continue
		}
		log.Println("tcp listen loop", conn)
		// go handleAll(conn)
		go tserver.handleAll(conn)

	}
}

//some bugs here???//todo
func ReadByteLoop(conn net.Conn) []byte {
	connbuf := bufio.NewReader(conn)
	b, _ := connbuf.ReadByte()
	var msgData []byte
	if connbuf.Buffered() > 0 {

		msgData = append(msgData, b)
		for connbuf.Buffered() > 0 {
			b, err := connbuf.ReadByte()
			if err != nil {
				fmt.Print("unreadable readbyte or null ", err)
				break

			}
			msgData = append(msgData, b)
		}

	}
	return msgData
}

/*
Use length description:
-----------------------+---------------------------------+------------------------+-----------+-------------------
header(8bytes)        +length(4 bytes)=Data1's length   +  IF COMPRESSED(1 byte) +   Data1  +   length(2)//todo
-----------------------+---------------------------------+------------------------+------------+-------------------
*/
//not used temp
// func reader(conn net.Conn, readerChannel chan []byte) {

// 	log.Println("readerchannel begin length-----------------------", len(readerChannel))
// 	// realData,ok := <-readerChannel
// 	// if ok{
// 	// 	log.Println("------------already pass value---------------------------------")
// 	// }
// 	// log.Println("readerchannel now----------", readerChannel)
// 	select {
// 	case realData := <-readerChannel:
// 		{
// 			log.Println("readerchannel length after retrieve-----------------------", len(readerChannel))
// 			log.Println("----------tcp channel reader block here------------")
// 			toServerD := &data.ToServerData{}
// 			dataErr := proto.Unmarshal(realData, toServerD)
// 			fmt.Println("readchannel", toServerD)

// 			if dataErr != nil {
// 				fmt.Println("proto", dataErr)
// 				panic(dataErr)
// 			}
// 			if toServerD.GetCtype() == Util.LOGOUTCODE {
// 				fmt.Println("----------logout tcp get reader------------")
// 			}

// 			FunctionCode[toServerD.GetCtype()](conn, toServerD)

// 		}
// 		// case <-time.After(time.Second * 10):
// 		// 	fmt.Println("10 seconds after reader channel")
// 	}

// }
func (tserver *TCPServer) handleAll(conn net.Conn) {
	fmt.Println("handleall coming")
	tmpBuffer := make([]byte, 0) //save splitted data(not included RealData but header,zipField)
	// readerChannel := make(chan []byte, 1024) //realdata
	// tserver.ReaderChannel = make(chan []byte, 1024)
	// go reader(conn, readerChannel)
	// defer close(readerChannel)
	buffer := make([]byte, 1024)
	// for { //continue read buffer from socket
	// var wg sync.WaitGroup
	// wg.Add(2)
	size, err := conn.Read(buffer) //all data at once,but maybe the data is not complete
	if err != nil {
		if err == io.EOF {
			fmt.Println("eof", err)
			// break
		}
		fmt.Println("error why handleall", err)
		return
	}
	// tmpBuffer = Util.Unpack(append(tmpBuffer, buffer[:size]...), readerChannel)
	// tserver.TcpMutex.Lock()

	tmpBuffer = Util.Unpack(append(tmpBuffer, buffer[:size]...)) //, readerChannel
	// tserver.TcpMutex.Unlock()

	// log.Println("readerchannel length after retrieve-----------------------", len(readerChannel))
	//--------------------------------------not used channel below------------------------------
	log.Println("----------tcp channel reader block here------------")
	toServerD := &data.ToServerData{}
	dataErr := proto.Unmarshal(tmpBuffer, toServerD)
	fmt.Println("readchannel", toServerD)
	// dataErr := proto.Unmarshal(buff[:int(size)], toServerD)
	// dataErr := proto.Unmarshal(readData, toServerD)
	if dataErr != nil {
		fmt.Println("proto", dataErr)
		panic(dataErr)
	}
	if toServerD.GetCtype() == Util.LOGOUTCODE {
		fmt.Println("----------logout tcp get reader------------")
	}

	FunctionCode[toServerD.GetCtype()](conn, toServerD)
	// wg.Wait()
	// tserver.TcpMutex.Unlock()

	// }
	fmt.Println("handleall ___ end0--------------")

	// Msglength := make([]byte, 4)

	// _, cerr := con.Read(Msglength)
	// if cerr != nil {
	// 	if cerr == io.EOF {
	// 		fmt.Println("eof read ")

	// 	}
	// 	fmt.Println("buferr", cerr)
	// 	panic(cerr)
	// 	// break
	// }

	// realSize := binary.BigEndian.Uint64(Msglength)
	// fmt.Println("data size:", realSize)
	// //-------------check compress or not----------
	// checkCompress := make([]byte, 1)
	// _, cerrCompress := con.Read(checkCompress)
	// if cerrCompress != nil {
	// 	if cerrCompress == io.EOF {
	// 		fmt.Println("eof read ")

	// 	}
	// 	fmt.Println("buferr", cerrCompress)
	// 	panic(cerrCompress)
	// 	// break
	// }
	// x := binary.BigEndian.Uint16(checkCompress)
	// if Depress(x) {

	// }
	// //---------------decode real data-------------------
	// realData := make([]byte, realSize)
	// _, cerr2 := con.Read(realData)
	// if cerr2 != nil {
	// 	if cerr2 == io.EOF {
	// 		fmt.Println("eof read ")

	// 	}
	// 	fmt.Println("buferr", cerr2)
	// 	panic(cerr2)
	// 	// break
	// }

	// toServerD := &data.ToServerData{}
	// dataErr := proto.Unmarshal(realData, toServerD)
	// // dataErr := proto.Unmarshal(buff[:int(size)], toServerD)
	// // dataErr := proto.Unmarshal(readData, toServerD)
	// if dataErr != nil {
	// 	fmt.Println("proto", dataErr)
	// 	panic(dataErr)
	// }
	// FunctionCode[realSize](con, toServerD)
	//according to Ctype to set the response
	// switch *toServerD.Ctype {
	// case "login":
	// 	fmt.Println("login enter before")
	// 	tserver.ServiceMap["login"](tserver.Con, toServerD)
	// 	// service.LoginHandle(conn, *tmpdata)
	// //Home tcp
	// case "home":
	// 	tserver.ServiceMap["home"](tserver.Con, toServerD)

	// log.Println("home tcp decode data", tmpdata)
	// service.HomeHandle(conn, tmpdata.GetUsername(), tmpdata.GetToken())

	// case "uploadAvatar":
	// 	tmpdata := &data.InfoWithUsername{}
	// 	tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
	// 	if tmpErr != nil {
	// 		panic(tmpErr)
	// 	}

	// 	fmt.Println("tcp upload file decode data", tmpdata)
	// 	service.UploadHandle(conn, tmpdata.GetUsername(), tmpdata.GetInfo(), tmpdata.GetToken())
	// 	// uploadHandle(conn, , tmpdata.Info, tmpdata.Token)

	// case "changeNickName":
	// 	tmpdata := &data.InfoWithUsername{}
	// 	tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
	// 	if tmpErr != nil {
	// 		panic(tmpErr)
	// 	}
	// 	fmt.Println("tcp change nickname decode data ", tmpdata)
	// 	service.ChangeNickNameHandle(conn, tmpdata.GetUsername(), string(tmpdata.GetInfo()[:]), tmpdata.GetToken())

	// case "logout":
	// 	tmpdata := &data.InfoWithUsername{}
	// 	tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
	// 	if tmpErr != nil {
	// 		fmt.Println("logout err:", tmpErr)
	// 		panic(tmpErr)
	// 	}
	// 	fmt.Println("tcp change logout decode data ", tmpdata)
	// 	service.LogoutHandle(conn, tmpdata.GetUsername(), tmpdata.GetToken())
	// }
	// //flush in case of appendding
	// buff = make([]byte, 2048)

	// }
}

func init() {
	_, fp, _, _ := runtime.Caller(0)
	p := path.Dir(fp)
	p = path.Dir(p)
	log.Println("log path", p)
	Conf.LoadConf(filepath.Join(p, "Conf/config.json"))
	// log.Println("dafas", Conf.Config)
}
func main() {
	dao.InitDB()
	// dao.RedisInit()
	fmt.Println("tcp start ", Conf.Config.Connect.Tcphost)
	fmt.Println("tcp start ", Conf.Config.Connect.Tcpport)
	tcpServer := &TCPServer{Proto: "tcp", Addr: Conf.Config.Connect.Tcpport}
	tcpServer.Run()

}

// func handleAll(conn net.Conn) {
// 	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))
// 	buff := make([]byte, 2048)
// 	// c := bufio.NewReader(conn)
// 	defer conn.Close()
// 	//this for loop is for one connection with multiple requests!!!
// 	for {
// 		size, cerr := conn.Read(buff)
// 		if cerr != nil {
// 			fmt.Println("buferr", cerr)
// 			panic(cerr)
// 			// break
// 		}
// 		if size == 0 {
// 			fmt.Println("no message")
// 			break
// 		}
// 		// _, ioerr := io.ReadFull(c, buff[:int(size)])
// 		// if ioerr != nil {
// 		// 	fmt.Println(ioerr)
// 		// 	panic(ioerr)
// 		// }

// 		toServerD := &data.ToServerData{}
// 		dataErr := proto.Unmarshal(buff[:int(size)], toServerD)
// 		if dataErr != nil {
// 			fmt.Println("proto", dataErr)
// 			panic(dataErr)
// 		}

// 		//according to Ctype to set the response
// 		switch *toServerD.Ctype {
// 		case "login":
// 			tmpdata := &data.User{}
// 			tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
// 			if tmpErr != nil {
// 				fmt.Println("login err:", tmpErr)
// 				panic(tmpErr)
// 			}
// 			// tmpdata := data.Httpdata
// 			service.LoginHandle(conn, *tmpdata)
// 		//Home tcp
// 		case "home":
// 			tmpdata := &data.InfoWithUsername{}
// 			tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
// 			if tmpErr != nil {
// 				fmt.Println("login err:", tmpErr)
// 				panic(tmpErr)
// 			}

// 			log.Println("home tcp decode data", tmpdata)
// 			service.HomeHandle(conn, tmpdata.GetUsername(), tmpdata.GetToken())

// 		case "uploadAvatar":
// 			tmpdata := &data.InfoWithUsername{}
// 			tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
// 			if tmpErr != nil {
// 				panic(tmpErr)
// 			}

// 			fmt.Println("tcp upload file decode data", tmpdata)
// 			service.UploadHandle(conn, tmpdata.GetUsername(), tmpdata.GetInfo(), tmpdata.GetToken())
// 			// uploadHandle(conn, , tmpdata.Info, tmpdata.Token)

// 		case "changeNickName":
// 			tmpdata := &data.InfoWithUsername{}
// 			tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
// 			if tmpErr != nil {
// 				panic(tmpErr)
// 			}
// 			fmt.Println("tcp change nickname decode data ", tmpdata)
// 			service.ChangeNickNameHandle(conn, tmpdata.GetUsername(), string(tmpdata.GetInfo()[:]), tmpdata.GetToken())

// 		case "logout":
//
// 			fmt.Println("tcp change logout decode data ", tmpdata)
// 			service.LogoutHandle(conn, tmpdata.GetUsername(), tmpdata.GetToken())
// 		}
// 		//flush in case of appendding
// 		buff = make([]byte, 2048)

// 	}
// }

// func loginHandle(conn net.Conn, ruser data.User) {
// 	service.LoginHandle(conn, ruser)
// }

// func logoutHandle(conn net.Conn, username string, token interface{}) {
// 	service.LogoutHandle(conn, username, token)
//

// }
// func homeHandle(conn net.Conn, username string, token interface{}) {
// 	service.HomeHandle(conn, username, token)

// }

// //
// func uploadHandle(conn net.Conn, username string, avatar interface{}, token string) {
// 	service.UploadHandle(conn, username, avatar, token)
// 	// Util.FailSafeCheckErr("uploadfile encode err", errreturn)
// 	// return success
// }

// //
// func changeNickNameHandle(conn net.Conn, username string, nickname interface{}, token string) {
// 	service.ChangeNickNameHandle(conn, username, nickname, token)

// }
