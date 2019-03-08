package services

import (
	"fabric-sdk-go/server/sdkprovider"
	"golang.org/x/net/context"

	pb "fabric-sdk-go/protos"
)

type ChaincdoeService struct {
	provider sdkprovider.SdkProvider
}

func NewChaincodeService() *ChaincdoeService {
	return &ChaincdoeService{
		provider: GetSdkProvider(),
	}
}

func (c *ChaincdoeService) InstallCC(ctx context.Context, r *pb.InstallCCRequest) (*pb.ServerStatus, error) {
	code, err := c.provider.InstallCC(r.CcId, r.CcVersion, r.CcPath)
	return &pb.ServerStatus{Status: code}, err
}

func (c *ChaincdoeService) InstantiateCC(ctx context.Context, r *pb.InstantiateCCRequest) (*pb.InstantiateCCResponse, error) {
	transactionID, code, err := c.provider.InstantiateCC(r.ChannelId, r.CcId, r.CcVersion, r.CcPath, r.Args)
	return &pb.InstantiateCCResponse{Status: code, TransactionId: transactionID}, err
}
