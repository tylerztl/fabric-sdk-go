package grpchandler

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
	transactionID, code, err := c.provider.InstantiateCC(r.ChannelId, r.CcId, r.CcVersion, r.CcPath, r.CcPolicy, r.Args)
	return &pb.InstantiateCCResponse{Status: code, TransactionId: string(transactionID)}, err
}

func (c *ChaincdoeService) UpgradeCC(ctx context.Context, r *pb.UpgradeCCRequest) (*pb.UpgradeCCResponse, error) {
	transactionID, code, err := c.provider.UpgradeCC(r.ChannelId, r.CcId, r.CcVersion, r.CcPath, r.CcPolicy, r.Args)
	return &pb.UpgradeCCResponse{Status: code, TransactionId: string(transactionID)}, err
}

func (c *ChaincdoeService) InvokeCC(ctx context.Context, r *pb.InvokeCCRequest) (*pb.InvokeCCResponse, error) {
	payload, transactionID, code, err := c.provider.InvokeCC(r.ChannelId, r.CcId, r.Func, r.Args)
	return &pb.InvokeCCResponse{Status: code, TransactionId: string(transactionID), Payload: payload}, err
}

func (c *ChaincdoeService) QueryCC(ctx context.Context, r *pb.QueryCCRequest) (*pb.QueryCCResponse, error) {
	payload, code, err := c.provider.QueryCC(r.ChannelId, r.CcId, r.Func, r.Args)
	return &pb.QueryCCResponse{Status: code, Payload: payload}, err
}
