package celestia_app

import (
	"fmt"
	"github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"os"
	"testing"
	"time"
)

func TestPoolSync(t *testing.T) {
	// Setup

	executor, err := knuu.NewExecutor()
	if err != nil {
		t.Fatalf("Error creating executor: %v", err)
	}

	validator, err := knuu.NewInstance("validator")
	if err != nil {
		t.Fatalf("Error creating instance '%v':", err)
	}
	err = validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")
	if err != nil {
		t.Fatalf("Error setting image: %v", err)
	}
	err = validator.AddPortTCP(26656)
	if err != nil {
		t.Fatalf("Error adding port: %v", err)
	}
	err = validator.AddPortTCP(26657)
	if err != nil {
		t.Fatalf("Error adding port: %v", err)
	}
	err = validator.AddFile("resources/genesis.sh", "/opt/genesis.sh", "0:0")
	if err != nil {
		t.Fatalf("Error adding file: %v", err)
	}
	_, err = validator.ExecuteCommand("/bin/sh", "/opt/genesis.sh")
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	err = validator.SetArgs("start", "--home=/root/.celestia-app", "--rpc.laddr=tcp://0.0.0.0:26657")
	if err != nil {
		t.Fatalf("Error setting args: %v", err)
	}
	err = validator.Commit()
	if err != nil {
		t.Fatalf("Error committing instance: %v", err)
	}

	full, err := knuu.NewInstance("full")
	if err != nil {
		t.Fatalf("Error creating instance '%v':", err)
	}
	err = full.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")
	if err != nil {
		t.Fatalf("Error setting image: %v", err)
	}
	err = full.AddPortTCP(26657)
	if err != nil {
		t.Fatalf("Error adding port: %v", err)
	}
	genesis, err := validator.GetFileBytes("/root/.celestia-app/config/genesis.json")
	if err != nil {
		t.Fatalf("Error getting genesis: %v", err)
	}
	err = full.AddFileBytes(genesis, "/root/.celestia-app/config/genesis.json", "0:0")
	if err != nil {
		t.Fatalf("Error adding file: %v", err)
	}
	err = full.Commit()
	if err != nil {
		t.Fatalf("Error committing instance: %v", err)
	}
	// new InstancePool struct
	fullNodes := &knuu.InstancePool{}

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
			t.Fatalf("Error destroying validator: %v", err)
		}
		err = full.Destroy()
		if err != nil {
			t.Fatalf("Error destroying full: %v", err)
		}
		err = fullNodes.Destroy()
		if err != nil {
			t.Fatalf("Error destroying full nodes: %v", err)
		}
	})

	// Test logic

	err = validator.Start()
	if err != nil {
		t.Fatalf("Error starting validator: %v", err)
	}
	err = validator.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for validator to be running: %v", err)
	}
	validatorIP, err := validator.GetIP()
	if err != nil {
		t.Fatalf("Error getting validator IP: %v", err)
	}

	status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", validatorIP))
	if err != nil {
		t.Fatalf("Error getting status: %v", err)
	}

	id, err := utils.NodeIdFromStatus(status)
	if err != nil {
		t.Fatalf("Error getting node id: %v", err)
	}

	persistentPeers := id + "@" + validatorIP + ":26656"

	err = full.SetArgs("start", "--home=/root/.celestia-app", "--rpc.laddr=tcp://0.0.0.0:26657", "--p2p.persistent_peers", persistentPeers)
	if err != nil {
		t.Fatalf("Error setting args: %v", err)
	}
	fullNodes, err = full.CreatePool(5)
	if err != nil {
		t.Fatalf("Error creating pool: %v", err)
	}
	err = fullNodes.Start()
	if err != nil {
		t.Fatalf("Error starting full nodes: %v", err)
	}
	err = fullNodes.WaitInstancePoolIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for full nodes to be running: %v", err)
	}

	// Wait until validator reaches block height 3 or more
	t.Log("Waiting for validator to reach block height 3")
	for {
		status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", validatorIP))
		if err != nil {
			t.Fatalf("Error getting status: %v", err)
		}
		blockheightVal, err := utils.LatestBlockHeightFromStatus(status)
		if err != nil {
			t.Fatalf("Error getting latest block height: %v", err)
		}
		if blockheightVal >= 3 {
			t.Logf("Validator reached block height 3")
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Wait until full node reaches block height 3 or more but error out if it takes too long
	t.Log("Waiting for full nodes to reach block height 3")
	maxWait := 5 * time.Second
	start := time.Now()
	for _, full := range fullNodes.Instances() {
		for {
			fullNodeIP, err := full.GetIP()
			if err != nil {
				t.Fatalf("Error getting full node IP: %v", err)
			}
			status, err := executor.ExecuteCommand("wget", "-q", "-O", "-", fmt.Sprintf("%s:26657/status", fullNodeIP))
			if err != nil {
				t.Fatalf("Error getting status: %v", err)
			}
			blockheightFull, err := utils.LatestBlockHeightFromStatus(status)
			if err != nil {
				t.Fatalf("Error getting latest block height: %v", err)
			}
			if blockheightFull >= 3 {
				t.Log("Full node reached block height 3")
				break
			}
			time.Sleep(1 * time.Second)
			if time.Since(start) > maxWait {
				t.Fatalf("Full node did not reach block height 3 in time")
			}
		}
	}
}
