package auth

import (
	"crypto/md5"
	"fmt"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	AuthGrpcAccessDenied  = fmt.Errorf("Auth: access denied")
	AuthGrpcEmptyMetadata = fmt.Errorf("Auth: empty metadata")
	grpcServer            = "localhost"
	grpcSalt              = "Uojee5ah Sae4aili eeWoo8uj aTi3EeF9 ulu1Aiph Neup4uch ohchai2A phuog9Wa"
)

func SetGrpcServerName(name string) {
	grpcServer = name
}

func GetGrpcServerName() string {
	return grpcServer
}

// https://github.com/grpc/grpc-go/issues/106

func NewGrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(StreamInterceptor),
		grpc.UnaryInterceptor(UnaryInterceptor),
	)
}

type GrpcCreds struct {
	Type string
	Code string
}

func GrpcAuthFor(server string) *GrpcCreds {
	return &GrpcCreds{Type: "md5", Code: generateCode(server)}
}

func (c *GrpcCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"code": c.Code, "type": c.Type}, nil
}

func (c *GrpcCreds) RequireTransportSecurity() bool {
	return false
}

func StreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorizeGrpc(stream.Context()); err != nil {
		return err
	}
	return handler(srv, stream)
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorizeGrpc(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func authorizeGrpc(ctx context.Context) error {
	if md, ok := metadata.FromContext(ctx); ok {
		if len(md["type"]) > 0 && md["type"][0] == "md5" {
			if len(md["code"]) > 0 && md["code"][0] == generateCode(grpcServer) {
				return nil
			}
		}
		return AuthGrpcAccessDenied
	}
	return AuthGrpcEmptyMetadata
}

func generateCode(server string) string {
	h := md5.New()
	io.WriteString(h, server)
	io.WriteString(h, grpcSalt)
	return fmt.Sprintf("%x", h.Sum(nil))
}
