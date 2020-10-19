package main

import (
	"context"
	"encoding/json"
	"Tcp-HttpServerDemo/Conf"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"path"
	"runtime"
	"time"

	data "Tcp-HttpServerDemo/Data"
	Util "Tcp-HttpServerDemo/Util"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	// pool "gopkg.in/fatih/pool.v2"
	// pool "Tcp-HttpServerDemo/ServerPool"
)

// var connpool *pool.GRPCPool
var globalcon net.Conn
var addr string

// var conn *grpc.ClientConn

func init() {

}

var loginTemplate *template.Template
var homeTemplate *template.Template

func main() {
	_, filepath, _, _ := runtime.Caller(0)
	p := path.Dir(filepath)
	p = path.Dir(p)

	log.Println("log path", p)
	// var err error

	// defer conn.Close()
	Conf.LoadConf(p + "/Conf/config.json")

	// tcpconn, err = net.Dial("tcp", Util.Tcpaddress+":"+Util.Tcpport)
	// tcpconn.SetReadDeadline(time.Now().Add(Util.TimeoutDuration))

	addr = Conf.Config.Connect.Tcphost + ":" + Conf.Config.Connect.Tcpport
	// options := &pool.Options{
	// 	InitTargets:  []string{addr},
	// 	InitCap:      5,
	// 	MaxCap:       100,
	// 	DialTimeout:  time.Second * 5,
	// 	IdleTimeout:  time.Second * 60,
	// 	ReadTimeout:  time.Second * 5,
	// 	WriteTimeout: time.Second * 5,
	// }
	// var perr error
	// connpool, perr = pool.NewGRPCPool(options, grpc.WithInsecure())
	// if perr != nil {
	// 	panic(perr)
	// }
	// defer connpool.Close()
	// var err error
	// conn, err = grpc.Dial(addr, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("logout client fail", err)
	// }
	// tcpAddr, addErr := net.ResolveTCPAddr("tcp4", addr)
	// if addErr != nil {
	// 	fmt.Fprintf(os.Stderr, "Fatal error: %s", addErr.Error())
	// 	os.Exit(1)
	// }
	// factory := func() (net.Conn, error) {
	// 	return net.DialTCP("tcp", nil, tcpAddr)
	// }
	// var err error
	// connpool, err = pool.NewChannelPool(Conf.Config.Chanpool.Initsize, Conf.Config.Chanpool.Maxsize, factory)
	// if err != nil {
	// 	panic(err)
	// }

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
	httpTmp := &http.Server{
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         Conf.Config.Connect.Httphost + ":" + Conf.Config.Connect.Httpport,
	}
	errhttp := httpTmp.ListenAndServe()
	log.Fatalf("http listen server: %v", errhttp)

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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodPost {
		//todo
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Println("----------------logout-------------------")

		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("logout client fail", err)
		}
		defer conn.Close()

		c := data.NewAuthenticateClient(conn)

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

		// for {
		//go to tcp to invalid the cache
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		res, errR := c.Logout(ctx, tmp)
		if errR != nil {
			log.Fatalf("logout res failed", errR)
		}
		successlogout := res.GetSuccess()
		// _, successlogout := readServer(w, r, tcpconn, Util.LOGOUTCODE)
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

// func do(i int) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }

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

		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("grpc login fail: %v", err)
		}
		defer conn.Close()

		c := data.NewAuthenticateClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println("enter!!!!!!")
		username := r.FormValue("username")
		password := r.FormValue("password")
		// fmt.Println("front username:", username)
		// fmt.Println("front pwd:", password)

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
		res, errR := c.Login(ctx, tmpdata)
		if errR != nil {
			log.Fatalf("login response failed %v", errR)
		}
		successlogin := res.GetSuccess()
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
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("home getfrom server fail", err)
		}
		defer conn.Close()

		c := data.NewAuthenticateClient(conn)

		log.Println("home rendering")

		//send to tcp server
		tokenwithusername := &data.InfoWithUsername{Username: proto.String(cookieuser.Value), Token: proto.String(cookietoken.Value)}
		tokenwithusernameData, terr := proto.Marshal(tokenwithusername)
		if terr != nil {
			panic(terr)
		}
		tmp := &data.ToServerData{Ctype: proto.Int32(Util.HOMECODE), Httpdata: tokenwithusernameData}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, errR := c.Home(ctx, tmp)
		if errR != nil {
			log.Fatalf("home response failed %v", errR)
		}
		successHome := res.GetSuccess()
		//token correct
		if successHome {
			tmpData := res.GetTcpData()
			ruser := &data.RealUser{}
			errU := proto.Unmarshal(tmpData, ruser)
			if errU != nil {
				log.Fatalf("home get user msg fail %v", errU)
			}
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
// if r.Method != http.MethodPost {
// 	w.WriteHeader(http.StatusMethodNotAllowed)
// 	return
// }
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
// if r.Method != http.MethodPost {
// 	w.WriteHeader(http.StatusMethodNotAllowed)
// 	return
// }
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
