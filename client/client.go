package main

import (
	pb "fabric-sdk-go/protos"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Println("1", err)
	}

	c := pb.NewChannelClient(conn)
	context := context.Background()
	body := &pb.CreateChannelRequest{ChannelId: "mychannel"}

	r, err := c.CreateChannel(context, body)
	if err != nil {
		log.Println("2", err)
	}

	log.Println(r.Status, r.TransactionId)
}
