package sdkprovider

type SdkProvider interface {
	CreateChannel(channelID string) (transactionID string, err error)
}
