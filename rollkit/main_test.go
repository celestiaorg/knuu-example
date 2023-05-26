package rollkit

import (
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

var Executor *knuu.Executor

func TestMain(m *testing.M) {
	err := knuu.Initialize()
	if err != nil {
		logrus.Fatalf("error initializing knuu: %v", err)
	}
	executor, err := knuu.NewExecutor()
	if err != nil {
		logrus.Fatalf("error creating executor: %v", err)
	}
	Executor = executor

	prepareInstances()

	exitVal := m.Run()

	err = Executor.Destroy()

	if err != nil {
		logrus.Fatalf("error destroying executor: %v", err)
	}
	os.Exit(exitVal)
}

var Instances = map[string]*knuu.Instance{}

func prepareInstances() {
	localCelestiaDevnet, err := knuu.NewInstance("local-celestia-devnet")
	if err != nil {
		logrus.Fatalf("error creating instance: %v", err)
	}
	err = localCelestiaDevnet.SetImage("ghcr.io/rollkit/local-celestia-devnet:v0.10.0")
	if err != nil {
		logrus.Fatalf("error setting image: %v", err)
	}
	err = localCelestiaDevnet.AddPortTCP(26657)
	if err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	err = localCelestiaDevnet.AddPortTCP(26659)
	if err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	err = localCelestiaDevnet.Commit()

	Instances["local-celestia-devnet"] = localCelestiaDevnet

	gm, err := knuu.NewInstance("gm")
	if err != nil {
		logrus.Fatalf("error creating instance: %v", err)
	}
	err = gm.SetImage("ghcr.io/rollkit/gm:v0.1.0")
	if err != nil {
		logrus.Fatalf("error setting image: %v", err)
	}
	err = gm.AddPortTCP(26656)
	if err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	err = gm.AddPortTCP(26657)
	if err != nil {
		logrus.Fatalf("error adding port: %v", err)
	}
	err = gm.SetUser("root")
	if err != nil {
		logrus.Fatalf("error setting user: %v", err)
	}
	_, err = gm.ExecuteCommand("apk", "--no-cache", "add", "curl", "jq")
	if err != nil {
		logrus.Fatalf("error executing command: %v", err)
	}
	err = gm.AddFile("resources/setup.sh", "/home/setup.sh", "10001:10001")
	if err != nil {
		logrus.Fatalf("error adding file: %v", err)
	}
	err = gm.AddFile("resources/entrypoint.sh", "/home/entrypoint.sh", "10001:10001")
	if err != nil {
		logrus.Fatalf("error adding file: %v", err)
	}
	err = gm.SetUser("rollkit")
	if err != nil {
		logrus.Fatalf("error setting user: %v", err)
	}
	_, err = gm.ExecuteCommand("/bin/sh", "/home/setup.sh")
	if err != nil {
		logrus.Fatalf("error executing command: %v", err)
	}
	err = gm.SetCommand([]string{"/bin/sh", "/home/entrypoint.sh"})
	if err != nil {
		logrus.Fatalf("error setting command: %v", err)
	}
	err = gm.Commit()
	if err != nil {
		logrus.Fatalf("error committing instance: %v", err)
	}
	Instances["gm"] = gm
}
