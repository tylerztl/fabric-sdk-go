package services

import (
	"fabric-sdk-go/server/sdkprovider"
)

type Handler struct {
	Provider sdkprovider.SdkProvider
}

var hanlder *Handler

func Init() {
	provider, err := sdkprovider.NewFabSdkProvider()
	if err != nil {
		panic(err)
	}
	hanlder = &Handler{Provider: provider}
}

func GetSdkProvider() sdkprovider.SdkProvider {
	return hanlder.Provider
}
