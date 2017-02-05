package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	daemon "daemon"
	store "store"
	www "web"
)

var (
	listen  = flag.String("listen", "0.0.0.0:8080", "listen address")
	static  = flag.String("static", "/usr/share/dullivery/static", "www static files")
	data    = flag.String("data", "/var/lib/dullivery", "data directory")
	logfile = flag.String("log-file", "", "path to log file")
	pid     = flag.String("pid", "", "path to pid-file")
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	var fd *os.File
	var err error
	if *logfile != "" {
		if fd, err = os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640); err != nil {
			fmt.Printf("Open log file error: %s\n", err.Error())
			os.Exit(1)
		} else {
			log.SetOutput(fd)
			daemon.Daemonize(fd)
		}
	}

	if *pid != "" {
		if err = ioutil.WriteFile(*pid, []byte(fmt.Sprintf("%d", os.Getpid())), 0640); err != nil {
			fmt.Printf("Write pid file error: %s\n", err.Error())
			os.Exit(2)
		}
	}

	storage, err := store.NewStorage(*data)
	if err != nil {
		fmt.Printf("Start web server error: %s\n", err.Error())
		os.Exit(3)
	}
	store.Box = storage

	if err := www.Serve(*listen, *static, fd); err != nil {
		fmt.Printf("Start web server error: %s\n", err.Error())
		os.Exit(4)
	}
}
