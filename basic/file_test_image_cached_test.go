package basic

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/celestiaorg/knuu/pkg/knuu"
)

func TestFileCached(t *testing.T) {
	t.Parallel()
	// Setup

	executor, err := knuu.NewExecutor()
	if err != nil {
		t.Fatalf("Error creating executor: %v", err)
	}

	const numberOfInstances = 10
	instances := make([]*knuu.Instance, numberOfInstances)
	errs := make(chan error, numberOfInstances)

	//var wg sync.WaitGroup
	for i := 0; i < numberOfInstances; i++ {
		instanceName := fmt.Sprintf("web%d", i+1)
		instance, err := knuu.NewInstance(instanceName)
		if err != nil {
			errs <- fmt.Errorf("Error creating instance '%v': %v", instanceName, err)
			return
		}
		err = instance.SetImage("docker.io/nginx:latest")
		if err != nil {
			errs <- fmt.Errorf("Error setting image for '%v': %v", instanceName, err)
			return
		}
		instance.AddPortTCP(80)
		_, err = instance.ExecuteCommand("mkdir", "-p", "/usr/share/nginx/html")
		if err != nil {
			errs <- fmt.Errorf("Error executing command for '%v': %v", instanceName, err)
			return
		}
		err = instance.Commit()
		if err != nil {
			errs <- fmt.Errorf("Error committing instance '%v': %v", instanceName, err)
			return
		}

		instances[i] = instance // Stores the instance reference for later use
	}

	var wgFolders sync.WaitGroup
	for i, instance := range instances {
		wgFolders.Add(1)
		go func(i int, instance *knuu.Instance) {
			defer wgFolders.Done()
			instanceName := fmt.Sprintf("web%d", i+1)
			// adding the folder after the Commit, it will help us to use a cached image.
			err = instance.AddFile("resources/html/index.html", "/usr/share/nginx/html/index.html", "0:0")
			if err != nil {
				errs <- fmt.Errorf("Error adding file to '%v': %v", instanceName, err)
				return
			}
		}(i, instance)
	}
	wgFolders.Wait()

	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("%v", err)
		}
	}

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") != "true" {
			err := executor.Destroy()
			if err != nil {
				t.Fatalf("Error destroying executor: %v", err)
			}

			for _, instance := range instances {
				if instance != nil {
					err := instance.Destroy()
					if err != nil {
						t.Fatalf("Error destroying instance: %v", err)
					}
				}
			}
		}
	})

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") != "true" {
			err := executor.Destroy()
			if err != nil {
				t.Fatalf("Error destroying executor: %v", err)
			}

			for _, instance := range instances {
				if instance != nil {
					err := instance.Destroy()
					if err != nil {
						t.Fatalf("Error destroying instance: %v", err)
					}
				}
			}
		}
	})

	// Test logic
	for _, instance := range instances {
		err = instance.StartWithoutWait()
		if err != nil {
			t.Fatalf("Error waiting for instance to be running: %v", err)
		}
	}

	for _, instance := range instances {
		webIP, err := instance.GetIP()
		if err != nil {
			t.Fatalf("Error getting IP: %v", err)
		}

		err = instance.WaitInstanceIsRunning()
		if err != nil {
			t.Fatalf("Error waiting for instance to be running: %v", err)
		}

		wget, err := executor.ExecuteCommand("wget", "-q", "-O", "-", webIP)
		if err != nil {
			t.Fatalf("Error executing command: %v", err)
		}

		assert.Equal(t, "Hello World!\n", wget)
	}
}
