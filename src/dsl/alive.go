package dsl

import (
	"fmt"
	"log"
	"net"

	"github.com/yuin/gopher-lua"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api_health "api/health"
	auth "auth"
)

func (d *Config) dslAlive(L *lua.LState) int {

	address := L.CheckString(1)
	server := address
	if host, _, err := net.SplitHostPort(address); err != nil {
		address = fmt.Sprintf("%s:%d", address, DSLPosixDefaultPort)
	} else {
		// попарсился host:port
		server = host
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithPerRPCCredentials(auth.GrpcAuthFor(server)),
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Printf("[INFO] dial %s: %s\n", address, err.Error())
		L.Push(lua.LBool(false))
		return 1
	}
	defer conn.Close()

	client := api_health.NewHealthClient(conn)
	_, err = client.Ping(context.Background(), &api_health.PingMsg{})

	if err != nil {
		log.Printf("[INFO] ping %s: %s\n", address, err.Error())
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(true))
	return 1

}
