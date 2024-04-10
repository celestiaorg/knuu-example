package basic

import (
	"fmt"
	"os"
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

	const numberOfInstances = 10 // Define how many instances you want to create
	instances := make([]*knuu.Instance, numberOfInstances)

	for i := 0; i < numberOfInstances; i++ {
		instanceName := fmt.Sprintf("web%d", i+1) // Generates a unique name for each instance
		instance, err := knuu.NewInstance(instanceName)
		if err != nil {
			t.Fatalf("Error creating instance '%v': %v", instanceName, err)
		}

		err = instance.SetImage("docker.io/nginx:latest")
		if err != nil {
			t.Fatalf("Error setting image for '%v': %v", instanceName, err)
		}

		instance.AddPortTCP(80)

		_, err = instance.ExecuteCommand("mkdir", "-p", "/usr/share/nginx/html")
		if err != nil {
			t.Fatalf("Error executing command for '%v': %v", instanceName, err)
		}

		err = instance.Commit()
		if err != nil {
			t.Fatalf("Error committing instance '%v': %v", instanceName, err)
		}

		err = instance.AddFile("resources/html/index.html", "/usr/share/nginx/html/index.html", "0:0")
		if err != nil {
			t.Fatalf("Error adding file to '%v': %v", instanceName, err)
		}

		instances[i] = instance // Stores the instance reference for later use
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

	// Test logic
	webdifferent, err := knuu.NewInstance("webdifferent")
	if err != nil {
		t.Fatalf("Error creating instance '%v':", err)
	}
	err = webdifferent.SetImage("docker.io/nginx:latest")
	if err != nil {
		t.Fatalf("Error setting image '%v':", err)
	}
	webdifferent.AddPortTCP(80)
	_, err = webdifferent.ExecuteCommand("mkdir", "-p", "/usr/share/nginx/html")
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	err = webdifferent.Commit()
	if err != nil {
		t.Fatalf("Error committing instance: %v", err)
	}

	err = webdifferent.AddFile("resources/html/smuu.html", "/usr/share/nginx/html/index.html", "0:0")
	if err != nil {
		t.Fatalf("Error adding file '%v':", err)
	}

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
			t.Log("Skipping cleanup")
			return
		}

		err = executor.Destroy()
		if err != nil {
			t.Fatalf("Error destroying executor: %v", err)
		}

		err = webdifferent.Destroy()
		if err != nil {
			t.Fatalf("Error destroying instance: %v", err)
		}
	})

	// Test logic
	for _, instance := range instances {
		webIP, err := instance.GetIP()
		if err != nil {
			t.Fatalf("Error getting IP: %v", err)
		}

		err = instance.Start()
		if err != nil {
			t.Fatalf("Error starting instance: %v", err)
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

	webIPdifferent, err := webdifferent.GetIP()
	if err != nil {
		t.Fatalf("Error getting IP '%v':", err)
	}

	err = webdifferent.Start()
	if err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}
	err = webdifferent.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for instance to be running: %v", err)
	}

	wget2, err := executor.ExecuteCommand("wget", "-q", "-O", "-", webIPdifferent)
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	assert.Equal(t, wget2, "It's-a Me, Smuu!\n")
}
