package Util

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

const (
	PACK_CLIENT      = "CLIENT"
	PACK_HEARTBEAT   = "HEART"
	PackHeaderLength = 8
	PackDataLength   = 4
	PackZipLength    = 1
)

func Pack(source string, message []byte, IsCompress bool) []byte {
	fmt.Println("pack")
	bitCompress := byte(1)
	if !IsCompress {
		bitCompress = byte(0)
	}
	sourceMsg := make([]byte, PackHeaderLength)

	copy(sourceMsg[:], []byte(source)) //1.pack header content
	fmt.Println("header", string(sourceMsg))
	fmt.Println("header length", len(sourceMsg))
	compressMsg := make([]byte, PackZipLength)
	compressMsg[0] = bitCompress
	sourceMsg = append(sourceMsg, IntToBytes(len(message))...) //2. pack Data length
	sourceMsg = append(sourceMsg, compressMsg...)              //3.pack compress info
	sourceMsg = append(sourceMsg, message...)                  //4.pack RealData
	return sourceMsg
}

/*
Use length description:
-----------------------+---------------------------------+------------------------+-----------+-------------------
header(8bytes)         +length(4 bytes)=Data1's length   +  IF COMPRESSED(1 byte) +   Data1  +   length(2)//todo
-----------------------+---------------------------------+------------------------+------------+-------------------
*/
// func Unpack(buffer []byte) []byte { //, readerChannel chan []byte
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	// defer wg.Done()
	length := len(buffer)
	fmt.Println("tcp buffer length", length)
	var i int
	for i = 0; i < length; i++ {
		fmt.Println("------------------unpack begin-----------------------------")
		if length < i+PackHeaderLength+PackDataLength+PackZipLength {
			// <-readerChannel
			fmt.Println("unpack length not complete")
			break
		}
		x := string(bytes.Trim(buffer[i:i+PackHeaderLength], "\x00")) //remove unknown "\x00" in byte array
		// fmt.Println("buff pack header", len(string(x)))
		//real msg pack
		if strings.Compare(x, PACK_CLIENT) == 0 {
			fmt.Println("msg from client ")
			datalength := BytestoInt(buffer[i+PackHeaderLength : i+PackHeaderLength+PackDataLength]) //get data length
			fmt.Println("tcp realdata length", datalength)
			if length < i+PackHeaderLength+PackDataLength+PackZipLength+datalength { //package.length > buffer.length
				fmt.Println("pack length > buffer length")
				break
			}
			if BytestoInt(buffer[i+PackHeaderLength+PackDataLength:i+PackHeaderLength+PackDataLength+PackZipLength]) == COMPRESS {
				fmt.Println("compresss")
				//do extraction//todo
			} else {
				fmt.Println("nocompresss")
				data := buffer[i+PackHeaderLength+PackDataLength+PackZipLength : i+PackHeaderLength+PackDataLength+PackZipLength+datalength]
				log.Println("data------------------>", data)
				log.Println("readerchannel before data------------------>", len(readerChannel))
				// return data
				readerChannel <- data

				log.Println("readerchannel after data------------------>", len(readerChannel))
				log.Println("data in channel------------------>", data)
				i += PackHeaderLength + PackDataLength + PackZipLength + datalength
			}

		} else if strings.Compare(x, PACK_HEARTBEAT) == 0 {
			//todo
			fmt.Println("msg from heartbeat ")
		} else {
			fmt.Println("w's the hell wrong with it")
		}

	}
	if i == length { //no data
		return make([]byte, 0)
	}
	return buffer[i:]

}

// /**/
// func UnpackHttp(buffer []byte) []byte {
// 	length := len(buffer)
// }

//todo
func Depress(compress uint16) bool {
	fmt.Println("compress mark ", compress)
	return false
}

// for {

// 	// size, cerr := con.Read(buff)
// 	fmt.Println(buff[:size])
// 	if cerr != nil {
// 		if cerr == io.EOF {
// 			fmt.Println("eof read ")
// 			break
// 		}
// 		fmt.Println("buferr", cerr)
// 		panic(cerr)
// 		// break
// 	}
// 	realData = append(realData, buff[:size]...)
// 	// if size == 0 {
// 	// 	fmt.Println("no message")
// 	// 	break
// 	// }
// }
