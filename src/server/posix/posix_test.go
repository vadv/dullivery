package posix

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "api/posix"
	ut "utils"
)

const testingFileServerAddress = "127.0.0.7:7775"

func prepareTestingBigFile() {
	bigfile := filepath.Join(getTestingContentDir(), "big_file_1")
	fd, err := os.Create(bigfile)
	if err != nil {
		panic(err)
	}
	fd.Seek(1024*1024*25, 0)
	fd.Write([]byte("ok"))
	fd.Sync()
	fd.Close()
	ioutil.WriteFile(bigfile+".md5", []byte(`cca460c24c46f11cabd200ce41daad57 file`), 0640)
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

func getTestingStateFile() string {
	return filepath.Join(getTestingContentDir(), "state.json")
}

func readMd5Sum(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	content := strings.Split(string(data), " ")
	result := content[0]
	return result
}

func startTestingServer() {
	listener, err := net.Listen("tcp", testingFileServerAddress)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	posix_file_server, err := NewPosixServer(getTestingContentDir(), getTestingStateFile())
	if err != nil {
		panic(err)
	}
	api.RegisterFileServer(server, posix_file_server)
	reflection.Register(server)
	server.Serve(listener)
}

func TestPosix(t *testing.T) {

	prepareTestingBigFile()
	go startTestingServer()

	time.Sleep(time.Second)

	conn, err := grpc.Dial(testingFileServerAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial: %s\n", err.Error())
	}
	defer conn.Close()
	client_file := api.NewFileClient(conn)

	// загружаем файл
	fd, err := os.Open(filepath.Join(getTestingContentDir(), "big_file_1"))
	if err != nil {
		t.Fatalf("open data file: %s", err.Error())
	}
	hash, err := ut.Md5FD(fd)
	if err != nil {
		t.Fatalf("calculate hash error: %s", err.Error())
	}
	t.Log("calculate hash sum of big file done")

	// проверяем чексумму
	realHash := readMd5Sum(filepath.Join(getTestingContentDir(), "big_file_1.md5"))
	if hash != realHash {
		t.Fatal("calculate hash error expected", realHash, "got", hash)
	}
	fd.Seek(0, 0)
	if info, err := StreamFileClient(client_file, fd, "new_file", hash); err != nil {
		t.Fatalf("stream file error: %s\n", err.Error())
	} else {
		if info.State != api.Info_OK {
			t.Error("upload hash error expected", api.Info_OK, "got", info.State)
		}
		if info.Path != "new_file" {
			t.Error("upload name error: expected", "new_file", "got", info.Path)
		}
	}
	fd.Close()
	t.Log("new file uploaded")

	// находим старый файл
	if list, err := client_file.Find(context.Background(), &api.Filter{PathMatch: "^big_file_1$"}); err != nil {
		t.Fatal(err.Error())
	} else {
		if len(list.Files) == 0 {
			t.Fatal("can't recieve information about old file")
		}
		file := list.Files[0]
		if file.Md5 != realHash {
			t.Error("find md5: expected", realHash, "got", file.Md5)
		}
	}
	t.Log("find old file completed")

	// сохраняем этот файл
	new_file := filepath.Join(getTestingContentDir(), "big_file_2")
	fdStream, err := os.Create(new_file)
	if err != nil {
		t.Fatalf("create file: %s\n", err.Error())
	}
	if err := ReceiveFileClient(client_file, "big_file_1", fdStream); err != nil {
		t.Fatalf("failed to rec file: %s\n", err.Error())
	}
	fdStream.Seek(0, 0)
	if new_file_hash, err := ut.Md5FD(fdStream); err != nil {
		t.Fatalf("hash: %s\n", err.Error())
	} else {
		if new_file_hash != realHash {
			t.Error("stream md5: expected", realHash, "got", new_file_hash)
		}
	}
	t.Log("stream file completed")

}
