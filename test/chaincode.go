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
	fmt.Printf("StatusCode: %s, err: %v\n", r.Status, err)
	return r.Status, err
}

func InstantiateCC(channelID, ccID, ccVersion, ccPath, ccPolicy string, args [][]byte) (
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
		CcPolicy:  ccPolicy,
		Args:      args}

	r, err := c.InstantiateCC(context, body)
	fmt.Printf("StatusCode: %s, transaction id: %s, err: %v\n", r.Status, r.TransactionId, err)
	return r.Status, err
}

func UpgradeCC(channelID, ccID, ccVersion, ccPath, ccPolicy string, args [][]byte) (
	code pb.StatusCode, err error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChaincodeClient(conn)
	context := context.Background()
	body := &pb.UpgradeCCRequest{
		ChannelId: channelID,
		CcId:      ccID,
		CcVersion: ccVersion,
		CcPath:    ccPath,
		CcPolicy:  ccPolicy,
		Args:      args}

	r, err := c.UpgradeCC(context, body)
	fmt.Printf("StatusCode: %s, transaction id: %s, err: %v\n", r.Status, r.TransactionId, err)
	return r.Status, err
}

func InvokeCC(channelID, ccID, function string, args [][]byte) (
	code pb.StatusCode, err error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChaincodeClient(conn)
	context := context.Background()
	body := &pb.InvokeCCRequest{
		ChannelId: channelID,
		CcId:      ccID,
		Func:      function,
		Args:      args}

	r, err := c.InvokeCC(context, body)
	fmt.Printf("StatusCode: %s, transaction id: %s, err: %v\n", r.Status, r.TransactionId, err)
	return r.Status, err
}

func QueryCC(channelID, ccID, function string, args [][]byte) (
	code pb.StatusCode, err error) {
	conn := NewConn()
	defer conn.Close()

	c := pb.NewChaincodeClient(conn)
	context := context.Background()
	body := &pb.QueryCCRequest{
		ChannelId: channelID,
		CcId:      ccID,
		Func:      function,
		Args:      args}

	r, err := c.QueryCC(context, body)
	fmt.Printf("StatusCode: %s, payload: %s, err: %v\n", r.Status, r.Payload, err)
	return r.Status, err
}
