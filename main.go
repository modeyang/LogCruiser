package main

import (
	"runtime"
	"os"
	"flag"
	"fmt"
	"context"
	"log"
	"github.com/modeyang/LogCruiser/config"
	"github.com/modeyang/LogCruiser/module"
	"runtime/pprof"
	_ "net/http/pprof"
	"net/http"
)

var (
	confFile 	string
	help 		bool
	cpuProfile 	string
	port 		int
)

func usage() {
	fmt.Fprintf(os.Stderr, `
		Usage: cruiser -s <config>
	`)
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.StringVar(&confFile, "c", "Log.yml", "config file")
	flag.BoolVar(&help, "h", false, "tool help")
	flag.StringVar(&cpuProfile, "cpuprofile", "", "write cpu profile to file")
	flag.IntVar(&port, "port", 6100, "listen port")
}

func initHttpProfile() {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	for {
		http.ListenAndServe(address, nil)
	}
}

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if help {
		flag.Usage()
	}
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	conf, err := config.LoadFromFile(confFile)
	if err != nil {
		log.Println(err)
		return
	}
	module.InitModule()

	ctx , cancel := context.WithCancel(context.TODO())
	defer cancel()

	log.Println("start log process...")
	if err = conf.Start(ctx); err != nil {
		log.Println(err)
		return
	}
	go initHttpProfile()
	log.Println("wait end ...")
	if err = conf.Wait(); err != nil {
		return
	}
	log.Println("end ...")
}
