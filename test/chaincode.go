package test

import (
	pb "fabric-sdk-go/protos"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func InstallCC(ccID, ccVersion, ccPath string) (pb.StatusCode, error) {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return pb.StatusCode_FAILED, err
	}

	c := pb.NewChaincodeClient(conn)
	context := context.Background()
	body := &pb.InstallCCRequest{CcId: ccID, CcVersion: ccVersion, CcPath: ccPath}

	r, err := c.InstallCC(context, body)
	fmt.Printf("StatusCode: %s, err: %v", r.Status, err)
	return r.Status, err
}
