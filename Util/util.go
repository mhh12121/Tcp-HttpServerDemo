package Util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

//to tcp server
const Tcpaddress = "localhost"
const Tcpport = "8081"
const Httpport = "8080"
const TimeoutDuration = 5 * time.Minute
const UploadPath = "../images/"
const ResSuccessStr = "success"
const ResFailStr = "fail"
const ResWrongStr = "wrong password or account"

//redis info
const RedisAddr = "localhost:6379"

var TokenExpires = int64(1e11)
var CookieExpires = time.Now().Add(1 * time.Hour)

// var TokenExpires = time.Now().Add(1 * time.Minute)

//uniform data to tcp server
type ToServerData struct {
	Ctype    string
	HttpData interface{}
}

//Uniform data from tcp server
type ResponseFromServer struct {
	Success bool
	TcpData interface{}
}

//success response from server

//RealUser is the home page data
type RealUser struct {
	Username string
	Nickname string
	Avatar   string
}

//User is login data
type User struct {
	Username string
	Password string
	Token    string
}

//Info is for changing avatar data and nickname
type InfoWithUsername struct {
	Username string
	Info     interface{}
	Token    string
}

func FailSafeCheckErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
func FailFastCheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//rename the uploaded files
func GetFileName(fileName string, ext string) string {
	h := md5.New()
	h.Write([]byte(fileName + strconv.FormatInt(time.Now().Unix(), 10)))
	return hex.EncodeToString(h.Sum(nil)) + ext
}
