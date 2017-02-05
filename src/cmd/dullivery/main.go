package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	auth "auth"
	daemon "daemon"

	api_ftp "api/ftp"
	api_health "api/health"
	api_posix "api/posix"
	server_ftp "server/ftp"
	server_health "server/health"
	server_posix "server/posix"
)

var (
	listen   = flag.String("listen", "0.0.0.0:5472", "listen address")
	content  = flag.String("content", "/content", "path to directory with content")
	statedb  = flag.String("state-db", "/var/lib/dullivery/server-state.json", "path to state file")
	hostname = flag.String("hostname", "", "Set hostname of dullivery server")
	logfile  = flag.String("log-file", "", "path to log file")
	pid      = flag.String("pid", "", "path to pid-file")
)

func main() {

	if !flag.Parsed() {
		flag.Parse()
	}

	if osHostname, err := os.Hostname(); err == nil && *hostname != "" {
		hostname = &osHostname
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	if *logfile != "" {
		if fd, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640); err != nil {
			fmt.Printf("Open log file error: %s\n", err.Error())
			os.Exit(1)
		} else {
			daemon.Daemonize(fd)
			log.SetOutput(fd)
		}
	}

	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		fmt.Printf("Start server error: %s\n", err.Error())
		os.Exit(2)
	}

	auth.SetGrpcServerName(*hostname)
	server := auth.NewGrpcServer()

	posix, err := server_posix.NewPosixServer(*content, *statedb)
	if err != nil {
		fmt.Printf("Start server error: %s\n", err.Error())
		os.Exit(3)
	}

	if *pid != "" {
		if err := ioutil.WriteFile(*pid, []byte(fmt.Sprintf("%d", os.Getpid())), 0640); err != nil {
			fmt.Printf("Write pid file error: %s\n", err.Error())
			os.Exit(4)
		}
	}

	api_posix.RegisterFileServer(server, posix)
	api_health.RegisterHealthServer(server, server_health.NewHealthServer(""))
	api_ftp.RegisterFileServer(server, server_ftp.NewFtpServer(*content))
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		fmt.Printf("Start server error: %s\n", err.Error())
		os.Exit(5)
	}

}
