package Util

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"

	"strconv"
	"time"
)

// func init() {
// 	file, err := os.Open("../config/config.json")
// 	if err != nil {
// 		panic("file config wrong")
// 	}
// 	defer file.Close()
// 	decoder := json.NewDecoder(file)
// }

//to tcp server

const (
	TimeoutDuration = 5 * time.Minute
	UploadPath      = "../images/"
)

var TokenExpires = int64(1e11)
var CookieExpires = time.Now().Add(1 * time.Hour)

var ErrCode map[int]string = map[int]string{
	4001: "mysql record not found",
	4002: "redis record not found",
	4003: "cookie record not found",
	5001: "unmarshal error",
	5002: "no connection tcp",
	5003: "no connection http",
	6001: "username error",
	6002: "password error",
}

const (
	LOGINCODE    = 2
	LOGOUTCODE   = 4
	HOMECODE     = 6
	UPDATENICK   = 8
	UPLOADAVATAR = 10
	COMPRESS     = 1
	NOCOMPRESS   = 0
)

func IntToBytes(n int) []byte {
	res := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, res)
	return bytesBuffer.Bytes()
}
func BytestoInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)

}

//rename the uploaded files
func GetFileName(fileName string, ext string) string {
	h := md5.New()
	h.Write([]byte(fileName + strconv.FormatInt(time.Now().Unix(), 10)))
	return hex.EncodeToString(h.Sum(nil)) + ext
}
