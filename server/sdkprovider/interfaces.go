package sdkprovider

import "fabric-sdk-go/protos"

type SdkProvider interface {
	CreateChannel(channelID string) (transactionID string, code protos.StatusCode, err error)
	JoinChannel(channelID string) (code protos.StatusCode, err error)
}
