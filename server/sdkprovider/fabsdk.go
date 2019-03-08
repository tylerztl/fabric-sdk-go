package sdkprovider

import (
	pb "fabric-sdk-go/protos"
	"fabric-sdk-go/server/helpers"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var logger = helpers.GetLogger()
var appConf = helpers.GetAppConf().Conf

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

func (f *FabSdkProvider) CreateChannel(channelID string) (string, pb.StatusCode, error) {
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := f.Sdk.Context(fabsdk.WithUser(appConf.OrgAdmin), fabsdk.WithOrg(appConf.OrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Error("Failed to create channel management client: %s", err)
		return "", pb.StatusCode_FAILED_NEW_CLIENT, err
	}
	mspClient, err := mspclient.New(f.Sdk.Context(), mspclient.WithOrg(appConf.OrgName))
	if err != nil {
		logger.Error("New mspclient err: %s", err)
		return "", pb.StatusCode_FAILED_NEW_MSP_CLIENT, err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(appConf.OrgAdmin)
	if err != nil {
		logger.Error("MspClient getSigningIdentity err: %s", err)
		return "", pb.StatusCode_FAILED_GET_SIGNING_IDENTITY, err
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: helpers.GetChannelConfigPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(appConf.OrdererEndpoint))
	if err != nil {
		logger.Error("Failed SaveChannel: %s", err)
		return "", pb.StatusCode_FAILED_CREATE_CHANNEL, err
	}
	logger.Debug("Successfully created channel: %s", channelID)
	return string(txID.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) JoinChannel(channelID string) (pb.StatusCode, error) {
	//prepare context
	adminContext := f.Sdk.Context(fabsdk.WithUser(appConf.OrgAdmin), fabsdk.WithOrg(appConf.OrgName))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		logger.Error("Failed to create channel management client: %s", err)
		return pb.StatusCode_FAILED_NEW_CLIENT, err
	}

	// Org peers join channel
	err = orgResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(appConf.OrdererEndpoint))
	if err != nil {
		logger.Error("Org peers failed to JoinChannel: %v", err)
		return pb.StatusCode_FAILED_JOIN_CHANNEL, err
	}
	logger.Debug("Successfully joined channel: %s", channelID)
	return pb.StatusCode_SUCCESS, err
}

func (f *FabSdkProvider) InstallCC(ccID, ccVersion, ccPath string) (pb.StatusCode, error) {
	//prepare context
	adminContext := f.Sdk.Context(fabsdk.WithUser(appConf.OrgAdmin), fabsdk.WithOrg(appConf.OrgName))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		logger.Error("Failed to create channel management client: %s", err)
		return pb.StatusCode_FAILED_NEW_CLIENT, err
	}

	ccPkg, err := packager.NewCCPackage(ccPath, helpers.GetDeployPath())
	if err != nil {
		logger.Error("New cc package err: %s", err)
		return pb.StatusCode_FAILED_NEW_CCPACKAGE, err
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed InstallCC: %s", err)
		return pb.StatusCode_FAILED_INSTALL_CC, err
	}
	logger.Debug("Successfully install cc: %s-%s", ccID, ccVersion)
	return pb.StatusCode_SUCCESS, err
}
