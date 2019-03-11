package test

import (
	pb "fabric-sdk-go/protos"
	"fmt"

	"golang.org/x/net/context"
)

func CreateChannel(channelId string) (pb.StatusCode, error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChannelClient(conn)
	context := context.Background()
	body := &pb.CreateChannelRequest{ChannelId: channelId}

	r, err := c.CreateChannel(context, body)
	fmt.Printf("StatusCode: %s, transaction id: %s, err: %v\n", r.Status, r.TransactionId, err)
	return r.Status, err
}

func JoinChannel(channelId string) (pb.StatusCode, error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChannelClient(conn)
	context := context.Background()
	body := &pb.JoinChannelRequest{ChannelId: channelId}

	r, err := c.JoinChannel(context, body)
	fmt.Printf("StatusCode: %s, err: %v\n", r.Status, err)
	return r.Status, err
}
