package grpchandler

import (
	"fabric-sdk-go/server/sdkprovider"
	"golang.org/x/net/context"

	pb "fabric-sdk-go/protos"
)

type ChannelService struct {
	provider sdkprovider.SdkProvider
}

func NewChannelService() *ChannelService {
	return &ChannelService{
		provider: GetSdkProvider(),
	}
}

func (c *ChannelService) CreateChannel(ctx context.Context, r *pb.CreateChannelRequest) (*pb.CreateChannelResponse, error) {
	transactionID, code, err := c.provider.CreateChannel(r.ChannelId)
	return &pb.CreateChannelResponse{Status: code, TransactionId: string(transactionID)}, err
}

func (c *ChannelService) JoinChannel(ctx context.Context, r *pb.JoinChannelRequest) (*pb.ServerStatus, error) {
	code, err := c.provider.JoinChannel(r.ChannelId)
	return &pb.ServerStatus{Status: code}, err
}
