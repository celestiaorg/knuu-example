package celestia_app

import (
	"github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"os"
	"testing"
)

func TestPoolSync(t *testing.T) {
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

	id, err := utils.NodeIdFromNode(executor, validator)
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
	err = utils.WaitForHeight(executor, validator, 3)
	if err != nil {
		t.Fatalf("Error waiting for validator to reach block height 3: %v", err)
	}

	// Wait until full node reaches block height 3 or more but error out if it takes too long
	t.Log("Waiting for full nodes to reach block height 3")
	for _, full := range fullNodes.Instances() {
		err = utils.WaitForHeight(executor, full, 3)
		if err != nil {
			t.Fatalf("Error waiting for full node to reach block height 3: %v", err)
		}
	}
}
