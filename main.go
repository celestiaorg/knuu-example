package main

import (
	"github.com/celestiaorg/knuu-example/cmd/basic"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/sirupsen/logrus"
)

func main() {

	err := knuu.Initialize()
	if err != nil {
		// This is very important. If this is not here, the program will fail
		if err.Error() == "InitReexec triggered re-exec" {
			return
		}
		logrus.Errorf("Error initializing: %w", err)
	}

	// TODO: INVARIANT: name of container must be unique

	testRuns := &knuu.TestRuns{}

	// Basic tests
	testRuns.AddTest(&basic.BasicTest{})
	// testRuns.AddTest(&basic.FileTest{})

	// testRuns.AddTest(&ExampleTestRun{})
	// testRuns.AddTest(&celestia_app.ExampleTestSync{})
	// testRuns.AddTest(&celestia_app.ExampleTestUpdate{})
	// testRuns.AddTest(&celestia_app.ExampleTestPoolSync{})

	err = testRuns.RunAll()
	if err != nil {
		logrus.Errorf("Error running testruns: %w", err)
	}
}
