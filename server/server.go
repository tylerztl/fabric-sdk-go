package server

import (
	pb "fabric-sdk-go/protos"
	"fabric-sdk-go/server/grpchandler"
	"fabric-sdk-go/server/helpers"
	"net"

	"google.golang.org/grpc"
)

var (
	ServerPort string
	EndPoint   string
)

var logger = helpers.GetLogger()

func Run() (err error) {
	EndPoint = ":" + ServerPort
	conn, err := net.Listen("tcp", EndPoint)
	if err != nil {
		logger.Error("TCP Listen err:%s", err)
	}

	srv := newGrpc()
	logger.Info("gRPC and https listen on: %s", ServerPort)

	if err = srv.Serve(conn); err != nil {
		logger.Error("ListenAndServe: %s", err)
	}

	return err
}

func newGrpc() *grpc.Server {
	server := grpc.NewServer()
	// TODO
	pb.RegisterChannelServer(server, grpchandler.NewChannelService())
	pb.RegisterChaincodeServer(server, grpchandler.NewChaincodeService())

	return server
}
