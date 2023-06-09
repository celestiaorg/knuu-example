package celestia_app

import (
	"github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"os"
	"testing"
	"time"
)

func TestNetworking(t *testing.T) {
	t.Parallel()
	// Setup

	executor, err := knuu.NewExecutor()
	if err != nil {
		t.Fatalf("Error creating executor: %v", err)
	}

	validator, err := Instances["validator"].Clone()
	if err != nil {
		t.Fatalf("Error cloning instance: %v", err)
	}
	full, err := Instances["full"].Clone()
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
		err = full.Destroy()
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
	// TODO: replace persistentPeers by utils call
	validatorIP, err := validator.GetIP()
	if err != nil {
		t.Fatalf("Error getting validator IP: %v", err)
	}
	id, err := utils.NodeIdFromNode(executor, validator)
	if err != nil {
		t.Fatalf("Error getting node id: %v", err)
	}
	persistentPeers := id + "@" + validatorIP + ":26656"
	err = full.SetArgs("start", "--home=/root/.celestia-app", "--rpc.laddr=tcp://0.0.0.0:26657", "--p2p.persistent_peers", persistentPeers)
	if err != nil {
		t.Fatalf("Error setting args: %v", err)
	}

	t.Log("Starting full nodes")
	err = full.Start()
	if err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}
	err = full.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for instance to be running: %v", err)
	}

	// Wait until validator reaches block height 1 or more
	err = utils.WaitForHeight(executor, full, 1)
	if err != nil {
		t.Fatalf("Error waiting for height: %v", err)
	}

	// Disable networking
	t.Log("Disabling networking")
	err = full.DisableNetwork()
	if err != nil {
		t.Fatalf("Error disabling network: %v", err)
	}

	// Get current block height
	height, err := utils.GetHeight(executor, full)
	if err != nil {
		t.Fatalf("Error getting height: %v", err)
	}

	// Fail if height increases more than 1 for next 1 minute
	t.Log("Waiting for height to not increase for 1 minute")
	timeout := time.After(1 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			goto afterTimeout
		case <-tick:
			newHeight, err := utils.GetHeight(executor, full)
			if err != nil {
				t.Fatalf("Error getting height: %v", err)
			}
			if newHeight > height+1 {
				t.Fatalf("Height increased from %d to %d", height, newHeight)
			}
			height = newHeight
		}
	}

afterTimeout:
	// Enable networking
	t.Log("Enabling networking")
	err = full.EnableNetwork()
	if err != nil {
		t.Fatalf("Error enabling network: %v", err)
	}

	// Check if blockheight is increasing, timeout after some time
	t.Log("Waiting for height to increase")
	height, err = utils.GetHeight(executor, validator)
	if err != nil {
		t.Fatalf("Error getting block height: %v", err)
	}
	err = utils.WaitForHeight(executor, validator, height+1)
	if err != nil {
		t.Fatalf("Error waiting for height: %v", err)
	}

}
