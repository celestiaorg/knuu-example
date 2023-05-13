package basic

import (
	knuu "github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

type BasicTest struct {
	knuu.BaseTestRun
}

func (t *BasicTest) Prepare() error {
	// Setup your test
	logrus.Infof("MyTest: Preparing test")
	instance := knuu.NewInstance("alpine")
	t.AddInstance(instance)
	instance.SetImage("docker.io/library/alpine:latest")
	instance.SetCommand([]string{"sleep", "infinity"})
	instance.Commit()
	return nil
}

func (t *BasicTest) Test() error {
	// Run your test
	instance := t.Instances()[0]
	instance.Start()
	instance.WaitInstanceIsRunning()
	hello := instance.ExecuteCommand([]string{"echo", "Hello World!"})
	if hello != "Hello World!\n" {
		logrus.Fatalf("Hello World! not found")
	}
	logrus.Infof("MyTest: Running test")
	return nil
}
