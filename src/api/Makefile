all:
	cd posix && rm -f *.pb.go && protoc --go_out=plugins=grpc:. *.proto
	cd ftp && rm -f *.pb.go && protoc --go_out=plugins=grpc:. *.proto
	cd health && rm -f *.pb.go && protoc --go_out=plugins=grpc:. *.proto

.DEFAULT_GOAL: all
