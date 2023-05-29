package utils

import (
	"fmt"
	"github.com/celestiaorg/knuu/pkg/knuu"
)

func CreateAndStartConsensus(executor *knuu.Executor) (*knuu.Instance, error) {
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
	WaitForHeight(executor, consensus, 1)

	return consensus, nil
}
