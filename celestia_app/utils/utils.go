package utils

import (
	"fmt"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"time"
)

type JSONRPCError struct {
	Code    int
	Message string
	Data    string
}

func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("JSONRPC Error - Code: %d, Message: %s, Data: %s", e.Code, e.Message, e.Data)
}

func NodeIdFromNode(executor *knuu.Executor, node *knuu.Instance) (string, error) {
	nodeIP, err := node.GetIP()
	if err != nil {
		return "", fmt.Errorf("error getting node ip: %w", err)
	}
	status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", nodeIP))
	if err != nil {
		return "", fmt.Errorf("error getting status: %v", err)
	}

	id, err := nodeIdFromStatus(status)
	if err != nil {
		return "", fmt.Errorf("error getting node id: %v", err)
	}
	return id, nil
}

func GetHeight(executor *knuu.Executor, app *knuu.Instance) (int64, error) {
	appIP, err := app.GetIP()
	if err != nil {
		return 0, fmt.Errorf("error getting app ip: %w", err)
	}
	status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", appIP))
	if err != nil {
		return 0, fmt.Errorf("error executing command: %w", err)
	}
	blockHeight, err := latestBlockHeightFromStatus(status)
	if err != nil {
		return 0, fmt.Errorf("error getting block height: %w", err)
	}
	return blockHeight, nil
}

func WaitForHeight(executor *knuu.Executor, app *knuu.Instance, height int64) error {
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(1 * time.Second)

	maxRetries := 5

	appIP, err := app.GetIP()
	if err != nil {
		return fmt.Errorf("error getting app ip: %w", err)
	}

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for the first block")
		case <-tick:
			var blockHeight int64
			var err error
			for retries := 0; retries < maxRetries; retries++ {
				status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", appIP))
				if err != nil {
					return fmt.Errorf("error executing command: %w", err)
				}
				blockHeight, err = latestBlockHeightFromStatus(status)
				if err != nil {
					if _, ok := err.(*JSONRPCError); !ok {
						// If the error is not a JSONRPCError, stop and return it
						return fmt.Errorf("error getting block height: %w", err)
					}
					// If the error is a JSONRPCError, sleep and then retry
					if retries < maxRetries-1 {
						time.Sleep(1 * time.Second)
					}
				} else {
					// If no error, break the retry loop
					break
				}
			}
			if err != nil {
				return fmt.Errorf("error getting block height: %w", err)
			}
			if blockHeight >= height {
				// Reached block height `height`
				return nil
			}
		}
	}
}

func ChainId(executor *knuu.Executor, app *knuu.Instance) (string, error) {
	appIP, err := app.GetIP()
	if err != nil {
		return "", fmt.Errorf("error getting app ip: %w", err)
	}

	status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", appIP))
	if err != nil {
		return "", fmt.Errorf("error executing command: %w", err)
	}
	chainId, err := chainIdFromStatus(status)
	if err != nil {
		return "", fmt.Errorf("error getting chain id: %w", err)
	}
	return chainId, nil
}

func GenesisHash(executor *knuu.Executor, app *knuu.Instance) (string, error) {
	appIP, err := app.GetIP()
	if err != nil {
		return "", fmt.Errorf("error getting app ip: %w", err)
	}
	block, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/block?height=1", appIP))
	if err != nil {
		return "", fmt.Errorf("error getting block: %v", err)
	}
	genesisHash, err := hashFromBlock(block)
	if err != nil {
		return "", fmt.Errorf("error getting hash from block: %v", err)
	}
	return genesisHash, nil
}

func GetPersistentPeers(executor *knuu.Executor, apps []*knuu.Instance) (string, error) {
	var persistentPeers string
	for _, app := range apps {
		validatorIP, err := app.GetIP()
		if err != nil {
			return "", fmt.Errorf("error getting validator IP: %v", err)
		}
		id, err := NodeIdFromNode(executor, app)
		if err != nil {
			return "", fmt.Errorf("error getting node id: %v", err)
		}
		persistentPeers += id + "@" + validatorIP + ":26656" + ","
	}
	return persistentPeers[:len(persistentPeers)-1], nil
}
