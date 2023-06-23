package main

import (
	app_utils "github.com/celestiaorg/knuu-example/celestia_app/utils"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

func main() {

	err := knuu.Initialize()
	if err != nil {
		logrus.Fatalf("error initializing knuu: %v", err)
	}
	identifier := knuu.Identifier()
	logrus.Infof("Unique identifier: %s", identifier)

	executor, err := knuu.NewExecutor()
	if err != nil {
		logrus.Fatalf("Error creating executor: %v", err)
	}

	//// VALIDATOR ////
	validator, err := knuu.NewInstance("validator")
	if err != nil {
		logrus.Fatalf("error creating instance: %v", err)
	}
	if err := validator.SetImage("ghcr.io/celestiaorg/celestia-app:v1.0.0-rc4"); err != nil {
		logrus.Fatalf("error setting image: %v", err)
	}
	if err := validator.AddPortTCP(26656); err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	if err := validator.AddPortTCP(26657); err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	if err := validator.AddFile("resources/genesis.sh", "/opt/genesis.sh", "10001:10001"); err != nil {
		logrus.Fatalf("error adding file: %v", err)
	}
	result, err := validator.ExecuteCommand("/bin/sh", "/opt/genesis.sh")
	if err != nil {
		logrus.Errorf("command result: %s", result)
		logrus.Fatalf("Error executing command: %v", err)
	}
	if err := validator.SetMemory("200Mi", "200Mi"); err != nil {
		logrus.Fatalf("error setting memory: %v", err)
	}
	if err := validator.SetCPU("300m"); err != nil {
		logrus.Fatalf("error setting cpu: %v", err)
	}
	if err := validator.Commit(); err != nil {
		logrus.Fatalf("error committing instance: %v", err)
	}
	genesisFile, err := validator.GetFileBytes("/home/celestia/.celestia-app/config/genesis.json")
	if err != nil {
		logrus.Fatalf("error getting genesis: %v", err)
	}
	if err := validator.SetArgs("start", "--rpc.laddr=tcp://0.0.0.0:26657"); err != nil {
		logrus.Fatalf("error setting args: %v", err)
	}
	if err := validator.Start(); err != nil {
		logrus.Fatalf("error starting instance: %v", err)
	}

	//// FULL ////
	full, err := knuu.NewInstance("full")

	if err := full.SetImage("ghcr.io/celestiaorg/celestia-app:v1.0.0-rc4"); err != nil {
		logrus.Fatalf("error running knuu: %v", err)
	}
	if err := full.AddPortTCP(26656); err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	if err := full.AddPortTCP(26657); err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	err = full.AddFileBytes(genesisFile, "/home/celestia/.celestia-app/config/genesis.json", "10001:10001")
	if err != nil {
		logrus.Fatalf("error adding file: %v", err)
	}
	if err := full.SetMemory("200Mi", "200Mi"); err != nil {
		logrus.Fatalf("error setting memory: %v", err)
	}
	if err := full.SetCPU("300m"); err != nil {
		logrus.Fatalf("error setting cpu: %v", err)
	}
	if err := full.Commit(); err != nil {
		logrus.Fatalf("error committing instance: %v", err)
	}
	persistentPeers, err := app_utils.GetPersistentPeers(executor, []*knuu.Instance{validator})
	if err != nil {
		logrus.Fatalf("error getting persistent peers: %v", err)
	}
	if err := full.SetArgs("start", "--rpc.laddr=tcp://0.0.0.0:26657", "--p2p.persistent_peers", persistentPeers); err != nil {
		logrus.Fatalf("error setting args: %v", err)
	}
	full_pool, err := full.CreatePool(5)
	if err != nil {
		logrus.Fatalf("error creating pool: %v", err)
	}
	if err := full_pool.Start(); err != nil {
		logrus.Fatalf("error starting instance: %v", err)
	}

	//// BRIDGE ////

}
