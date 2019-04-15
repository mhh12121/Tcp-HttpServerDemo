package main

import (
	"bufio"
	"encoding/json"
	"entry_task/Conf"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"

	data "entry_task/Data"
	Util "entry_task/Util"
	service "entry_task/services"

	"github.com/golang/protobuf/proto"
	pool "gopkg.in/fatih/pool.v2"
)

var connpool pool.Pool
var globalcon net.Conn

func init() {

}

var loginTemplate *template.Template
var homeTemplate *template.Template
var rpc *service.RpcHandle

func main() {
	_, filepath, _, _ := runtime.Caller(0)
	p := path.Dir(filepath)
	p = path.Dir(p)

	log.Println("log path", p)

	Conf.LoadConf(p + "/Conf/config.json")
	// var err error
	// tcpconn, err = net.Dial("tcp", Util.Tcpaddress+":"+Util.Tcpport)
	// tcpconn.SetReadDeadline(time.Now().Add(Util.TimeoutDuration))
	addr := Conf.Config.Connect.Tcphost + ":" + Conf.Config.Connect.Tcpport
	tcpAddr, addErr := net.ResolveTCPAddr("tcp4", addr)
	if addErr != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", addErr.Error())
		os.Exit(1)
	}
	factory := func() (net.Conn, error) {
		return net.DialTCP("tcp", nil, tcpAddr)
	}
	var err error
	connpool, err = pool.NewChannelPool(Conf.Config.Chanpool.Initsize, Conf.Config.Chanpool.Maxsize, factory)
	if err != nil {
		panic(err)
	}

	// con, err := net.("tcp", nil, tcpAddr)
	// if err != nil {
	// 	panic(err)
	// }
	// globalcon = con
	// http.HandleFunc("/", viewHandler)
	// loginTemplate = template.Must(template.ParseFiles("../view/login.html"))
	// homeTemplate = template.Must(template.ParseFiles("../view/Home.html"))

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(Util.UploadPath))))
	http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/login/auth", authHandler)
	http.HandleFunc("/Home", homeHandler)
	// http.HandleFunc("/Home/upload", uploadHandler)
	http.HandleFunc("/", testHandler)
	// http.HandleFunc("/Home/change", changeNickNameHandler)
	http.HandleFunc("/Home/logout", logoutHandler)
	http.HandleFunc("/test", testHandler)
	errhttp := http.ListenAndServe(Conf.Config.Connect.Httphost+":"+Conf.Config.Connect.Httpport, nil)
	log.Fatal(errhttp)

}
func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := template.Must(template.ParseFiles("../view/test.html"))
		t.Execute(w, nil)
	}

}

//generate simple token
func GenerateToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
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

//get response from server
func readServer(w http.ResponseWriter, r *http.Request, tcpconn net.Conn, ctype int32) (interface{}, bool) {
	buffer := make([]byte, 1024)
	dataServer := make([]byte, 0)
	// rw := bufio.NewReadWriter(bufio.NewReader(tcpconn), bufio.NewWriter(tcpconn))
	// c := bufio.NewReader(tcpconn)
	// defer tcpconn.Close()
	// <-readchan
	// var wg sync.WaitGroup
	fmt.Println("pass readserver")
	// wg.Add(1)
	switch ctype {
	case Util.LOGINCODE:
		{
			size, httperr := tcpconn.Read(buffer)
			fmt.Println("http util login read size---------------:", size)
			// dataServer := ReadByteLoop(tcpconn)
			dataServer = Util.Unpack(append(dataServer, buffer[:size]...))
			// for {

			// // size, httperr := globalcon.Read(dataServer)

			if httperr != nil {
				panic(httperr)
			}
			// fmt.Println("http get from server length", size)
			// dataServer, rErr := ioutil.ReadAll(tcpconn) //can it really read to the end???//#todo
			// fmt.Println("dataserver", dataServer)
			// if rErr != nil {
			// 	fmt.Fprintf(os.Stderr, "Fatal error: %s", rErr.Error())
			// 	os.Exit(1)

			// }
			// dataServer = bytes.Trim(dataServer, "\x00")
			dataResp := &data.ResponseFromServer{}
			dataErr := proto.Unmarshal(dataServer, dataResp) //illeagl tag 0
			if dataErr != nil {
				fmt.Println("proto login unmarshal", dataErr)
				panic(dataErr)
			}

			//something wrong in tcp server
			//or login fail
			fmt.Println("http read server:", dataResp)
			if dataResp.GetSuccess() {
				return nil, true
			}
			// else {
			// 	break
			// }

			// }
			return nil, false
		}
	case Util.HOMECODE:
		{
			size, httperr := tcpconn.Read(buffer)
			// dataServer := ReadByteLoop(tcpconn)
			dataServer = Util.Unpack(append(dataServer, buffer[:size]...))
			// _, httperr := tcpconn.Read(dataServer)
			// // _, httperr := globalcon.Read(dataServer)
			if httperr != nil {
				panic(httperr)
			}
			// // fmt.Println("http get from server length", size)
			// dataServer = bytes.Trim(dataServer, "\x00")
			// dataServer := ReadByteLoop(tcpconn)
			tmp := &data.ResponseFromServer{}
			tmpErr := proto.Unmarshal(dataServer, tmp)
			if tmpErr != nil {
				fmt.Println("http home unmarshal err", tmpErr)
				panic(tmpErr)
			}

			if tmp.GetSuccess() {
				//token expires or not correct
				if tmp.GetTcpData() == nil {
					return nil, false
				}
				tcpData := &data.RealUser{}
				tcpErr := proto.Unmarshal(tmp.GetTcpData(), tcpData)
				if tcpErr != nil {
					panic(tcpErr)
				}
				//token pass
				//but
				fmt.Println("redis cache not update")

				return tcpData, true
			}
			// }
			// return nil, false
		}
	case Util.LOGOUTCODE:
		{
			log.Println("--------------logout read from tcp-------------------")
			// _, cerr := tcpconn.Read(dataServer)
			// // _, cerr := globalcon.Read(dataServer)
			// // .Read(buff)
			// log.Println("--------------logout read data-------------------")
			// if cerr != nil {
			// 	fmt.Println("logout buferr", cerr)
			// 	panic(cerr)
			// }
			// dataServer = bytes.Trim(dataServer, "\x00")
			dataServer := ReadByteLoop(tcpconn)
			tmp := &data.ResponseFromServer{}
			tmpErr := proto.Unmarshal(dataServer, tmp)
			if tmpErr != nil {
				fmt.Println("http logout unmarshal err", tmpErr)
				panic(tmpErr)
			}
			if tmp.GetSuccess() {
				return tmp.GetTcpData(), tmp.GetSuccess()
			}

			return nil, false
		}
		// case "changeNickName":
		// 	{
		// 		for {
		// 			// bufio.NewReader(tcpconn)
		// 			size, cerr := tcpconn.Read(buff)
		// 			// .Read(buff)
		// 			if cerr != nil {
		// 				fmt.Println("nickname buferr", cerr)
		// 				panic(cerr)
		// 			}
		// 			if size == 0 {
		// 				fmt.Println("nickname nothing in conn")
		// 				return nil, false
		// 			}

		// 			tmp := &data.ResponseFromServer{}
		// 			tmpErr := proto.Unmarshal(buff[:size], tmp)
		// 			if tmpErr != nil {
		// 				fmt.Println("http nickname unmarshal err", tmpErr)
		// 				panic(tmpErr)
		// 			}

		// 			if tmp.GetSuccess() {

		// 				// if tmp.GetTcpData() == nil {
		// 				// 	break
		// 				// }

		// 				return nil, true
		// 			} else { //token expires or not correct
		// 				break
		// 			}
		// 		}
		// 		return nil, false
		// 	}
		// case "uploadAvatar":
		// 	{
		// 		for {
		// 			size, cerr := tcpconn.Read(buff)
		// 			if cerr != nil {
		// 				panic(cerr)
		// 			}
		// 			if size == 0 {
		// 				return nil, false
		// 			}

		// 			tmp := &data.ResponseFromServer{}
		// 			tmpErr := proto.Unmarshal(buff[:size], tmp)
		// 			if tmpErr != nil {
		// 				fmt.Println("avatar unmarshal")
		// 				panic(tmpErr)
		// 			}
		// 			if tmp.GetSuccess() {
		// 				return nil, true
		// 			} else {
		// 				break
		// 			}

		// 		}

		// 		return nil, false
		// 	}

	}
	return nil, false

}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//todo
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("----------------logout-------------------")
		tcpconn, errget := connpool.Get()
		defer tcpconn.Close()
		if errget != nil {
			panic(errget)
		}
		log.Println("logout pool length---------------", connpool.Len())
		// defer tcpconn.Close()
		user, erruser := r.Cookie("username")
		if erruser != nil {
			log.Println("logout no user cookie")
		}

		token, errtoken := r.Cookie("token")
		if errtoken != nil {
			log.Println("logout no token cookie")
		}
		httpWrap := &data.InfoWithUsername{Username: proto.String(user.Value), Info: []byte(""), Token: proto.String(token.Value)}
		httpData, hErr := proto.Marshal(httpWrap)
		if hErr != nil {
			fmt.Println("logout marshal", hErr)
			panic(hErr)
		}

		tmp := &data.ToServerData{Ctype: proto.Int32(Util.LOGOUTCODE), Httpdata: httpData}
		tmpdata, tErr := proto.Marshal(tmp)
		if tErr != nil {
			panic(tErr)
		}
		wrappedSend := Util.Pack(Util.PACK_CLIENT, tmpdata, false)
		// _, writeErr := globalcon.Write(wrappedSend)
		_, writeErr := tcpconn.Write(wrappedSend)
		if writeErr != nil {
			panic(writeErr)
		}
		// for {
		//go to tcp to invalid the cache

		_, successlogout := readServer(w, r, tcpconn, Util.LOGOUTCODE)
		// _, successlogout := readServer(w, r, globalcon, Util.LOGOUTCODE)
		log.Println("-----------logout successlogout ", successlogout)
		if successlogout { //to clear all cookie
			//temp struct
			logoutReturn := struct {
				Ok   bool
				Data interface{}
			}{
				true,
				"",
			}
			cookieuser := http.Cookie{
				Name:   "username",
				MaxAge: -1,
				Path:   "/",
			}
			cookietoken := http.Cookie{
				Name:   "token",
				MaxAge: -1,
				Path:   "/",
			}
			http.SetCookie(w, &cookieuser)
			http.SetCookie(w, &cookietoken)
			b, err := json.Marshal(logoutReturn)
			if err != nil {
				log.Println("logout http return to browser err")
				panic(err)
			}

			w.Header().Set("content-type", "application/json")
			w.Write(b)
			return
		}
		fmt.Println("logout:fail delete cache in redis")
		return
		// }
	}
}

//for login Get render
func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		log.Println("login enter")
		t := template.Must(template.ParseFiles("../view/login.html"))
		// usernamecookie, erruser := r.Cookie("username")
		// tokencookie, errtoken := r.Cookie("token")
		// //not found cookie
		// if erruser != nil || errtoken != nil {
		// 	fmt.Println("no cookie login", erruser)
		// 	t.Execute(w, nil)
		// 	return
		// }
		//There exists one problem:
		//1. Cookie exists but token in redis already expires
		//solution:also needs to authorize again?
		//if found username and token (no matter right or wrong)
		//tmpcommand
		// tcpconn, errget := connpool.Get()
		// if errget != nil {
		// 	panic(errget)
		// }
		// fmt.Println("http login get tcpconn", tcpconn)
		// defer tcpconn.Close()
		// // tcpconn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// fmt.Println("tcp conn:", tcpconn.RemoteAddr().String())
		// defer tcpconn.Close()
		// tempuser := &data.User{Username: proto.String(usernamecookie.Value), Password: proto.String(""), Token: proto.String(tokencookie.Value)}
		// tempuserData, tmpErr := proto.Marshal(tempuser)
		// if tmpErr != nil {
		// 	fmt.Println("login marshal tempuser err:", tmpErr)
		// 	panic(tmpErr)
		// }
		// tmpdata := &data.ToServerData{Ctype: proto.String("login"), Httpdata: tempuserData}
		// tmpdataSend, _ := proto.Marshal(tmpdata)
		// //----rpc call login------------

		// _, err := tcpconn.Write(tmpdataSend)
		// if err != nil {
		// 	fmt.Println("login write err(may):", err)
		// 	panic(err)
		// }

		// fmt.Println("encode usename pwd:", tmpdata)
		// // //loop to listen from server

		// // time.Sleep(2 * time.Second)
		// _, successlogin := readServer(w, r, tcpconn, "login") //tcpconn or rpc.Con
		// //success login

		// if successlogin {
		// 	fmt.Println("login success!!http")
		// 	//, MaxAge: Util.CookieExpires
		// 	log.Println("login cookie expr", Util.CookieExpires)

		// 	cookie := http.Cookie{Name: "username", Value: usernamecookie.Value, Path: "/", Expires: Util.CookieExpires}
		// 	http.SetCookie(w, &cookie)
		// 	cookie = http.Cookie{Name: "token", Value: tokencookie.Value, Path: "/", Expires: Util.CookieExpires}
		// 	http.SetCookie(w, &cookie)
		// 	http.Redirect(w, r, "/Home", http.StatusFound)

		// 	return
		// }

		t.Execute(w, nil)
		return
	}
	//login authentication
	if r.Method == http.MethodPost {

		tcpconn, errget := connpool.Get()
		defer tcpconn.Close()
		if errget != nil {
			panic(errget)
		}

		// tcpconn.SetReadDeadline(time.Now().Add(5 * time.Second))
		// fmt.Println("tcp conn and http conn", tcpconn.RemoteAddr().String(), tcpconn.LocalAddr().String())
		// fmt.Println("tcp conn and http conn", globalcon.RemoteAddr().String(), globalcon.LocalAddr().String())

		fmt.Println("enter!!!!!!")
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println("front username:", username)
		fmt.Println("front pwd:", password)

		//Wrap the data
		//this token here may be destroyed
		temptoken := GenerateToken(5)
		//test
		// temptoken := "test"

		tempuser := &data.User{Username: proto.String(username), Password: proto.String(password), Token: proto.String(temptoken)}
		tempuserData, tmpErr := proto.Marshal(tempuser)
		if tmpErr != nil {
			fmt.Println("login marshal tempuser err:", tmpErr)
			panic(tmpErr)
		}
		tmpdata := &data.ToServerData{Ctype: proto.Int32(Util.LOGINCODE), Httpdata: tempuserData}
		tmpdataSend, _ := proto.Marshal(tmpdata)
		wrappedSend := Util.Pack(Util.PACK_CLIENT, tmpdataSend, false)
		// fmt.Println("encode wrappedsend pwd:", wrappedSend)
		fmt.Println("encode wrappedsend length", len(wrappedSend))
		_, err := tcpconn.Write(wrappedSend)
		// _, err := globalcon.Write(wrappedSend)
		if err != nil {
			fmt.Println("login write err:", err)
			panic(err)
		}

		fmt.Println("encode usename pwd:", tmpdata)
		// //loop to listen from server

		// for {
		// time.Sleep(2 * time.Second)
		_, successlogin := readServer(w, r, tcpconn, Util.LOGINCODE)
		// _, successlogin := readServer(w, r, globalcon, Util.LOGINCODE)
		//success login
		// tcpconn.Close()
		if successlogin {
			fmt.Println("login success!!http")
			//, MaxAge: Util.CookieExpires
			log.Println("login cookie expr", Util.CookieExpires)

			//------------------------comment for test------------------------------
			cookie := http.Cookie{Name: "username", Value: username, Path: "/", Expires: Util.CookieExpires}
			http.SetCookie(w, &cookie)
			cookie = http.Cookie{Name: "token", Value: temptoken, Path: "/", Expires: Util.CookieExpires}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/Home", http.StatusFound)
			//------------------------comment for test above------------------------------
			return
		}
		//wrong password
		http.Redirect(w, r, "/login", http.StatusFound)

		// w.WriteHeader(http.StatusForbidden)
		// w.Write([]byte(Util.ResWrongStr))
		return

		// }

	}

}

//after login //todo
func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		t := template.Must(template.ParseFiles("../view/Home.html"))
		cookieuser, erruser := r.Cookie("username")
		cookietoken, errtoken := r.Cookie("token")
		if erruser != nil {
			fmt.Println("cookie home user", erruser)
			http.Redirect(w, r, "/login", http.StatusFound)
			t.Execute(w, nil)
			return
		}
		if errtoken != nil {
			fmt.Println("cookie home token", errtoken)
			http.Redirect(w, r, "/login", http.StatusFound)
			t.Execute(w, nil)
			return
		}

		tcpconn, errget := connpool.Get()
		defer tcpconn.Close()
		if errget != nil {
			panic(errget)
		}
		log.Println("home rendering")

		//send to tcp server
		tokenwithusername := &data.InfoWithUsername{Username: proto.String(cookieuser.Value), Token: proto.String(cookietoken.Value)}
		tokenwithusernameData, terr := proto.Marshal(tokenwithusername)
		if terr != nil {
			panic(terr)
		}
		tmp := &data.ToServerData{Ctype: proto.Int32(Util.HOMECODE), Httpdata: tokenwithusernameData}
		tmpData, tmpErr := proto.Marshal(tmp)
		if tmpErr != nil {
			panic(tmpErr)
		}
		fmt.Println("http cookie ", cookieuser.Value)
		//----------------wrap handle code-------------
		// wrappedSend := make([]byte, len(tmpData))
		// binary.BigEndian.PutUint32(wrappedSend, HOMECODE)
		// //----------------wrap compress code-------------
		// compressMark := make([]byte, 1)
		// binary.BigEndian.PutUint16(compressMark, NOCOMPRESS)
		wrappedSend := Util.Pack(Util.PACK_CLIENT, tmpData, false)
		_, werr := tcpconn.Write(wrappedSend)
		// _, werr := globalcon.Write(wrappedSend)
		if werr != nil {
			panic(werr)
		}
		log.Println("home render loop", tmpData)
		datar, successHome := readServer(w, r, tcpconn, Util.HOMECODE)
		// datar, successHome := readServer(w, r, globalcon, Util.HOMECODE)

		//token correct
		if successHome {
			ruser := datar.(*data.RealUser)
			t.Execute(w, &ruser)
			return
		}
		//token not correct
		//clear cookie and then redirect
		log.Println("token expires home page")
		setcookieuser := http.Cookie{
			Name:   "username",
			MaxAge: -1,
			Path:   "/",
		}
		setcookietoken := http.Cookie{
			Name:   "token",
			MaxAge: -1,
			Path:   "/",
		}
		http.SetCookie(w, &setcookieuser)
		http.SetCookie(w, &setcookietoken)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
		// }
		// }

	}

}

//upload avatar handler
// func uploadHandler(w http.ResponseWriter, r *http.Request) {
// 	//simple one
// 	// var errget error
// 	if r.Method == http.MethodPost {
// 		cookieuser, erruser := r.Cookie("username")
// 		if erruser != nil {
// 			log.Println("http upload file user err:", erruser)
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}
// 		cookietoken, errtoken := r.Cookie("token")
// 		if errtoken != nil {
// 			log.Println("http upload file token err:", errtoken)
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}

// 		tcpconn, errget := connpool.Get()
// 		if errget != nil {
// 			panic(errget)
// 		}
// 		defer tcpconn.Close()
// 		// Util.FailFastCheckErr(errget)
// 		// defer tcpconn.Close()
// 		file, handler, err := r.FormFile("profile")
// 		defer file.Close()
// 		if err != nil {
// 			fmt.Println("http upload get file err", err)
// 			http.Redirect(w, r, "/Home", http.StatusFound)
// 			return
// 		}
// 		//check if file format is correct
// 		//todo
// 		//some kinds of file may cause page crash
// 		filename, isLegal := checkAndCreateFileName(handler.Filename)
// 		if !isLegal {
// 			log.Println("illegal file format")
// 			http.Redirect(w, r, "/Home", http.StatusFound)
// 			return
// 		}

// 		// fmt.Fprintf(w, "%v", handler.Header)
// 		f, err := os.OpenFile(filepath.Join(Util.UploadPath, filename), os.O_WRONLY|os.O_CREATE, 0666)
// 		if err != nil {
// 			fmt.Println("http openfile fail", err)
// 			return
// 		}
// 		defer f.Close()
// 		io.Copy(f, file)

// 		//get username from cookie

// 		tempAvatar := &data.InfoWithUsername{Username: proto.String(cookieuser.Value), Info: []byte(filename), Token: proto.String(cookietoken.Value)}
// 		tempAvatarData, tErr := proto.Marshal(tempAvatar)
// 		if tErr != nil {
// 			panic(tErr)
// 		}
// 		uploadToServer := &data.ToServerData{Ctype: proto.String("uploadAvatar"), Httpdata: tempAvatarData}
// 		uploadServerData, uErr := proto.Marshal(uploadToServer)
// 		if uErr != nil {
// 			panic(uErr)
// 		}
// 		_, wErr := tcpconn.Write(uploadServerData)
// 		if wErr != nil {
// 			panic(wErr)
// 		}
// 		//listen response from tcp server
// 		// for {
// 		// readchan := make(chan int)
// 		_, successupload := readServer(w, r, tcpconn, "uploadAvatar")
// 		if successupload {
// 			http.Redirect(w, r, "/Home", http.StatusFound)
// 			return
// 		}
// 		// w.WriteHeader(http.StatusUnauthorized)
// 		//if db crash or token wrong
// 		http.Redirect(w, r, "/Home", http.StatusFound)
// 		return
// 		// }
// 	}

// }
// func changeNickNameHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		tcpconn, errget := connpool.Get()
// 		if errget != nil {
// 			panic(errget)
// 		}
// 		// Util.FailFastCheckErr(errget)
// 		defer tcpconn.Close()
// 		newnickname := r.FormValue("newnickname")
// 		log.Println("homenickname", newnickname)
// 		cookieuser, erruser := r.Cookie("username")
// 		cookietoken, errtoken := r.Cookie("token")
// 		if erruser != nil {
// 			//cookie not exists or be destroyed
// 			fmt.Println("change nickname get cookie fail", erruser)
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}
// 		if errtoken != nil {
// 			fmt.Println("change nickname get cookie fail", erruser)
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			return
// 		}

// 		tempMap := &data.InfoWithUsername{Username: proto.String(cookieuser.Value), Info: []byte(newnickname), Token: proto.String(cookietoken.Value)}
// 		tempMapData, tmpErr := proto.Marshal(tempMap)
// 		if tmpErr != nil {
// 			panic(tmpErr)
// 		}
// 		uploadToServer := &data.ToServerData{Ctype: proto.String("changeNickName"), Httpdata: tempMapData}
// 		uploadToServerData, uErr := proto.Marshal(uploadToServer)
// 		if uErr != nil {
// 			panic(uErr)
// 		}
// 		_, wErr := tcpconn.Write(uploadToServerData)
// 		if wErr != nil {
// 			fmt.Println("nickname write marshal")
// 			panic(wErr)
// 		}

// 		_, success := readServer(w, r, tcpconn, "changeNickName")
// 		if success {
// 			http.Redirect(w, r, "/Home", http.StatusFound)
// 			return
// 		} else { //token expires or cookie expires
// 			// w.WriteHeader(http.StatusUnauthorized)
// 			fmt.Println("token wrong???")
// 		}

// 	}
// }

func checkAndCreateFileName(oldName string) (newName string, isLegal bool) {
	ext := path.Ext(oldName)
	// uppercase
	//todo
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		newName = Util.GetFileName(oldName, ext)
		isLegal = true
	}
	return newName, isLegal
}
