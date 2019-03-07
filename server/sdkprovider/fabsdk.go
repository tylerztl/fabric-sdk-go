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

type FabSdkProvider struct {
	AppConf *helpers.AppConf
	Sdk     *fabsdk.FabricSDK
}

func NewFabSdkProvider() (*FabSdkProvider, error) {
	appConf, err := helpers.LoadAppConf()
	if err != nil {
		logger.Errorf("Failed to load appConf: %s", err)
		return nil, err
	}
	configOpt := config.FromFile(helpers.GetConfigPath("config.yaml"))
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		logger.Errorf("Failed to create new SDK: %s", err)
		return nil, err
	}

	return &FabSdkProvider{AppConf: appConf, Sdk: sdk}, nil
}

func (f *FabSdkProvider) CreateChannel(channelID string) (string, error) {
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := f.Sdk.Context(fabsdk.WithUser(f.AppConf.Conf.OrgAdmin), fabsdk.WithOrg(f.AppConf.Conf.OrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Errorf("Failed to create channel management client: %s", err)
		return "", err
	}
	mspClient, err := mspclient.New(f.Sdk.Context(), mspclient.WithOrg(f.AppConf.Conf.OrgName))
	if err != nil {
		logger.Error(err)
		return "", err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(f.AppConf.Conf.OrgAdmin)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: helpers.GetChannelConfigPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(helpers.OrdererEndpoint))
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return string(txID.TransactionID), nil
}
