package celestia_app

import (
	"time"

	"github.com/celestiaorg/knuu-example/cmd/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

type ExampleTestPoolSync struct {
	knuu.BaseTestRun
}

func (t *ExampleTestPoolSync) Prepare() error {
	// Prepare test run, e.g. creating instances, setting configurations, etc.
	validator := knuu.NewInstance("validator")
	validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")
	validator.AddPortTCP(26656)
	validator.AddPortTCP(26657)
	validator.AddFile("resources/genesis.sh", "/opt/genesis.sh", "0:0")
	validator.ExecuteCommand([]string{"/bin/sh", "/opt/genesis.sh"})

	validator.Commit()

	t.AddInstance(validator)

	full := knuu.NewInstance("full")
	full.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")
	full.AddPortTCP(26657)
	genesis := validator.GetFileBytes("/root/.celestia-app/config/genesis.json")
	full.AddFileBytes(genesis, "/root/.celestia-app/config/genesis.json", "0:0")
	full.Commit()

	t.AddInstance(full)

	return nil
}

func (t *ExampleTestPoolSync) Test() error {
	// Perform the actual tests

	validator := t.Instances()[0]
	full := t.Instances()[1]
	validator.SetArgs([]string{"start", "--home=/root/.celestia-app"})
	validator.Start()
	validator.WaitInstanceIsRunning()

	// TODO: Instead of executing wget on the instance, we could curl the status from the pupeteer
	status := validator.ExecuteCommand([]string{"wget", "-q", "-O", "-", "localhost:26657/status"})

	id := utils.NodeIdFromStatus(status)

	persistentPeers := id + "@" + validator.GetIP() + ":26656"

	full.SetArgs([]string{"start", "--home=/root/.celestia-app", "--p2p.persistent_peers", persistentPeers})
	fullNodes := full.CreatePool(5)
	t.AddInstancePool(fullNodes)
	fullNodes.Start()
	fullNodes.WaitInstancePoolIsRunning()

	// Wait until validator reaches block height 3 or more
	logrus.Infof("Waiting for validator to reach block height 3")
	for {
		status = validator.ExecuteCommand([]string{"wget", "-q", "-O", "-", "localhost:26657/status"})
		blockheightVal := utils.LatestBlockHeightFromStatus(status)
		if blockheightVal >= 3 {
			logrus.Infof("Validator reached block height 3")
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Wait until full node reaches block height 3 or more but error out if it takes too long
	logrus.Infof("Waiting for full nodes to reach block height 3")
	maxWait := 5 * time.Second
	start := time.Now()
	for _, full := range fullNodes.Instances() {
		for {
			status = full.ExecuteCommand([]string{"wget", "-q", "-O", "-", "localhost:26657/status"})
			blockheightFull := utils.LatestBlockHeightFromStatus(status)
			if blockheightFull >= 3 {
				logrus.Infof("Full node reached block height 3")
				break
			}
			time.Sleep(1 * time.Second)
			if time.Since(start) > maxWait {
				logrus.Fatalf("Full node did not reach block height 3 in time")
			}
		}
	}

	return nil
}
