package basic

import (
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

type FileTest struct {
	knuu.BaseTestRun
}

func (t *FileTest) Prepare() error {
	// Prepare test run, e.g. creating instances, setting configurations, etc.
	web := knuu.NewInstance("web")
	t.AddInstance(web)

	web.SetImage("nginx")
	web.AddPortTCP(80)
	web.ExecuteCommand([]string{"mkdir", "-p", "/usr/share/nginx/html"})
	web.AddFile("resources/index.html", "/usr/share/nginx/html/index.html", "0:0")
	web.Commit()

	checker := knuu.NewInstance("checker")
	t.AddInstance(checker)

	checker.SetImage("alpine")
	checker.SetCommand([]string{"sleep", "infinity"})
	checker.Commit()

	return nil
}

func (t *FileTest) Test() error {
	// Perform the actual tests

	web := t.Instances()[0]
	checker := t.Instances()[1]

	webIP := web.GetIP()

	web.Start()
	web.WaitInstanceIsRunning()
	checker.Start()
	checker.WaitInstanceIsRunning()

	wget := checker.ExecuteCommand([]string{"wget", "-q", "-O", "-", webIP})

	if wget != "Hello World!\n" {
		logrus.Fatalf("Hello World! not found")
	}

	return nil
}
