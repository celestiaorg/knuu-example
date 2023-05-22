package basic

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
	exitVal := m.Run()
	os.Exit(exitVal)
}
