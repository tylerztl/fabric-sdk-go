package test

import "google.golang.org/grpc"

const (
	ServerAddress string = "localhost:8080"
)

func NewConn() *grpc.ClientConn {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}
