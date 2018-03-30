package main

import (
	"runtime"
	"os"
	"flag"
	"fmt"
	"context"
	"github.com/modeyang/LogCruiser/config"
	"log"
	"github.com/modeyang/LogCruiser/module"
)

var (
	confFile 	string
	help 	bool
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
}

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if help {
		flag.Usage()
	}
	conf, err := config.LoadFromFile(confFile)
	if err != nil {
		log.Println(err)
		return
	}
	module.InitModule()

	ctx , cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("start log process...")
	if err = conf.Start(ctx); err != nil {
		log.Println(err)
		return
	}

	log.Println("wait end ...")
	if err = conf.Wait(); err != nil {
		return
	}
	log.Println("end ...")
}
