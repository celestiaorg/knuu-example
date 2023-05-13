package utils

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/sirupsen/logrus"
)

func NodeIdFromStatus(status string) string {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(status), &result)
	if err != nil {
		log.Fatal(err)
	}
	resultData := result["result"].(map[string]interface{})
	nodeInfo := resultData["node_info"].(map[string]interface{})
	id := nodeInfo["id"].(string)
	return id
}

func LatestBlockHeightFromStatus(status string) int64 {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(status), &result)
	if err != nil {
		log.Fatal(err)
	}
	resultData := result["result"].(map[string]interface{})
	syncInfo := resultData["sync_info"].(map[string]interface{})
	latestBlockHeight := syncInfo["latest_block_height"].(string)
	latestBlockHeightInt, err := strconv.ParseInt(latestBlockHeight, 10, 64)
	if err != nil {
		logrus.Fatalf("Error converting latest block height to int: %w", err)
	}
	return latestBlockHeightInt
}
