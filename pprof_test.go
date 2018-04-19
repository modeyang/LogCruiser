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
)

func counter() {
	list := []int{1}
	c := 1
	for i := 0; i < 100; i++ {
		length := httpGet()
		c += length
		list = append(list, c)
	}
	fmt.Println(c)
	fmt.Println(list[0])
}

func httpGet() int{
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
	return len(body)
}

func worker(wg *sync.WaitGroup) {
	counter()
	wg.Done()
}


func TestPProf(T *testing.T) {
	profile := "test.prof"
	f, err := os.Create(profile)
	if err != nil {
		T.Error(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	var wg sync.WaitGroup
	wg.Add(1)
	for i:=0; i<100; i++ {
		go worker(&wg)
	}
	wg.Wait()
	log.Println("end..")

}
