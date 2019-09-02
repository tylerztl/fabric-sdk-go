package sdkprovider

import (
	pb "fabric-sdk-go/protos"
	"fabric-sdk-go/server/helpers"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

var logger = helpers.GetLogger()
var appConf = helpers.GetAppConf().Conf

type OrgInstance struct {
	Admin       string
	User        string
	AdminClient *resmgmt.Client
}

type FabSdkProvider struct {
	Sdk        *fabsdk.FabricSDK
	Org        map[string]*OrgInstance
	DefaultOrg string
}

func NewFabSdkProvider() (*FabSdkProvider, error) {
	configOpt := config.FromFile(helpers.GetConfigPath("config.yaml"))
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		logger.Error("Failed to create new SDK: %s", err)
		return nil, err
	}

	provider := &FabSdkProvider{
		Sdk: sdk,
		Org: make(map[string]*OrgInstance),
	}
	for _, org := range appConf.OrgInfo {
		//clientContext allows creation of transactions using the supplied identity as the credential.
		adminContext := sdk.Context(fabsdk.WithUser(org.Admin), fabsdk.WithOrg(org.Name))

		// Resource management client is responsible for managing channels (create/update channel)
		// Supply user that has privileges to create channel (in this case orderer admin)
		adminClient, err := resmgmt.New(adminContext)
		if err != nil {
			logger.Error("Failed to new resource management client: %s", err)
			return nil, err
		}
		provider.Org[org.Name] = &OrgInstance{org.Admin, org.User, adminClient}
		if org.Default {
			provider.DefaultOrg = org.Name
		}
	}

	return provider, nil
}

func (f *FabSdkProvider) CreateChannel(channelID string) (helpers.TransactionID, pb.StatusCode, error) {
	orgName := f.DefaultOrg
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return "", pb.StatusCode_INVALID_ADMIN_CLIENT, fmt.Errorf("Not found admin client for org:  %v", orgName)
	}
	mspClient, err := mspclient.New(f.Sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		logger.Error("New mspclient err: %s", err)
		return "", pb.StatusCode_FAILED_NEW_MSP_CLIENT, err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgInstance.Admin)
	if err != nil {
		logger.Error("MspClient getSigningIdentity err: %s", err)
		return "", pb.StatusCode_FAILED_GET_SIGNING_IDENTITY, err
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: helpers.GetChannelConfigPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := orgInstance.AdminClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(appConf.OrdererEndpoint))
	if err != nil {
		logger.Error("Failed SaveChannel: %s", err)
		return "", pb.StatusCode_FAILED_CREATE_CHANNEL, err
	}
	logger.Debug("Successfully created channel: %s", channelID)
	return helpers.TransactionID(txID.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) JoinChannel(channelID string) (pb.StatusCode, error) {
	orgName := f.DefaultOrg
	// Org resource management client
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return pb.StatusCode_INVALID_ADMIN_CLIENT, fmt.Errorf("Not found admin client for org:  %v", orgName)
	}

	// Org peers join channel
	err := orgInstance.AdminClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(appConf.OrdererEndpoint))
	if err != nil {
		logger.Error("Org peers failed to JoinChannel: %v", err)
		return pb.StatusCode_FAILED_JOIN_CHANNEL, err
	}
	logger.Debug("Successfully joined channel: %s", channelID)
	return pb.StatusCode_SUCCESS, err
}

func (f *FabSdkProvider) InstallCC(ccID, ccVersion, ccPath string) (pb.StatusCode, error) {
	orgName := f.DefaultOrg
	// Org resource management client
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return pb.StatusCode_INVALID_ADMIN_CLIENT, fmt.Errorf("Not found admin client for org:  %v", orgName)
	}

	ccPkg, err := packager.NewCCPackage(ccPath, helpers.GetDeployPath())
	if err != nil {
		logger.Error("New cc package err: %s", err)
		return pb.StatusCode_FAILED_NEW_CCPACKAGE, err
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Package: ccPkg}
	_, err = orgInstance.AdminClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed InstallCC: %s", err)
		return pb.StatusCode_FAILED_INSTALL_CC, err
	}
	logger.Debug("Successfully install chaincode [%s:%s]", ccID, ccVersion)
	return pb.StatusCode_SUCCESS, err
}

func (f *FabSdkProvider) InstantiateCC(channelID, ccID, ccVersion, ccPath string, args [][]byte) (helpers.TransactionID, pb.StatusCode, error) {
	orgName := f.DefaultOrg
	// Org resource management client
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return "", pb.StatusCode_INVALID_ADMIN_CLIENT, fmt.Errorf("Not found admin client for org:  %v", orgName)
	}
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := orgInstance.AdminClient.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Args: args, Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		logger.Error("Failed InstantiateCC: %s", err)
		return "", pb.StatusCode_FAILED_INSTANTIATE_CC, err
	}
	logger.Debug("Successfully instantiate chaincode  [%s:%s]", ccID, ccVersion)
	return helpers.TransactionID(resp.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) UpgradeCC(channelID, ccID, ccVersion, ccPath string, args [][]byte) (helpers.TransactionID, pb.StatusCode, error) {
	orgName := f.DefaultOrg
	// Org resource management client
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return "", pb.StatusCode_INVALID_ADMIN_CLIENT, fmt.Errorf("Not found admin client for org:  %v", orgName)
	}
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := orgInstance.AdminClient.UpgradeCC(
		channelID,
		resmgmt.UpgradeCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Args: args, Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		logger.Error("Failed UpgradeCC: %s", err)
		return "", pb.StatusCode_FAILED_UPGRADE_CC, err
	}
	logger.Debug("Successfully upgrade chaincode  [%s:%s]", ccID, ccVersion)
	return helpers.TransactionID(resp.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) InvokeCC(channelID, ccID, function string, args [][]byte) ([]byte, helpers.TransactionID, pb.StatusCode, error) {
	orgName := f.DefaultOrg
	// Org resource management client
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return nil, "", pb.StatusCode_NOT_FOUND_ORG_INSTANCE, fmt.Errorf("Not found org instance for org:  %v", orgName)
	}
	//prepare context
	userContext := f.Sdk.ChannelContext(channelID, fabsdk.WithUser(orgInstance.User), fabsdk.WithOrg(orgName))
	//get channel client
	chClient, err := channel.New(userContext)
	if err != nil {
		logger.Error("Failed to create new channel client: %v", err)
		return nil, "", pb.StatusCode_INVALID_USER_CLIENT, fmt.Errorf("Failed to create new channel client:  %s", orgName)
	}
	// Synchronous transaction
	response, err := chClient.Execute(
		channel.Request{
			ChaincodeID: ccID,
			Fcn:         function,
			Args:        args,
		},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logger.Error("Failed InvokeCC: %s", err)
		return nil, "", pb.StatusCode_FAILED_INVOKE_CC, err
	}
	logger.Debug("Successfully invoke chaincode  ccName[%s] func[%v] txId[%v] payload[%v]",
		ccID, function, response.TransactionID, response.Payload)
	return response.Payload, helpers.TransactionID(response.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) QueryCC(channelID, ccID, function string, args [][]byte) ([]byte, pb.StatusCode, error) {
	orgName := f.DefaultOrg
	// Org resource management client
	orgInstance, ok := f.Org[orgName]
	if !ok {
		logger.Error("Not found resource management client for org: %s", orgName)
		return nil, pb.StatusCode_NOT_FOUND_ORG_INSTANCE, fmt.Errorf("Not found  org instance for org:  %v", orgName)
	}
	//prepare context
	userContext := f.Sdk.ChannelContext(channelID, fabsdk.WithUser(orgInstance.User), fabsdk.WithOrg(orgName))
	//get channel client
	chClient, err := channel.New(userContext)
	if err != nil {
		logger.Error("Failed to create new channel client: %v", err)
		return nil, pb.StatusCode_INVALID_USER_CLIENT, fmt.Errorf("Failed to create new channel client:  %s", orgName)
	}

	response, err := chClient.Query(channel.Request{ChaincodeID: ccID, Fcn: function, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logger.Error("Failed QueryCC: %s", err)
		return nil, pb.StatusCode_FAILED_QUERY_CC, err
	}

	logger.Debug("Successfully query chaincode  ccName[%s] func[%v] payload[%v]",
		ccID, function, response.Payload)
	return response.Payload, pb.StatusCode_SUCCESS, nil
}
