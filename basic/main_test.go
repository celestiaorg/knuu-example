package basic

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	t := time.Now()
	identifier := fmt.Sprintf("%s_%03d", t.Format("20060102_150405"), t.Nanosecond()/1e6)
	err := knuu.InitializeWithIdentifier(identifier)
	if err != nil {
		logrus.Fatalf("error initializing knuu: %v", err)
	}
	logrus.Infof("Unique identifier: %s", identifier)
	exitVal := m.Run()
	os.Exit(exitVal)
}
