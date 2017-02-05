package dsl

import (
	"fmt"
	"log"
	"time"

	"github.com/secsy/goftp"
	"github.com/yuin/gopher-lua"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api_posix "api/posix"
	auth "auth"
)

func (d *Config) dslRemove(L *lua.LState) int {

	src, err := parseUrl(L.CheckString(1))
	if err != nil {
		log.Printf("[ERROR] bad url: %s\n", err.Error())
		L.ArgError(1, fmt.Sprintf("can't parse url: %s", err.Error()))
		return -1
	}
	log.Printf("[INFO] start remove: %v\n", src)
	server := src.Host

	switch src.Scheme {

	// удаление с ftp
	case DSLFtpScheme:
		ftpConfig := goftp.Config{Timeout: 5 * time.Second}
		if src.Url.User != nil {
			ftpConfig.User = src.Url.User.Username()
			if passwd, ok := src.Url.User.Password(); ok {
				ftpConfig.Password = passwd
			}
		}
		ftpClient, err := goftp.DialConfig(ftpConfig, server)
		if err != nil {
			log.Printf("[ERROR] dial ftp %s: %s\n", server, err.Error())
			L.RaiseError("dial ftp %s: %s\n", server, err.Error())
			return -1
		}
		if err := ftpClient.Delete(src.Path); err != nil {
			log.Printf("[ERROR] delete from ftp: %s\n", err.Error())
			L.RaiseError("delete from ftp: %s\n", server, err.Error())
			return -1
		}
		return 0

	// удаление с posix
	case DSLPosixScheme:
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
		client := api_posix.NewFileClient(conn)
		info, err := client.LocalOps(
			context.Background(),
			&api_posix.LocalOperation{
				Type: api_posix.LocalOperation_DELETE,
				File: &api_posix.Info{Path: src.Path},
			})
		if err != nil {
			log.Printf("[ERROR] delete operation: %s\n", err.Error())
			L.RaiseError("delete error: %s", err.Error())
			return -1
		}
		if info.State != api_posix.LocalOperation_OK {
			log.Printf("[ERROR] delete operation: %s\n", info.Error)
			L.RaiseError("delete error: %s", info.Error)
			return -1
		}
		return 0

	// здесь нам пришло что-то не понятное
	default:
		panic("program bug")

	}
}
