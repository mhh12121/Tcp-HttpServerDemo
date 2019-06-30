package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"path"
	"path/filepath"
	"runtime"

	"entry_task/Conf"
	dao "entry_task/DAO"
	data "entry_task/Data"
	service "entry_task/services"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

// var FunctionCode map[int32]func(net.Conn, *data.ToServerData) = map[int32]func(net.Conn, *data.ToServerData){
// 	Util.LOGINCODE:  service.LoginHandle,
// 	Util.LOGOUTCODE: service.LogoutHandle,
// 	Util.HOMECODE:   service.HomeHandle,

// 	// 8003: service.LogoutHandle,
// }

type TCPServer struct {
	Proto string
	Addr  string
	// TcpMutex sync.Mutex
	// ServiceMap map[string]func(net.Conn, *data.ToServerData)
	// handler func(c *net.Conn)
}

func (tserver *TCPServer) Login(ctx context.Context, toServerD *data.ToServerData) (*data.ResponseFromServer, error) {

	return service.LoginHandle(toServerD)
}
func (tserver *TCPServer) Home(ctx context.Context, toServerD *data.ToServerData) (*data.ResponseFromServer, error) {
	return service.HomeHandle(toServerD)
}
func (tserver *TCPServer) Logout(ctx context.Context, toServerD *data.ToServerData) (*data.ResponseFromServer, error) {
	return service.LogoutHandle(toServerD)
}

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	_, fp, _, _ := runtime.Caller(0)
	p := path.Dir(fp)
	p = path.Dir(p)
	log.Println("log path", p)
	Conf.LoadConf(filepath.Join(p, "Conf/config.json"))
	// log.Println("dafas", Conf.Config)
}
func main() {

	dao.InitDB()
	tcpServer := &TCPServer{Proto: "tcp", Addr: Conf.Config.Connect.Tcpport}
	fmt.Println("tcp start ", Conf.Config.Connect.Tcphost)
	fmt.Println("tcp start ", Conf.Config.Connect.Tcpport)
	s := grpc.NewServer()
	data.RegisterAuthenticateServer(s, tcpServer)
	defer s.Stop()
	addr := Conf.Config.Connect.Tcphost + ":" + Conf.Config.Connect.Tcpport
	tcpAddr, addErr := net.ResolveTCPAddr("tcp4", addr)
	if addErr != nil {
		panic(addErr)
	}
	//----grpc already set loop listen------------------
	// for {
	log.Println("start listen loop tcp:")
	ln, err := net.ListenTCP("tcp", tcpAddr)
	// ln, err := net.Listen(tserver.Proto, ":"+tserver.Addr)

	if err != nil {
		fmt.Println("tcp listen failed:", err)
	}
	if errs := s.Serve(ln); errs != nil {
		log.Fatalf("grpc serve failed~~~~~~ %v", errs)
	}

}
