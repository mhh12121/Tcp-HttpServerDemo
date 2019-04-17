package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	testName := "testData"
	var f *os.File
	var errF error
	if !fileIsExist(testName) {
		f, errF = os.Create(testName)
		if errF != nil {
			panic(errF)
		}
	} else {
		f, errF = os.OpenFile(testName, os.O_APPEND, 0666)
		if errF != nil {
			panic(errF)
		}
	}
	w := bufio.NewWriter(f)
	for i := 1; i <= 200; i++ {
		num, err := w.WriteString("username=mhh" + strconv.Itoa(i) + "&password=a123456\n")
		if err != nil {
			log.Println(num)
			panic(err)
		}
	}
	w.Flush()
	f.Close()

}
func fileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
