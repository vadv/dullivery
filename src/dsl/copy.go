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

func (d *Config) dslCopy(L *lua.LState) int {

	src, err := parseUrl(L.CheckString(1))
	if err != nil {
		log.Printf("[ERROR] source url: %s\n", err.Error())
		L.ArgError(1, fmt.Sprintf("can't parse source url: %s", err.Error()))
		return -1
	}

	dst, err := parseUrl(L.CheckString(2))
	if err != nil {
		log.Printf("[ERROR] dst url: %s\n", err.Error())
		L.ArgError(2, fmt.Sprintf("can't parse dst url: %s", err.Error()))
		return -1
	}

	if src.Scheme == DSLFtpScheme && dst.Scheme == DSLFtpScheme {
		log.Printf("[ERROR] copy from ftp to ftp is not supported\n")
		L.RaiseError("copy from ftp to ftp is not supported")
		return -1
	}

	server := src.Host
	address := src.PosixAddress()
	if src.Scheme == DSLFtpScheme {
		// если нужно скопировать с ftp куда-то на другой сервер, то
		// выбираем указанный сервер для того чтобы он мог аплоадить
		server = dst.Host
		address = dst.PosixAddress()
	}

	// make connection
	conn, err := grpc.Dial(
		address,
		grpc.WithPerRPCCredentials(auth.GrpcAuthFor(server)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[ERROR] dial to %s: %s\n", server, err.Error())
		L.RaiseError("dial %s: %s", server, err.Error())
		return -1
	}
	defer conn.Close()

	switch {
	case src.Scheme == DSLFtpScheme && dst.Scheme == DSLPosixScheme:
		if err := downloadFromFtp(conn, src, dst); err != nil {
			L.RaiseError("download from ftp: %s", err.Error())
			return -1
		}

	case src.Scheme == DSLPosixScheme && dst.Scheme == DSLFtpScheme:
		if err := uploadToFtp(conn, src, dst); err != nil {
			L.RaiseError("upload to ftp: %s", err.Error())
			return -1
		}

	case src.Scheme == DSLPosixScheme && dst.Scheme == DSLPosixScheme:
		if err := copyPosix(src, dst); err != nil {
			L.RaiseError("posix copy: %s", err.Error())
			return -1
		}
	default:
		L.RaiseError("unsupported operation for copy('%s', '%s')", src.Scheme, dst.Scheme)
	}

	return 0
}

func copyPosix(src, dst *DSLUrl) error {

	log.Printf("[INFO] Run copy from `%s` to `%s`\n", src.ToString(), dst.ToString())

	// нужно подключиться к dst и попросить его скопировать файл с src
	conn, err := grpc.Dial(
		dst.PosixAddress(),
		grpc.WithPerRPCCredentials(auth.GrpcAuthFor(dst.Host)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Printf("[ERROR] dial to dst %s: %s\n", dst.PosixAddress(), err.Error())
		return err
	}
	defer conn.Close()
	client := api_posix.NewFileClient(conn)

	// локальное перемещение
	if src.PosixAddress() == dst.PosixAddress() {
		if src.Path == dst.Path {
			return fmt.Errorf("copy itself")
		}
		info, err := client.LocalOps(
			context.Background(),
			&api_posix.LocalOperation{
				Type:    api_posix.LocalOperation_COPY,
				File:    &api_posix.Info{Path: src.Path},
				DstFile: &api_posix.Info{Path: dst.Path},
			},
		)
		if err != nil {
			log.Printf("[ERROR] local copy on %s: %s\n", src.Host, err.Error())
			return err
		}
		if info.State != api_posix.LocalOperation_OK {
			log.Printf("[ERROR] local operations on server %s: %s\n", src.Host, info.Error)
			return fmt.Errorf("local operations on server %s: %s", src.Host, info.Error)
		}
		return nil
	}

	// удаленному серверу dst даем указание скопировать с src
	info, err := client.RemoteOps(
		context.Background(),
		&api_posix.RemoteOperation{
			Type:         api_posix.RemoteOperation_COPY_FROM,
			ToFile:       &api_posix.Info{Path: dst.Path}, // source file
			RemoteFile:   &api_posix.Info{Path: src.Path}, // local path for dst
			RemoteServer: src.Host,                        // source host
			RemotePort:   int64(src.Port),                 // source port
		})
	if err != nil {
		log.Printf("[ERROR] remote copy on `%s`: %s\n", src.Host, err.Error())
		return err
	}
	if info.State != api_posix.RemoteOperation_OK {
		log.Printf("[ERROR] remote operations on server %s: %s\n", src.Host, info.Error)
		return fmt.Errorf("remote operations on server %s: %s", src.Host, info.Error)
	}
	return nil
}

func downloadFromFtp(conn *grpc.ClientConn, ftp, dst *DSLUrl) error {
	client := api_ftp.NewFileClient(conn)
	// ищем указанный файл на ftp
	info := &api_ftp.Info{Path: ftp.Path, Connection: &api_ftp.Conn{Host: ftp.Host}}
	if ftp.Url.User != nil {
		info.Connection.User = ftp.Url.User.Username()
		if passwd, ok := ftp.Url.User.Password(); ok {
			info.Connection.Password = passwd
		}
	}
	info.LocalPath = dst.Path
	result, err := client.Download(context.Background(), info)
	if err != nil {
		log.Printf("[ERROR] while download: %s\n", err.Error())
		return err
	}
	if result.State != api_ftp.Info_OK {
		log.Printf("[ERROR] info state: %v, message: %s\n", result.State, result.Error)
		return fmt.Errorf("ftp download: %s", result.Error)
	}
	return nil
}

func uploadToFtp(conn *grpc.ClientConn, src, ftp *DSLUrl) error {
	client := api_ftp.NewFileClient(conn)
	// ищем указанный файл на ftp
	info := &api_ftp.Info{Path: ftp.Path, Connection: &api_ftp.Conn{Host: ftp.Host}}
	if ftp.Url.User != nil {
		info.Connection.User = ftp.Url.User.Username()
		if passwd, ok := ftp.Url.User.Password(); ok {
			info.Connection.Password = passwd
		}
	}
	info.LocalPath = src.Path
	result, err := client.Upload(context.Background(), info)
	if err != nil {
		log.Printf("[ERROR] while download: %s\n", err.Error())
		return err
	}
	if result.State != api_ftp.Info_OK {
		log.Printf("[ERROR] info state: %v, message: %s\n", result.State, result.Error)
		return fmt.Errorf("ftp upload: %s", result.Error)
	}
	return nil
}
