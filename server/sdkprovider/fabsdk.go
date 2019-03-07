package sdkprovider

import (
	"fabric-sdk-go/server/helpers"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var logger = helpers.GetLogger()
var appConf = helpers.GetAppConf()

type FabSdkProvider struct {
	Sdk *fabsdk.FabricSDK
}

func NewFabSdkProvider() (*FabSdkProvider, error) {
	configOpt := config.FromFile(helpers.GetConfigPath("config.yaml"))
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		logger.Error("Failed to create new SDK: %s", err)
		return nil, err
	}

	return &FabSdkProvider{Sdk: sdk}, nil
}

func (f *FabSdkProvider) CreateChannel(channelID string) (string, error) {
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := f.Sdk.Context(fabsdk.WithUser(appConf.Conf.OrgAdmin), fabsdk.WithOrg(appConf.Conf.OrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Error("Failed to create channel management client: %s", err)
		return "", err
	}
	mspClient, err := mspclient.New(f.Sdk.Context(), mspclient.WithOrg(appConf.Conf.OrgName))
	if err != nil {
		logger.Error("New mspclient err: %s", err)
		return "", err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(appConf.Conf.OrgAdmin)
	if err != nil {
		logger.Error("MspClient getSigningIdentity err: %s", err)
		return "", err
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: helpers.GetChannelConfigPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(helpers.OrdererEndpoint))
	if err != nil {
		logger.Error("Failed SaveChannel: %s", err)
		return "", err
	}
	return string(txID.TransactionID), nil
}
