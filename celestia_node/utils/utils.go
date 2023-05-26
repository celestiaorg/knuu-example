package utils

import (
	"context"
	"fmt"
	"github.com/celestiaorg/celestia-node/api/rpc/client"
	app_utils "github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

func initInstance(instanceName string, nodeType string, chainId string, genesisHash string) (*knuu.Instance, error) {
	instance, err := knuu.NewInstance(instanceName)
	if err != nil {
		return nil, fmt.Errorf("Error creating instance '%v':", err)
	}
	err = instance.SetImage("ghcr.io/celestiaorg/celestia-node:v0.10.0")
	if err != nil {
		return nil, fmt.Errorf("Error setting image: %v", err)
	}
	err = instance.AddPortTCP(2121)
	if err != nil {
		return nil, fmt.Errorf("Error adding port: %v", err)
	}
	err = instance.AddPortTCP(26658)
	if err != nil {
		return nil, fmt.Errorf("Error adding port: %v", err)
	}
	err = instance.SetEnvironmentVariable("CELESTIA_CUSTOM", fmt.Sprintf("%s:%s", chainId, genesisHash))
	if err != nil {
		return nil, fmt.Errorf("Error setting environment variable: %v", err)
	}
	_, err = instance.ExecuteCommand("celestia", nodeType, "init", "--node.store", "/home/celestia/.celestia-test")
	if err != nil {
		return nil, fmt.Errorf("Error executing command: %v", err)
	}
	err = instance.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error committing instance: %v", err)
	}
	return instance, nil
}

func CreateAndStartBridge(executor *knuu.Executor, instanceName string, consensus *knuu.Instance) (*knuu.Instance, error) {
	chainId, err := app_utils.ChainId(executor, consensus)
	if err != nil {
		return nil, fmt.Errorf("error getting chain ID: %w", err)
	}
	genesisHash, err := app_utils.GenesisHash(executor, consensus)
	if err != nil {
		return nil, fmt.Errorf("error getting genesis hash: %w", err)
	}
	consensusIP, err := consensus.GetIP()
	if err != nil {
	}

	bridge, err := initInstance(instanceName, "bridge", chainId, genesisHash)
	if err != nil {
		return nil, fmt.Errorf("error creating instance: %w", err)
	}

	err = bridge.SetCommand("celestia", "bridge", "start", "--node.store", "/home/celestia/.celestia-test", "--core.ip", consensusIP)
	if err != nil {
		return nil, fmt.Errorf("error setting command: %w", err)
	}

	bridge.Start()
	bridge.WaitInstanceIsRunning()

	return bridge, nil
}

func CreateAndStartNode(executor *knuu.Executor, instanceName string, nodeType string, consensus *knuu.Instance, trustedNode *knuu.Instance) (*knuu.Instance, error) {
	chainId, err := app_utils.ChainId(executor, consensus)
	if err != nil {
		return nil, fmt.Errorf("error getting chain ID: %w", err)
	}
	genesisHash, err := app_utils.GenesisHash(executor, consensus)
	if err != nil {
		return nil, fmt.Errorf("error getting genesis hash: %w", err)
	}

	node, err := initInstance(instanceName, nodeType, chainId, genesisHash)
	if err != nil {
		return nil, fmt.Errorf("error creating instance: %w", err)
	}

	adminAuthNode, err := trustedNode.ExecuteCommand("celestia", "full", "auth", "admin", "--node.store", "/home/celestia/.celestia-test")
	if err != nil {
		logrus.Fatalf("Error getting admin auth token: %v", err)
	}
	adminAuthNodeToken, err := authTokenFromAuth(adminAuthNode)

	p2pInfoNode, err := trustedNode.ExecuteCommand("celestia", "rpc", "p2p", "Info", "--auth", adminAuthNodeToken)
	if err != nil {
		logrus.Fatalf("Error getting p2p info: %v", err)
	}

	bridgeIP, err := trustedNode.GetIP()
	if err != nil {
		return nil, fmt.Errorf("error getting IP: %w", err)
	}
	bridgeID, err := iDFromP2PInfo(p2pInfoNode)
	if err != nil {
		return nil, fmt.Errorf("error getting ID: %w", err)
	}
	trustedPeers := fmt.Sprintf("/ip4/%s/tcp/2121/p2p/%s", bridgeIP, bridgeID)

	err = node.SetCommand("celestia", nodeType, "start", "--node.store", "/home/celestia/.celestia-test", "--headers.trusted-peers", trustedPeers)
	if err != nil {
		return nil, fmt.Errorf("error setting command: %w", err)
	}

	node.Start()
	node.WaitInstanceIsRunning()

	return node, nil
}

func GetRPCClient(node *knuu.Instance) (*client.Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	localPort, err := node.PortForwardTCP(26658)
	if err != nil {
		return nil, fmt.Errorf("error port forwarding: %w", err)
	}
	url := fmt.Sprintf("%s:%d", "0.0.0.0", localPort)
	adminAuthNode, err := node.ExecuteCommand("celestia", "full", "auth", "admin", "--node.store", "/home/celestia/.celestia-test")
	if err != nil {
		logrus.Fatalf("Error getting admin auth token: %v", err)
	}
	adminAuthNodeToken, err := authTokenFromAuth(adminAuthNode)
	rpcClient, err := client.NewClient(ctx, "http://"+url, adminAuthNodeToken)
	if err != nil {
		return nil, fmt.Errorf("error creating RPC client: %w", err)
	}
	return rpcClient, nil
}
