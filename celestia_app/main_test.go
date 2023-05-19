package celestia_app

import (
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := knuu.Initialize()
	if err != nil {
		logrus.Fatalf("error initializing knuu: %v", err)
	}
	prepareInstances(m)
	exitVal := m.Run()
	os.Exit(exitVal)
}

var Instances = map[string]*knuu.Instance{}

func prepareInstances(m *testing.M) {
	validator, err := knuu.NewInstance("validator")
	if err != nil {
		logrus.Fatalf("Error creating instance '%v':", err)
	}
	err = validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.1")
	if err != nil {
		logrus.Fatalf("Error setting image: %v", err)
	}
	err = validator.AddPortTCP(26656)
	if err != nil {
		logrus.Fatalf("Error adding port: %v", err)
	}
	err = validator.AddPortTCP(26657)
	if err != nil {
		logrus.Fatalf("Error adding port: %v", err)
	}
	err = validator.AddFile("resources/genesis.sh", "/opt/genesis.sh", "0:0")
	if err != nil {
		logrus.Fatalf("Error adding file: %v", err)
	}
	_, err = validator.ExecuteCommand("/bin/sh", "/opt/genesis.sh")
	if err != nil {
		logrus.Fatalf("Error executing command: %v", err)
	}
	err = validator.SetArgs("start", "--home=/root/.celestia-app", "--rpc.laddr=tcp://0.0.0.0:26657")
	if err != nil {
		logrus.Fatalf("Error setting args: %v", err)
	}
	err = validator.AddVolume("/root/.celestia-app", "1Gi")
	if err != nil {
		logrus.Fatalf("Error adding volume: %v", err)
	}
	err = validator.Commit()
	if err != nil {
		logrus.Fatalf("Error committing instance: %v", err)
	}

	Instances["validator"] = validator

	full, err := knuu.NewInstance("full")
	if err != nil {
		logrus.Fatalf("Error creating instance '%v':", err)
	}
	err = full.SetImage("ghcr.io/celestiaorg/celestia-app:v0.12.2")
	if err != nil {
		logrus.Fatalf("Error setting image: %v", err)
	}
	err = full.AddPortTCP(26656)
	if err != nil {
		logrus.Fatalf("Error adding port: %v", err)
	}
	err = full.AddPortTCP(26657)
	if err != nil {
		logrus.Fatalf("Error adding port: %v", err)
	}
	genesis, err := validator.GetFileBytes("/root/.celestia-app/config/genesis.json")
	if err != nil {
		logrus.Fatalf("Error getting genesis: %v", err)
	}
	err = full.AddFileBytes(genesis, "/root/.celestia-app/config/genesis.json", "0:0")
	if err != nil {
		logrus.Fatalf("Error adding file: %v", err)
	}
	err = full.Commit()
	if err != nil {
		logrus.Fatalf("Error committing instance: %v", err)
	}

	Instances["full"] = full
}
