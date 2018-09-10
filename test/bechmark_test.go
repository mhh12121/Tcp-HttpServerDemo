package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

// func benchmarkLoginReq(serverAddr string, n, c int32, isRan bool) (elapsed time.Duration) {
func BenchmarkLoginReq(b *testing.B) {
	// b.N = 1
	serverAddr := "http://localhost:8080/login"
	// n := int32(2)
	c := int32(200) //concurrency
	// isRan := true
	readyGo := make(chan bool)
	// test := make(chan int)
	var wg sync.WaitGroup
	log.Println("concurrency", c)
	wg.Add(int(c))

	// remaining := n

	var transport http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          int(c),
		MaxIdleConnsPerHost:   int(c),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	cliRoutine := func(no int) {
		// for atomic.AddInt32(&remaining, -1) > 0 {
		// continue

		data := url.Values{}

		var buffer bytes.Buffer
		buffer.WriteString("mhh")
		// rand
		// if isRan {
		// 	buffer.WriteString(strconv.Itoa(rand.Intn(1000000)))
		// } else {
		// 	buffer.WriteString("1")
		// }
		//here random number may be the same, so that it will regenerate token
		//however, token is checked in TCP server,if its not the same as the one in redis with the same username
		//then it will throw errors and redirect to login page
		buffer.WriteString(strconv.Itoa(no))
		log.Println("write string", no)
		username := buffer.String()

		data.Set("username", username)
		data.Set("password", "a123456")
		log.Println("account", data.Encode())
		// data.Set("nickname", "newbot")
		// req, err := http.NewRequest("GET", serverAddr, bytes.NewBufferString(data.Encode()))
		req, err := http.NewRequest("POST", serverAddr, bytes.NewBufferString(data.Encode()))
		// req.AddCookie(&http.Cookie{Name: "username", Value: username, Expires: time.Now().Add(120 * time.Second), Path: "/"})
		// req.AddCookie(&http.Cookie{Name: "token", Value: "test", Expires: time.Now().Add(120 * time.Second), Path: "/"})
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") //param=value
		if err != nil {
			log.Println(err)
		}
		<-readyGo

		resp, err := client.Do(req)

		if err != nil {
			log.Println("simulate", err)
		}
		_, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			log.Println("ioutil req err:", req, err1)
		}

		defer resp.Body.Close()
		// }

		wg.Done()
	}
	// for atomic.AddInt32(&c, -1) > 0 {

	// }
	for i := int32(0); i < c; i++ {
		// log.Println("run :", i)
		// test <- (int(i) + 1)
		tmp := int(i) + 1
		go cliRoutine(tmp)
	}
	// close(test)
	close(readyGo)
	start := time.Now()
	wg.Wait()
	log.Println(time.Since(start))
	// return time.Since(start)
}
