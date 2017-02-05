package dsl

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/yuin/gopher-lua"
	"google.golang.org/grpc/reflection"

	api_health "api/health"
	api_posix "api/posix"
	auth "auth"
	server_health "server/health"
	server_posix "server/posix"
)

func getTestingContentDir(suff string) string {
	tmpdir := os.Getenv("TEST_TMPDIR")
	if tmpdir == "" {
		if runtime.GOOS == "linux" {
			tmpdir = filepath.Join("/tmp/dullivery", suff)
		} else {
			tmpdir = filepath.Join(`c:\tmp`, suff)
		}
	}
	os.MkdirAll(tmpdir, 0750)
	return tmpdir
}

func prepareTestingBigFile(suff string) {
	bigfile := filepath.Join(getTestingContentDir(suff), "my_file")
	fd, err := os.Create(bigfile)
	if err != nil {
		panic(err)
	}
	fd.Seek(1024*1024*25, 0)
	fd.Write([]byte("ok"))
	fd.Sync()
	fd.Close()
}

func runTestingServer(address, suff string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	server := auth.NewGrpcServer()
	posix, err := server_posix.NewPosixServer(getTestingContentDir(suff), filepath.Join(getTestingContentDir(suff), "state.json"))
	if err != nil {
		panic(err)
	}
	api_posix.RegisterFileServer(server, posix)
	api_health.RegisterHealthServer(server, server_health.NewHealthServer(""))
	reflection.Register(server)
	server.Serve(listener)
}

func TestDsl(t *testing.T) {

	go runTestingServer("localhost:3131", "1")
	go runTestingServer("localhost:3132", "2")
	prepareTestingBigFile("1")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("§ ")

	time.Sleep(100 * time.Millisecond)

	storageDir := filepath.Join(getTestingContentDir("db"))

	state := lua.NewState()
	Register(state, &Config{LogFd: os.Stdout, StorageDir: storageDir})
	if err := state.DoFile("dsl_test.lua"); err != nil {
		t.Fatalf("error: %s\n", err.Error())
	}
}
