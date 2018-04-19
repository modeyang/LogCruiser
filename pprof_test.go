package main

import (
	"testing"
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"sync"
	_ "net/http/pprof"
)

func counter() {
	list := []int{1}
	c := 1
	lenChan := make(chan int, 1)
	for i := 0; i < 100; i++ {
		go httpGet(lenChan)
		select {
		case length := <- lenChan:
			c += length
		}
		list = append(list, c)
	}
	fmt.Println(c)
	fmt.Println(list[0])
}

func httpGet(lenChan chan int) int{
	resp, err := http.Get("http://www.163.com")
	if err != nil {
		log.Println(err)
		return 0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0
	}
	lenChan <- len(body)
	return len(body)
}

func worker(wg *sync.WaitGroup) {
	counter()
	wg.Done()
}

func testHttpProfile() {
	//http.HandleFunc("/debug/pprof/", tpprof.Index)
	//http.HandleFunc("/debug/pprof/cmdline", tpprof.Cmdline)
	//http.HandleFunc("/debug/pprof/profile", tpprof.Profile)
	//http.HandleFunc("/debug/pprof/symbol", tpprof.Symbol)
	//http.HandleFunc("/debug/pprof/trace", tpprof.Trace)
	log.Println(http.ListenAndServe("0.0.0.0:6100", nil))
}

func TestPProf(T *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	profile := "test.prof"
	f, err := os.Create(profile)
	if err != nil {
		T.Error(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	go testHttpProfile()

	var wg sync.WaitGroup
	wg.Add(10)
	for i:=0; i<10; i++ {
		go worker(&wg)
	}
	wg.Wait()
	log.Println("end..")

	//time.Sleep(60 * time.Second)
}
