package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"

	"entry_task/Conf"

	data "entry_task/Data"
	"entry_task/Util"
	service "entry_task/services"

	"github.com/golang/protobuf/proto"
	pool "gopkg.in/fatih/pool.v2"
)

var connpool pool.Pool

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
	// Util.FailFastCheckErr(err)
	// tcpconn.SetReadDeadline(time.Now().Add(Util.TimeoutDuration))
	factory := func() (net.Conn, error) {
		return net.Dial("tcp", Conf.Config.Connect.Tcphost+":"+Conf.Config.Connect.Tcpport)
	}
	var err error
	connpool, err = pool.NewChannelPool(Conf.Config.Chanpool.Initsize, Conf.Config.Chanpool.Maxsize, factory)
	if err != nil {
		panic(err)
	}
	// now you can get a connection from the pool, if there is no connection
	// available it will create a new one via the factory function.

	// http.HandleFunc("/", viewHandler)
	loginTemplate = template.Must(template.ParseFiles("../view/login.html"))
	homeTemplate = template.Must(template.ParseFiles("../view/Home.html"))

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(Util.UploadPath))))
	http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/login/auth", authHandler)
	// http.HandleFunc("/Home", homeHandler)
	// http.HandleFunc("/Home/upload", uploadHandler)
	http.HandleFunc("/", testHandler)
	// http.HandleFunc("/Home/change", changeNickNameHandler)
	// http.HandleFunc("/Home/logout", logoutHandler)
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

// func logoutHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		//todo
// 		//Actually not use it
// 		tmp := r.Context()
// 		fmt.Println("log out ", tmp)
// 		//
// 		tcpconn, errget := connpool.Get()
// 		if errget != nil {
// 			panic(errget)
// 		}

// 		user, erruser := r.Cookie("username")
// 		if erruser != nil {
// 			log.Println("logout no user cookie")
// 		}

// 		token, errtoken := r.Cookie("token")
// 		if errtoken != nil {
// 			log.Println("logout no token cookie")
// 		}
// 		httpdata := &Util.InfoWithUsername{Username: user.Value, Info: token.Value}
// 		tmpdata := &Util.ToServerData{Ctype: "logout", HttpData: httpdata}

// 		gob.Register(new(Util.InfoWithUsername))
// 		gob.Register(new(Util.ToServerData))
// 		encoder := gob.NewEncoder(tcpconn)
// 		err := encoder.Encode(tmpdata)
// 		if err != nil {
// 			log.Println("logout http encode err")
// 			panic(err)
// 		}
// 		// Util.FailSafeCheckErr("logout http encode err", err)

// 		// for {
// 		//go to tcp to invalid the cache
// 		readchan := make(chan int)
// 		_, successlogout := readServer(w, r, tcpconn, "logout", readchan)
// 		if successlogout { //to clear all cookie
// 			//temp struct
// 			logoutReturn := struct {
// 				Ok   bool
// 				Data interface{}
// 			}{
// 				true,
// 				"",
// 			}
// 			cookieuser := http.Cookie{
// 				Name:   "username",
// 				MaxAge: -1,
// 				Path:   "/",
// 			}
// 			cookietoken := http.Cookie{
// 				Name:   "token",
// 				MaxAge: -1,
// 				Path:   "/",
// 			}
// 			http.SetCookie(w, &cookieuser)
// 			http.SetCookie(w, &cookietoken)
// 			b, err := json.Marshal(logoutReturn)
// 			if err != nil {
// 				log.Println("logout http return to browser err")
// 				panic(err)
// 			}

// 			w.Header().Set("content-type", "application/json")
// 			w.Write(b)
// 			return
// 		}
// 		fmt.Println("logout:fail delete cache in redis")
// 		return
// 		// }

// 	}

// }

//generate simple token
func GenerateToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//get response from server
func readServer(w http.ResponseWriter, r *http.Request, tcpconn net.Conn, ctype string) (interface{}, bool) {
	// buff := make([]byte, 1024)
	// c := bufio.NewReader(tcpconn)
	// defer tcpconn.Close()
	// <-readchan
	// fmt.Println("pass readchan")
	switch ctype {
	case "login":
		{
			size, cerr := ioutil.ReadAll(tcpconn)
			// .Read(buff)
			if cerr != nil {
				fmt.Println("buferr", cerr)
				panic(cerr)
			}
			// _, ioerr := io.ReadFull(c, buff[:int(size)])
			// if ioerr != nil {
			// 	fmt.Println(ioerr)
			// 	panic(ioerr)
			// }
			data := &data.ResponseFromServer{}
			dataErr := proto.Unmarshal(size, data)
			if dataErr != nil {
				fmt.Println("proto login unmarshal", dataErr)
				panic(dataErr)
			}
			//decoder
			// gob.Register(new(Util.ResponseFromServer))
			// decoder := gob.NewDecoder(tcpconn)
			// var tmp Util.ResponseFromServer
			// decoder.Decode(&tmp)
			// log.Println("login http server", tmp)
			//tcp server response ok
			//and login success
			// if tmp.Success {
			// 	fmt.Println("get!!!", tmp)
			// 	return nil, true
			// }
			//something wrong in tcp server
			//or login fail
			fmt.Println("http read server:", data)
			// return Util.ResFailStr, false
			return nil, false
		}
		// case "home":
		// 	{
		// 		gob.Register(new(Util.ResponseFromServer))
		// 		gob.Register(new(Util.RealUser))
		// 		decoder := gob.NewDecoder(tcpconn)
		// 		var tmp Util.ResponseFromServer
		// 		err := decoder.Decode(&tmp)
		// 		if err != nil {
		// 			log.Println("home decode err")
		// 			panic(err)
		// 		}
		// 		// Util.FailFastCheckErr(err)
		// 		if !tmp.Success {
		// 			//token expires or not correct
		// 			if tmp.TcpData == nil {
		// 				return nil, false
		// 			}
		// 			//token pass
		// 			//but
		// 			fmt.Println("redis cache not update")
		// 			return tmp.TcpData.(*Util.RealUser), true
		// 		}
		// 		return tmp.TcpData.(*Util.RealUser), true
		// 	}
		// case "changeNickName":
		// 	{
		// 		decoder := gob.NewDecoder(tcpconn)
		// 		var tmp Util.ResponseFromServer
		// 		err := decoder.Decode(&tmp)
		// 		log.Println("changenickname recveive:", tmp)
		// 		if err != nil {
		// 			log.Println("changenickname decode err", err)
		// 		}
		// 		// Util.FailFastCheckErr(err)
		// 		if !tmp.Success {
		// 			//token expires or not correct
		// 			return nil, false
		// 		}

		// 		return nil, true
		// 	}
		// case "uploadAvatar":
		// 	{
		// 		gob.Register(new(Util.ResponseFromServer))
		// 		decoder := gob.NewDecoder(tcpconn)
		// 		var tmp Util.ResponseFromServer
		// 		err := decoder.Decode(&tmp)
		// 		if err != nil {
		// 			log.Println("uploadavatar decode err")
		// 			panic(err)
		// 		}
		// 		// Util.FailFastCheckErr(err)

		// 		return nil, true
		// 	}
		// case "logout":
		// 	{
		// 		gob.Register(new(Util.ResponseFromServer))

		// 		decoder := gob.NewDecoder(tcpconn)
		// 		var tmp Util.ResponseFromServer
		// 		err := decoder.Decode(&tmp)
		// 		if err != nil {
		// 			log.Println("logout decode err")
		// 			panic(err)
		// 		}
		// 		// Util.FailFastCheckErr(err)
		// 		return tmp.TcpData, tmp.Success
		// 	}
	}
	return nil, false

}

// func writeServer(data byte, tcpconn net.Conn, ctype string, readchan chan<- int, stopchan chan<- bool) (interface{}, bool) {
// 	// for {
// 	// 	tcpconn
// 	// }
// }

//for login Get render
func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		log.Println("login enter")
		t := template.Must(template.ParseFiles("../view/login.html"))
		usernamecookie, erruser := r.Cookie("username")
		tokencookie, errtoken := r.Cookie("token")
		//not found cookie
		if erruser != nil || errtoken != nil {
			fmt.Println("no cookie login", erruser)
			t.Execute(w, nil)
			return
		}
		// if errtoken != nil {
		// 	fmt.Println("no tokencookie", errtoken)
		// 	t.Execute(w, nil)
		// 	return
		// }

		//if found username and token (no matter right or wrong)
		//tmpcommand
		http.Redirect(w, r, "/Home", http.StatusFound)
		fmt.Println("login user cookie:", usernamecookie.Value)
		fmt.Println("login token cookie:", tokencookie.Value)
		t.Execute(w, nil)
	}
	//login authentication
	if r.Method == http.MethodPost {

		tcpconn, errget := connpool.Get()
		if errget != nil {
			panic(errget)
		}
		fmt.Println("tcp conn:", tcpconn.RemoteAddr().String())
		defer tcpconn.Close()

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
		tmpdata := &data.ToServerData{Ctype: proto.String("login"), Httpdata: tempuserData}
		tmpdataSend, _ := proto.Marshal(tmpdata)
		//----rpc call login------------
		// readchan := make(chan int)
		_, err := tcpconn.Write(tmpdataSend)
		if err != nil {
			fmt.Println("login write err:", err)
			panic(err)
		}

		//----------------------------
		//------old--------------
		// gob.Register(new(Util.User))
		// gob.Register(new(Util.RealUser))
		// gob.Register(new(Util.ToServerData))
		// encoder := gob.NewEncoder(tcpconn)
		// err := encoder.Encode(tmpdata)
		// if err != nil {
		// 	panic(err)
		// }
		//-------------------------

		fmt.Println("encode usename pwd:", tmpdata)
		// //loop to listen from server

		// for {
		// time.Sleep(2 * time.Second)
		_, successlogin := readServer(w, r, tcpconn, "login") //tcpconn or rpc.Con
		//success login

		if successlogin {
			fmt.Println("login success!!http")
			//, MaxAge: Util.CookieExpires
			log.Println("login cookie expr", Util.CookieExpires)

			cookie := http.Cookie{Name: "username", Value: username, Path: "/", Expires: Util.CookieExpires}
			http.SetCookie(w, &cookie)
			cookie = http.Cookie{Name: "token", Value: temptoken, Path: "/", Expires: Util.CookieExpires}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/Home", http.StatusFound)

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
// func homeHandler(w http.ResponseWriter, r *http.Request) {

// 	if r.Method == http.MethodGet {
// 		t := template.Must(template.ParseFiles("../view/Home.html"))
// 		cookieuser, erruser := r.Cookie("username")
// 		cookietoken, errtoken := r.Cookie("token")
// 		if erruser != nil {
// 			fmt.Println("cookie home user", erruser)
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			t.Execute(w, nil)
// 			return
// 		}
// 		if errtoken != nil {
// 			fmt.Println("cookie home token", errtoken)
// 			http.Redirect(w, r, "/login", http.StatusFound)
// 			t.Execute(w, nil)
// 			return
// 		}

// 		tcpconn, errget := connpool.Get()
// 		if errget != nil {
// 			panic(errget)
// 		}
// 		// Util.FailFastCheckErr(errget)
// 		log.Println("home rendering")

// 		//send to tcp server
// 		tokenwithusername := &Util.InfoWithUsername{Username: cookieuser.Value, Info: cookietoken.Value}
// 		tmpdata := &Util.ToServerData{}
// 		tmpdata.Ctype = "home"
// 		fmt.Println("http cookie ", cookieuser.Value)
// 		tmpdata.HttpData = tokenwithusername

// 		// fmt.Println("do sth ")
// 		gob.Register(new(Util.InfoWithUsername))
// 		gob.Register(new(Util.ToServerData))
// 		encoder := gob.NewEncoder(tcpconn)
// 		encoder.Encode(tmpdata)

// 		//loop listen response from tcp server
// 		// for {
// 		readchan := make(chan int)
// 		log.Println("home render loop", tmpdata)
// 		data, successHome := readServer(w, r, tcpconn, "home", readchan)
// 		//token correct
// 		if successHome {
// 			ruser := data.(*Util.RealUser)
// 			t.Execute(w, ruser)
// 			return
// 		}
// 		//token not correct
// 		//clear cookie and then redirect
// 		log.Println("token expires home page")
// 		setcookieuser := http.Cookie{
// 			Name:   "username",
// 			MaxAge: -1,
// 			Path:   "/",
// 		}
// 		setcookietoken := http.Cookie{
// 			Name:   "token",
// 			MaxAge: -1,
// 			Path:   "/",
// 		}
// 		http.SetCookie(w, &setcookieuser)
// 		http.SetCookie(w, &setcookietoken)
// 		http.Redirect(w, r, "/login", http.StatusFound)
// 		return
// 		// }
// 		// }

// 	}

// }

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
// 			// w.Write([]byte("wrong format!!!"))
// 			// http.Redirect(" ")
// 			http.Redirect(w, r, "/Home", http.StatusFound)
// 			return
// 		}

// 		// fmt.Fprintf(w, "%v", handler.Header)
// 		f, err := os.OpenFile(Util.UploadPath+filename, os.O_WRONLY|os.O_CREATE, 0666)
// 		if err != nil {
// 			fmt.Println("http openfile fail", err)
// 			return
// 		}
// 		defer f.Close()
// 		io.Copy(f, file)

// 		//get username from cookie

// 		tempAvatar := &Util.InfoWithUsername{Username: cookieuser.Value, Info: filename, Token: cookietoken.Value}
// 		uploadToServer := &Util.ToServerData{Ctype: "uploadAvatar", HttpData: tempAvatar}
// 		gob.Register(new(Util.InfoWithUsername))
// 		gob.Register(new(Util.ToServerData))
// 		encoder := gob.NewEncoder(tcpconn)
// 		uploaderr := encoder.Encode(uploadToServer)
// 		if uploaderr != nil {
// 			log.Println("upload encode err")
// 			panic(uploaderr)
// 		}
// 		// Util.FailFastCheckErr(uploaderr)
// 		//listen response from tcp server
// 		// for {
// 		readchan := make(chan int)
// 		_, successupload := readServer(w, r, tcpconn, "uploadAvatar", readchan)
// 		if successupload {
// 			http.Redirect(w, r, "/Home", http.StatusFound)
// 			return
// 		}
// 		//if db crash or token wrong
// 		http.Redirect(w, r, "/Home", http.StatusFound)
// 		// w.Write([]byte("upload file failed!"))
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
// 		gob.Register(new(Util.InfoWithUsername))
// 		gob.Register(new(Util.ToServerData))

// 		tempMap := &Util.InfoWithUsername{Username: cookieuser.Value, Info: newnickname, Token: cookietoken.Value}
// 		uploadToServer := &Util.ToServerData{Ctype: "changeNickName", HttpData: tempMap}
// 		encoder := gob.NewEncoder(tcpconn)
// 		err := encoder.Encode(uploadToServer)
// 		if err != nil {
// 			log.Println("change nickname encode err", err)
// 		}
// 		// Util.FailFastCheckErr(err)
// 		readchan := make(chan int)
// 		for {
// 			_, success := readServer(w, r, tcpconn, "changeNickName", readchan)
// 			if success {
// 				http.Redirect(w, r, "/Home", http.StatusFound)
// 				return
// 			}
// 		}
// 	}
// }
// func checkAndCreateFileName(oldName string) (newName string, isLegal bool) {
// 	ext := path.Ext(oldName)
// 	// uppercase
// 	//todo
// 	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
// 		newName = Util.GetFileName(oldName, ext)
// 		isLegal = true
// 	}
// 	return newName, isLegal
// }
