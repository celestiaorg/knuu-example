package utils

import (
	"encoding/json"
	"fmt"
	app_utils "github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
	"time"
)

func ChainIdFromStatus(status string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(status), &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling status: %w", err)
	}
	resultData := result["result"].(map[string]interface{})
	nodeInfo := resultData["node_info"].(map[string]interface{})
	chainId := nodeInfo["network"].(string)
	return chainId, nil
}

func HashFromBlock(block string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(block), &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling block: %w", err)
	}
	resultData := result["result"].(map[string]interface{})
	blockId := resultData["block_id"].(map[string]interface{})
	blockHash := blockId["hash"].(string)
	return blockHash, nil
}

func CreateAndStartConsensus() (*knuu.Instance, error) {
	consensus, err := knuu.NewInstance("consensus")
	if err != nil {
		return nil, fmt.Errorf("error creating instance: %w", err)
	}
	err = consensus.SetImage("ghcr.io/celestiaorg/celestia-app:v0.13.2")
	if err != nil {
		return nil, fmt.Errorf("error setting image: %w", err)
	}
	err = consensus.AddPortTCP(26656)
	if err != nil {
		return nil, fmt.Errorf("error adding port: %w", err)
	}
	err = consensus.AddPortTCP(26657)
	if err != nil {
		return nil, fmt.Errorf("error adding port: %w", err)
	}
	err = consensus.AddFile("resources/genesis.sh", "/opt/genesis.sh", "0:0")
	if err != nil {
		return nil, fmt.Errorf("error adding file: %w", err)
	}
	_, err = consensus.ExecuteCommand("/bin/sh", "/opt/genesis.sh")
	if err != nil {
		return nil, fmt.Errorf("error executing command: %w", err)
	}
	err = consensus.SetArgs("start", "--home=/root/.celestia-app", "--rpc.laddr=tcp://0.0.0.0:26657")
	if err != nil {
		return nil, fmt.Errorf("error setting args: %w", err)
	}
	err = consensus.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing instance: %w", err)
	}

	consensus.Start()
	consensus.WaitInstanceIsRunning()

	// Wait until validator reaches block height 1 or more
	for {
		status, err := consensus.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", "0.0.0.0"))
		if err != nil {
			return nil, fmt.Errorf("error executing command: %w", err)
		}
		blockHeight, err := app_utils.LatestBlockHeightFromStatus(status)
		if err != nil {
			return nil, fmt.Errorf("error getting block height: %w", err)
		}
		if blockHeight >= 1 {
			// Reached block height 1
			break
		}
		time.Sleep(1 * time.Second)
	}

	return consensus, nil
}

func ChainIdFromConsensus(consensus *knuu.Instance) (string, error) {
	status, err := consensus.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", "0.0.0.0"))
	if err != nil {
		logrus.Fatalf("Error getting status: %v", err)
	}
	chainId, err := ChainIdFromStatus(status)
	if err != nil {
		logrus.Fatalf("Error getting chain id from status: %v", err)
	}
	return chainId, nil
}

func GenesisHashFromConsensus(consensus *knuu.Instance) (string, error) {
	block, err := consensus.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/block?height=1", "0.0.0.0"))
	if err != nil {
		logrus.Fatalf("Error getting block: %v", err)
	}
	genesisHash, err := HashFromBlock(block)
	if err != nil {
		logrus.Fatalf("Error getting hash from block: %v", err)
	}
	return genesisHash, nil
}
