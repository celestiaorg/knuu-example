package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func NodeIdFromStatus(status string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(status), &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling status: %w", err)
	}
	resultData := result["result"].(map[string]interface{})
	nodeInfo := resultData["node_info"].(map[string]interface{})
	id := nodeInfo["id"].(string)
	return id, nil
}

func LatestBlockHeightFromStatus(status string) (int64, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(status), &result)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling status: %w", err)
	}
	resultData := result["result"].(map[string]interface{})
	syncInfo := resultData["sync_info"].(map[string]interface{})
	latestBlockHeight := syncInfo["latest_block_height"].(string)
	latestBlockHeightInt, err := strconv.ParseInt(latestBlockHeight, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting latest block height to int: %w", err)
	}
	return latestBlockHeightInt, nil
}

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
