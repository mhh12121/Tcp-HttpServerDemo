package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	serverAddr := "http://localhost:8080/login"
	var transport http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          int(10),
		MaxIdleConnsPerHost:   int(10),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
	}
	data := url.Values{}
	var buffer bytes.Buffer
	buffer.WriteString("mhh12345")

	username := buffer.String()

	data.Set("username", username)
	data.Set("password", "a123456")
	log.Println("account", data.Encode())
	// data.Set("nickname", "newbot")
	// req, err := http.NewRequest("GET", serverAddr, bytes.NewBufferString(data.Encode()))
	req, err := http.NewRequest("POST", serverAddr, bytes.NewBufferString(data.Encode()))
	// req.AddCookie(&http.Cookie{Name: "username", Value: username, Expires: time.Now().Add(120 * time.Second), Path: "/"})
	// req.AddCookie(&http.Cookie{Name: "token", Value: "test", Expires: time.Now().Add(120 * time.Second), Path: "/"})

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // This makes it workparam=value
	req.Close = true
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("simulate send request err", req, resp, err)
	}
	_, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		log.Println("ioutil req err:", req, err1)
	}
	resp.Body.Close()

}

func TestChangeNickName(t *testing.T) {
	//need set cookie mannually
	//username : mhh123456
	//cookie : test
	var transport http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          int(10),
		MaxIdleConnsPerHost:   int(10),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	// var buffer bytes.Buffer
	// buffer.WriteString("newnickname1")

	changeNikcNameAddr := "http://localhost:8080/Home/change"
	data := url.Values{}

	newnickname := "newnickname123"
	username := "mhh123456"

	data.Set("newnickname", newnickname)
	log.Println("updatenickname", data.Encode())
	// data.Set("nickname", "newbot")
	// req, err := http.NewRequest("GET", serverAddr, bytes.NewBufferString(data.Encode()))
	req, err := http.NewRequest("POST", changeNikcNameAddr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Println(err)
	}
	req.AddCookie(&http.Cookie{Name: "username", Value: username, Expires: time.Now().Add(120 * time.Second), Path: "/"})
	req.AddCookie(&http.Cookie{Name: "token", Value: "test", Expires: time.Now().Add(120 * time.Second), Path: "/"})

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // This makes it workparam=value

	resp, err := client.Do(req)
	if err != nil {
		log.Println("simulate send request err", req, resp, err)
	}
	_, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		log.Println("ioutil req err:", req, err1)
	}
	resp.Body.Close()

}
