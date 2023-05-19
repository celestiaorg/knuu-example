package celestia_app

import (
	"fmt"
	"github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"os"
	"testing"
	"time"
)

func TestUpgrade(t *testing.T) {
	// Setup

	executor, err := knuu.NewExecutor()
	if err != nil {
		t.Fatalf("Error creating executor: %v", err)
	}

	validator, err := Instances["validator"].Clone()
	if err != nil {
		t.Fatalf("Error cloning instance: %v", err)
	}

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
			t.Log("Skipping cleanup")
			return
		}

		err = executor.Destroy()
		if err != nil {
			t.Fatalf("Error destroying executor: %v", err)
		}

		err = validator.Destroy()
		if err != nil {
			t.Fatalf("Error destroying instance: %v", err)
		}
	})

	// Test logic

	err = validator.Start()
	if err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}
	err = validator.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for instance to be running: %v", err)
	}

	validatorIP, err := validator.GetIP()
	if err != nil {
		t.Fatalf("Error getting validator IP: %v", err)
	}

	// Wait until validator reaches block height 1 or more
	t.Log("Waiting for validator to reach block height 1")
	for {
		status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", validatorIP))
		if err != nil {
			t.Fatalf("Error executing command: %v", err)
		}
		blockheight, err := utils.LatestBlockHeightFromStatus(status)
		if err != nil {
			t.Fatalf("Error getting blockheight from status: %v", err)
		}
		if blockheight >= 1 {
			t.Log("Validator reached block height 1")
			break
		}
		time.Sleep(1 * time.Second)
	}

	t.Log("Updating validator to v0.12.2")

	err = validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")
	if err != nil {
		t.Fatalf("Error setting image: %v", err)
	}

	// Check if blockheight is increasing, timeout after some time
	maxWait := 1 * time.Minute
	start := time.Now()
	status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", validatorIP))
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	blockHeight, err := utils.LatestBlockHeightFromStatus(status)
	if err != nil {
		t.Fatalf("Error getting blockheight from status: %v", err)
	}
	t.Logf("Waiting for validator blockheight to increase. Current: %d", blockHeight)
	for {
		if time.Since(start) > maxWait {
			t.Error("Timeout while waiting for blockheight to increase")
		}
		time.Sleep(1 * time.Second)
		status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", validatorIP))
		if err != nil {
			t.Fatalf("Error executing command: %v", err)
		}
		newBlockHeight, err := utils.LatestBlockHeightFromStatus(status)
		if err != nil {
			t.Fatalf("Error getting blockheight from status: %v", err)
		}
		if newBlockHeight > blockHeight {
			t.Log("Validator blockheight increased")
			break
		}
	}

}
