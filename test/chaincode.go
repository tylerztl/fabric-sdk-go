package test

import (
	pb "fabric-sdk-go/protos"
	"fmt"

	"golang.org/x/net/context"
)

func InstallCC(ccID, ccVersion, ccPath string) (pb.StatusCode, error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChaincodeClient(conn)
	context := context.Background()
	body := &pb.InstallCCRequest{CcId: ccID, CcVersion: ccVersion, CcPath: ccPath}

	r, err := c.InstallCC(context, body)
	fmt.Printf("StatusCode: %s, err: %v", r.Status, err)
	return r.Status, err
}

func InstantiateCC(channelID, ccID, ccVersion, ccPath string, args [][]byte) (
	code pb.StatusCode, err error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChaincodeClient(conn)
	context := context.Background()
	body := &pb.InstantiateCCRequest{
		ChannelId: channelID,
		CcId:      ccID,
		CcVersion: ccVersion,
		CcPath:    ccPath,
		Args:      args}

	r, err := c.InstantiateCC(context, body)
	fmt.Printf("StatusCode: %s, transaction id: %s, err: %v", r.Status, r.TransactionId, err)
	return r.Status, err
}
