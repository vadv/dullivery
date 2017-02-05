package dsl

import (
	"fmt"
	"log"

	"github.com/yuin/gopher-lua"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api_ftp "api/ftp"
	api_posix "api/posix"
	auth "auth"
)

type dslFile struct {
	Name string
	Size int64
}

func (d *Config) dslFind(L *lua.LState) int {

	src, err := parseUrl(L.CheckString(1))
	if err != nil {
		log.Printf("[ERROR] source url: %s\n", err.Error())
		L.ArgError(1, fmt.Sprintf("can't parse source url: %s", err.Error()))
		return -1
	}

	server := src.Host
	log.Printf("[INFO] run on host: %s command: find on %s path: %s\n", server, src.Scheme, src.Path)
	if src.Scheme == DSLFtpScheme {
		log.Printf("[INFO] Choose localhost as server for ftp search")
		server = auth.GetGrpcServerName() // по ftp ищем с текущего сервера
	}
	// make connection
	conn, err := grpc.Dial(
		src.PosixAddress(),
		grpc.WithPerRPCCredentials(auth.GrpcAuthFor(server)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[ERROR] dial to %s: %s\n", server, err.Error())
		L.RaiseError("dial %s: %s", server, err.Error())
		return -1
	}
	defer conn.Close()

	var files []*dslFile
	switch src.Scheme {
	case DSLFtpScheme:
		files, err = listFilesFtp(src, conn)
	case DSLPosixScheme:
		files, err = listFilesPosix(src, conn)
	default:
		panic(fmt.Sprintf("unknown scheme for find: %s", src.Scheme))
	}
	if err != nil {
		L.RaiseError("find error: %s", err.Error())
		return -1
	}
	// transfer result files to lua structs:
	result := L.NewTable()
	for _, file := range files {
		newFile := L.NewTable()
		L.SetField(newFile, "name", lua.LString(file.Name))
		L.SetField(newFile, "size", lua.LNumber(file.Size))
		result.Append(newFile)
	}
	L.Push(result)
	return 1
}

func listFilesFtp(src *DSLUrl, conn *grpc.ClientConn) ([]*dslFile, error) {
	client := api_ftp.NewFileClient(conn)
	filter := &api_ftp.Filter{PathMatch: src.Path}
	if src.Url.User != nil {
		filter.Connection = &api_ftp.Conn{Host: src.Path, User: src.Url.User.Username()}
		if passwd, ok := src.Url.User.Password(); ok {
			filter.Connection.Password = passwd
		}
	} else {
		log.Printf("[INFO] username and password for ftp is not provided\n")
		filter.Connection = &api_ftp.Conn{Host: src.Path}
	}
	files, err := client.Find(context.Background(), &api_ftp.Filter{PathMatch: src.Path})
	if err != nil {
		log.Printf("[ERROR] find return: %s\n", err.Error())
		return nil, err
	}
	if files.State != api_ftp.List_OK {
		log.Printf("[ERROR] listing state: %v, message: %s\n", files.State, files.Error)
		return nil, fmt.Errorf("find: %s", files.Error)
	}
	result := make([]*dslFile, 0)
	for _, file := range files.Files {
		result = append(result, &dslFile{Name: file.Path, Size: file.Size})
	}
	return result, nil
}

func listFilesPosix(src *DSLUrl, conn *grpc.ClientConn) ([]*dslFile, error) {
	client := api_posix.NewFileClient(conn)
	files, err := client.Find(context.Background(), &api_posix.Filter{PathMatch: src.Path})
	if err != nil {
		log.Printf("[ERROR] find return: %s\n", err.Error())
		return nil, err
	}
	if files.State != api_posix.List_OK {
		log.Printf("[ERROR] listing state: %v, message: %s\n", files.State, files.Error)
		return nil, fmt.Errorf("find: %s", files.Error)
	}
	result := make([]*dslFile, 0)
	for _, file := range files.Files {
		result = append(result, &dslFile{Name: file.Path, Size: file.Size})
	}
	return result, nil
}
