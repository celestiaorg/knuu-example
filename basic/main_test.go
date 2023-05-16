package basic

import (
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

// DO NOT REMOVE THIS
func TestMain(m *testing.M) {
	err := knuu.Initialize()
	if err != nil {
		// This is very important. If this is not here, the program will fail
		if err.Error() == "InitReexec triggered re-exec" {
			return
		}
		logrus.Errorf("Error initializing: %v", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}
