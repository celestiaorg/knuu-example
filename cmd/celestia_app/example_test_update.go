package celestia_app

import (
	"fmt"
	"time"

	"github.com/celestiaorg/knuu-example/cmd/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

type ExampleTestUpdate struct {
	knuu.BaseTestRun
}

func (t *ExampleTestUpdate) Prepare() error {
	validator := knuu.NewInstance("validator")
	t.AddInstance(validator)
	validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.1")
	validator.AddFile("resources/genesis.sh", "/opt/genesis.sh", "0:0")
	validator.ExecuteCommand([]string{"/bin/sh", "/opt/genesis.sh"})
	validator.Commit()
	return nil
}

func (t *ExampleTestUpdate) Test() error {
	validator := t.Instances()[0]
	validator.SetArgs([]string{"start", "--home=/root/.celestia-app"})
	validator.AddVolume("/root/.celestia-app", "1Gi")
	validator.Start()
	validator.WaitInstanceIsRunning()

	// Wait until validator reaches block height 1 or more
	logrus.Infof("Waiting for validator to reach block height 1")
	for {
		status := validator.ExecuteCommand([]string{"wget", "-q", "-O", "-", "localhost:26657/status"})
		blockheight := utils.LatestBlockHeightFromStatus(status)
		if blockheight >= 1 {
			logrus.Infof("Validator reached block height 1")
			break
		}
		time.Sleep(1 * time.Second)
	}

	logrus.Infof("Updating validator to v0.12.2")

	validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")

	// Check if blockheight is increasing, timeout after some time
	maxWait := 1 * time.Minute
	start := time.Now()
	blockheight := utils.LatestBlockHeightFromStatus(validator.ExecuteCommand([]string{"wget", "-q", "-O", "-", "localhost:26657/status"}))
	logrus.Infof("Waiting for validator blockheight to increase. Current: %d", blockheight)
	for {
		if time.Since(start) > maxWait {
			return fmt.Errorf("Timeout while waiting for blockheight to increase")
		}
		time.Sleep(1 * time.Second)
		newBlockheight := utils.LatestBlockHeightFromStatus(validator.ExecuteCommand([]string{"wget", "-q", "-O", "-", "localhost:26657/status"}))
		if newBlockheight > blockheight {
			logrus.Infof("Validator blockheight increased")
			break
		}
	}

	return nil
}
