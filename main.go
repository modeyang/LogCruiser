package main

import (
	"runtime"
	"os"
	"flag"
	"fmt"
	"context"
	"github.com/modeyang/LogCruiser/config"
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
	flag.StringVar(&confFile, "c", "", "config file")
	flag.BoolVar(&help, "h", false, "tool help")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if help {
		flag.Usage()
	}
	conf, err := config.LoadFromFile(confFile)
	if err != nil {
		return
	}
	ctx , cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = conf.Start(ctx); err != nil {
		return
	}

	if err = conf.Wait(); err != nil {
		return
	}
}
