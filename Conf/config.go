package Conf

import (
	"encoding/json"
	"log"
	"os"
)

type Tconfig struct {
	Connect struct {
		Httphost string
		Tcphost  string
		Httpport string
		Tcpport  string
	}
	Chanpool struct {
		Initsize int
		Maxsize  int
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		Db       int
		Poolsize int
	}
	Mysql struct {
		Host     string
		Port     string
		Username string
		Password string
		Db       string
	}
}

var Config *Tconfig

func LoadConf(confpath string) {
	if Config == nil {
		Config = &Tconfig{}

		file, err := os.Open(confpath)
		defer file.Close()
		if err != nil {
			log.Println("loading conf err", err)
			return
		}
		decoder := json.NewDecoder(file)
		errdecode := decoder.Decode(Config)
		if errdecode != nil {
			log.Println("parsing json fail:", errdecode)
			return
		}
	}
}
