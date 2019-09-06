package sdkprovider

import (
	pb "fabric-sdk-go/protos"
	. "fabric-sdk-go/server/helpers"
)

type SdkProvider interface {
	CreateChannel(channelID string) (TransactionID, pb.StatusCode, error)
	JoinChannel(channelID string) (pb.StatusCode, error)
	InstallCC(ccID, ccVersion, ccPath string) (pb.StatusCode, error)
	InstantiateCC(channelID, ccID, ccVersion, ccPath, ccPolicy string, args [][]byte) (TransactionID, pb.StatusCode, error)
	UpgradeCC(channelID, ccID, ccVersion, ccPath, ccPolicy string, args [][]byte) (TransactionID, pb.StatusCode, error)
	InvokeCC(channelID, ccID, function string, args [][]byte) ([]byte, TransactionID, pb.StatusCode, error)
	QueryCC(channelID, ccID, function string, args [][]byte) ([]byte, pb.StatusCode, error)
}
