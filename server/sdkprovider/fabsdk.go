package sdkprovider

import (
	"errors"
	pb "fabric-sdk-go/protos"
	"fabric-sdk-go/server/helpers"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

var logger = helpers.GetLogger()
var appConf = helpers.GetAppConf().Conf

type OrgInstance struct {
	Config      *helpers.OrgInfo
	AdminClient *resmgmt.Client
	MspClient   *mspclient.Client
	Peers       []fab.Peer
}

type OrdererInstance struct {
	Config      *helpers.OrderderInfo
	AdminClient *resmgmt.Client
}

type FabSdkProvider struct {
	Sdk      *fabsdk.FabricSDK
	Orgs     []*OrgInstance
	Orderers []*OrdererInstance
}

func loadOrgPeers(org string, ctxProvider contextAPI.ClientProvider) ([]fab.Peer, error) {
	ctx, err := ctxProvider()
	if err != nil {
		return nil, err
	}

	orgPeers, ok := ctx.EndpointConfig().PeersConfig(org)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Failed to load org peers for %s", org))
	}
	peers := make([]fab.Peer, len(orgPeers))
	for i, val := range orgPeers {
		if peer, err := ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: val}); err != nil {
			return nil, err
		} else {
			peers[i] = peer
		}

	}
	return peers, nil
}

func NewFabSdkProvider() (*FabSdkProvider, error) {
	configOpt := config.FromFile(helpers.GetConfigPath("config.yaml"))
	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		logger.Error("Failed to create new SDK: %s", err)
		return nil, err
	}

	provider := &FabSdkProvider{
		Sdk:      sdk,
		Orgs:     make([]*OrgInstance, len(appConf.OrgInfo)),
		Orderers: make([]*OrdererInstance, len(appConf.OrderderInfo)),
	}
	for i, org := range appConf.OrgInfo {
		//clientContext allows creation of transactions using the supplied identity as the credential.
		adminContext := sdk.Context(fabsdk.WithUser(org.Admin), fabsdk.WithOrg(org.Name))

		mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(org.Name))
		if err != nil {
			logger.Error("Failed to create mspClient for %s, err: %v", org.Name, err)
			return nil, err
		}
		// Resource management client is responsible for managing channels (create/update channel)
		// Supply user that has privileges to create channel (in this case orderer admin)
		adminClient, err := resmgmt.New(adminContext)
		if err != nil {
			logger.Error("Failed to new resource management client: %s", err)
			return nil, err
		}

		orgPeers, err := loadOrgPeers(org.Name, adminContext)
		if err != nil {
			logger.Error("Failed to load peers for %s, err: %v", org.Name, err)
			return nil, err
		}

		provider.Orgs[i] = &OrgInstance{org, adminClient, mspClient, orgPeers}
	}

	if len(provider.Orgs) == 0 {
		logger.Error("Not provider org config in conf/app.yaml", err)
		return nil, errors.New("not provider org config")
	}

	for i, orderer := range appConf.OrderderInfo {
		//clientContext allows creation of transactions using the supplied identity as the credential.
		adminContext := sdk.Context(fabsdk.WithUser(orderer.Admin), fabsdk.WithOrg(orderer.Name))

		// Resource management client is responsible for managing channels (create/update channel)
		// Supply user that has privileges to create channel (in this case orderer admin)
		adminClient, err := resmgmt.New(adminContext)
		if err != nil {
			logger.Error("Failed to new resource management client: %s", err)
			return nil, err
		}
		provider.Orderers[i] = &OrdererInstance{orderer, adminClient}
	}

	return provider, nil
}

func (f *FabSdkProvider) CreateChannel(channelID string) (helpers.TransactionID, pb.StatusCode, error) {
	if len(f.Orderers) == 0 {
		return "", pb.StatusCode_FAILED, errors.New("not found orderers")
	}

	signingIdentities := make([]msp.SigningIdentity, len(f.Orgs))
	var err error
	for i, org := range f.Orgs {
		signingIdentities[i], err = org.MspClient.GetSigningIdentity(org.Config.Admin)
		if err != nil {
			logger.Error("MspClient getSigningIdentity err: %s", err)
			return "", pb.StatusCode_FAILED_GET_SIGNING_IDENTITY, err
		}
	}

	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: helpers.GetChannelConfigPath(channelID + ".tx"),
		SigningIdentities: signingIdentities}

	txID, err := f.Orderers[0].AdminClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(f.Orderers[0].Config.Endpoint))
	if err != nil {
		logger.Error("Failed SaveChannel: %s", err)
		return "", pb.StatusCode_FAILED_CREATE_CHANNEL, err
	}
	logger.Debug("Successfully created channel: %s", channelID)
	return helpers.TransactionID(txID.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) JoinChannel(channelID string) (pb.StatusCode, error) {
	if len(f.Orderers) == 0 {
		return pb.StatusCode_FAILED, errors.New("not found orderers")
	}

	for _, orgInstance := range f.Orgs {
		// Org peers join channel
		if err := orgInstance.AdminClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts),
			resmgmt.WithOrdererEndpoint(f.Orderers[0].Config.Endpoint)); err != nil {
			logger.Error("%s failed to JoinChannel: %v", orgInstance.Config.Name, err)
			return pb.StatusCode_FAILED_JOIN_CHANNEL, err
		}
		logger.Debug("%s joined channel: %s successfully", orgInstance.Config.Name, channelID)
	}

	return pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) InstallCC(ccID, ccVersion, ccPath string) (pb.StatusCode, error) {
	ccPkg, err := packager.NewCCPackage(ccPath, helpers.GetDeployPath())
	if err != nil {
		logger.Error("New cc package err: %s", err)
		return pb.StatusCode_FAILED_NEW_CCPACKAGE, err
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Package: ccPkg}

	for _, orgInstance := range f.Orgs {
		_, err = orgInstance.AdminClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			logger.Error("Failed InstallCC: %s to %s peers", err, orgInstance.Config.Name)
			return pb.StatusCode_FAILED_INSTALL_CC, err
		}
		logger.Debug("Successfully install chaincode [%s:%s] to %s peers", ccID, ccVersion, orgInstance.Config.Name)
	}

	return pb.StatusCode_SUCCESS, err
}

func (f *FabSdkProvider) InstantiateCC(channelID, ccID, ccVersion, ccPath, ccPolicy string, args [][]byte) (helpers.TransactionID, pb.StatusCode, error) {
	policy, err := cauthdsl.FromString(ccPolicy)
	if err != nil {
		logger.Error("Failed parse cc policy[%s], err:%v", ccPolicy, err)
		return "", pb.StatusCode_FAILED_INSTANTIATE_CC, err
	}

	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := f.Orgs[0].AdminClient.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Args: args, Policy: policy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		logger.Error("Failed InstantiateCC: %s", err)
		return "", pb.StatusCode_FAILED_INSTANTIATE_CC, err
	}
	logger.Debug("Successfully instantiate chaincode  [%s:%s]", ccID, ccVersion)
	return helpers.TransactionID(resp.TransactionID), pb.StatusCode_SUCCESS, nil
}

func (f *FabSdkProvider) UpgradeCC(channelID, ccID, ccVersion, ccPath, ccPolicy string, args [][]byte) (helpers.TransactionID, pb.StatusCode, error) {
	policy, err := cauthdsl.FromString(ccPolicy)
	if err != nil {
		logger.Error("Failed parse cc policy[%s], err:%v", ccPolicy, err)
		return "", pb.StatusCode_FAILED_UPGRADE_CC, err
	}

	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := f.Orgs[0].AdminClient.UpgradeCC(
		channelID,
		resmgmt.UpgradeCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Args: args, Policy: policy},
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
	//ledger.WithTargets(orgTestPeer0, orgTestPeer1)
	orgInstance := f.Orgs[0]
	//prepare context
	userContext := f.Sdk.ChannelContext(channelID, fabsdk.WithUser(orgInstance.Config.User), fabsdk.WithOrg(orgInstance.Config.Name))
	//get channel client
	chClient, err := channel.New(userContext)
	if err != nil {
		logger.Error("Failed to create new channel client: %v", err)
		return nil, "", pb.StatusCode_INVALID_USER_CLIENT, fmt.Errorf("Failed to create new channel client:  %s", orgInstance.Config.Name)
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
	orgInstance := f.Orgs[0]

	//prepare context
	userContext := f.Sdk.ChannelContext(channelID, fabsdk.WithUser(orgInstance.Config.User), fabsdk.WithOrg(orgInstance.Config.Name))
	//get channel client
	chClient, err := channel.New(userContext)
	if err != nil {
		logger.Error("Failed to create new channel client: %v", err)
		return nil, pb.StatusCode_INVALID_USER_CLIENT, fmt.Errorf("Failed to create new channel client:  %s", orgInstance.Config.Name)
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
