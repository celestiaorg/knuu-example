package celestia_app

import (
	"fmt"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := knuu.Initialize()
	if err != nil {
		fmt.Errorf("error initializing knuu: %v", err)
	}
	identifier := knuu.Identifier()
	logrus.Infof("Unique identifier: %s", identifier)
	exitVal := m.Run()
	os.Exit(exitVal)
}

var Instances = map[string]*knuu.Instance{}
