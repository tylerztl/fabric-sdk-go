package grpchandler

import "fabric-sdk-go/server/sdkprovider"

type Handler struct {
	Provider sdkprovider.SdkProvider
}

var hanlder = NewHandler()

func init() {
	provider, err := sdkprovider.NewFabSdkProvider()
	if err != nil {
		panic(err)
	}
	hanlder.Provider = provider
}

func NewHandler() *Handler {
	return &Handler{}
}

func GetSdkProvider() sdkprovider.SdkProvider {
	return hanlder.Provider
}
