package dao

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCheck(t *testing.T) {
	InitDB()
	var tests = []struct {
		username string
		password string
	}{
		{
			"mhh1", "a123456",
		}, {
			"mhh2", "a123456",
		}, {
			"mhha", "a123456",
		}, {
			"mhh3", "a1234567",
		}, {
			"m or 1", "a123456",
		},
	}
	for _, test := range tests {
		ok, err := Check(test.username, test.password)
		if err != nil || !ok {
			log.Printf("Test login auth fail,username:%s,password:%s", test.username, test.password)
		}
	}
}

func TestUpdateNickName(t *testing.T) {
	InitDB()
	var tests = []struct {
		username    string
		newnickname string
	}{
		{
			"mhh1", "dou1",
		}, {
			"mhh2", "dou2",
		}, {
			"mhha", "dou3",
		}, {
			"mhh3", "dou4",
		}, {
			"m or 1", "dou5",
		},
	}
	for _, test := range tests {
		ok, err := UpdateNickname(test.username, test.newnickname)
		if err != nil || !ok {
			log.Printf("Test updatenickname fail,username:%s,password:%s", test.username, test.newnickname)
		}
	}
}
