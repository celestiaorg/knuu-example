package rollkit

import (
	"fmt"
	app_utils "github.com/celestiaorg/knuu-example/celestia_app/utils"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	t.Parallel()
	// Setup

	localCelestiaDevnet, err := Instances["local-celestia-devnet"].Clone()
	if err != nil {
		t.Fatalf("error cloning instance: %v", err)
	}
	gmAggregator, err := Instances["gm"].Clone()
	if err != nil {
		t.Fatalf("error cloning instance: %v", err)
	}
	gmFull, err := Instances["gm"].Clone()
	if err != nil {
		t.Fatalf("error cloning instance: %v", err)
	}

	t.Log("Start local-celestia-devnet")
	err = localCelestiaDevnet.Start()
	if err != nil {
		t.Fatalf("error starting instance: %v", err)
	}
	err = localCelestiaDevnet.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("error waiting for instance to be running: %v", err)
	}
	err = app_utils.WaitForHeight(Executor, localCelestiaDevnet, 1)
	if err != nil {
		t.Fatalf("error waiting for first block: %v", err)
	}

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
			t.Log("Skipping cleanup")
			return
		}
		err = localCelestiaDevnet.Destroy()
		if err != nil {
			t.Fatalf("error destroying instance: %v", err)
		}
		err = gmAggregator.Destroy()
		if err != nil {
			t.Fatalf("error destroying instance: %v", err)
		}
		err = gmFull.Destroy()
		if err != nil {
			t.Fatalf("error destroying instance: %v", err)
		}
	})

	// Test logic

	t.Log("Start gmAggregator")
	devnetIP, err := localCelestiaDevnet.GetIP()
	if err != nil {
		t.Fatalf("error getting IP: %v", err)
	}
	err = gmAggregator.SetEnvironmentVariable("RPC", fmt.Sprintf("http://%s:26657/", devnetIP))
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmAggregator.SetEnvironmentVariable("DA_IP", devnetIP)
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmAggregator.SetEnvironmentVariable("START_ARGS", "--rollkit.aggregator")
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmFull.SetEnvironmentVariable("DA_BLOCK_HEIGHT", "1")
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmAggregator.Start()
	if err != nil {
		t.Fatalf("error starting instance: %v", err)
	}
	err = gmAggregator.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("error waiting for instance to be running: %v", err)
	}
	err = app_utils.WaitForHeight(Executor, gmAggregator, 1)
	if err != nil {
		t.Fatalf("error waiting for first block: %v", err)
	}

	//gmAggregatorNodeID, err := app_utils.NodeIdFromNode(executor, gmAggregator)
	//if err != nil {
	//	t.Fatalf("error getting node ID: %v", err)
	//}
	//gmAggregatorNodeIP, err := gmAggregator.GetIP()
	//if err != nil {
	//	t.Fatalf("error getting IP: %v", err)
	//}

	t.Log("Start gmFull")
	err = gmFull.SetEnvironmentVariable("RPC", fmt.Sprintf("http://%s:26657/", devnetIP))
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmFull.SetEnvironmentVariable("DA_IP", devnetIP)
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmFull.SetEnvironmentVariable("DA_BLOCK_HEIGHT", "1")
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	//persistentPeers := gmAggregatorNodeID + "@" + gmAggregatorNodeIP + ":26656"
	//err = gmFull.SetEnvironmentVariable("START_ARGS", fmt.Sprintf("--p2p.seeds=%s", persistentPeers))
	//if err != nil {
	//	t.Fatalf("error setting environment variable: %v", err)
	//}
	err = gmFull.SetEnvironmentVariable("START_ARGS", "--home /home/rollkit/.gm2")
	if err != nil {
		t.Fatalf("error setting environment variable: %v", err)
	}
	err = gmFull.Start()
	if err != nil {
		t.Fatalf("error starting instance: %v", err)
	}
	err = gmFull.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("error waiting for instance to be running: %v", err)
	}

	aggregatorHeight, err := app_utils.GetHeight(Executor, gmAggregator)
	if err != nil {
		t.Fatalf("error getting height: %v", err)
	}
	t.Log("Wait for full node to catch up to aggregator")
	err = app_utils.WaitForHeight(Executor, gmFull, aggregatorHeight)
	if err != nil {
		t.Fatalf("error waiting for height: %v", err)
	}

}
