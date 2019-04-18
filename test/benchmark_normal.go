package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func loginHandleHttp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// io.WriteString(w)
	// io.WriteString()
}
func BenchmarkWithHttp(serverAddr string, c int, isRan bool) time.Duration {
	var wg sync.WaitGroup

	cliRoutine := func(no int) {
		data := url.Values{}
		var buffer bytes.Buffer
		buffer.WriteString("mhh")
		buffer.WriteString(strconv.Itoa(no))

		username := buffer.String()
		data.Set("username", username)
		data.Set("password", "a123456")
		// handler := func(w http.ResponseWriter, r *http.Request) {
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Header().Set("Content-Type", "application/json")
		// 	io.WriteString()
		// }
		// req, err := http.NewRequest("POST", serverAddr, bytes.NewBufferString(data.Encode()))
		resp, err := http.PostForm(serverAddr, data)

		if err != nil {
			// if err == io.EOF {
			// 	fmt.Println("test post form eof")
			// } else {
			// 	panic(err)
			// }

			// handle error
		}

		// defer resp.Body.Close()
		_, errx := ioutil.ReadAll(resp.Body)
		if errx != nil {
			panic(errx)
			// handle error
		}
		resp.Body.Close()
		// fmt.Println(body)
		wg.Done()
	}
	for i := 1; i < c; i++ {
		wg.Add(1)
		go cliRoutine(i)
	}
	start := time.Now()
	wg.Wait()
	return time.Since(start)

}
