package ftp

import (
	"log"
	"net"
	"os"
	"runtime"

	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "api/ftp"
	ftpd "server/ftp/ftpd"
)

const testingServerAddress = "127.0.0.1:3030"

func runTestFtpServer() {

	ln, err := net.Listen("tcp", "127.0.0.1:2121")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		c := ftpd.NewFTPConn(conn, getTestingContentDir())
		go ftpd.HandleConnection(c)
	}

}

func getTestingContentDir() string {
	tmpdir := os.Getenv("TEST_TMPDIR")
	if tmpdir == "" {
		if runtime.GOOS == "linux" {
			tmpdir = "/tmp/dullivery"
		} else {
			tmpdir = `c:\tmp`
		}
	}
	os.MkdirAll(tmpdir, 0750)
	return tmpdir
}

func startTestingServer() {
	listener, err := net.Listen("tcp", testingServerAddress)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	ftpApiServer := NewFtpServer(getTestingContentDir())
	api.RegisterFileServer(server, ftpApiServer)
	reflection.Register(server)
	server.Serve(listener)
}

func TestPosix(t *testing.T) {

	go runTestFtpServer()
	go startTestingServer()

	conn, err := grpc.Dial(testingServerAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial: %s\n", err.Error())
	}
	defer conn.Close()
	client := api.NewFileClient(conn)

	ftpConn := &api.Conn{
		Host: "127.0.0.1:2121",
	}
	filter := &api.Filter{
		PathMatch:  "*",
		Connection: ftpConn,
	}
	if list, err := client.Find(context.Background(), filter); err != nil {
		t.Fatal(err)
	} else {
		log.Printf("list: %v\n", list)
	}

}
