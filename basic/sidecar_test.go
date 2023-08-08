package basic

import (
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSidecar(t *testing.T) {
	t.Parallel()
	// Setup

	executor, err := knuu.NewExecutor()
	if err != nil {
		t.Fatalf("Error creating executor: %v", err)
	}

	// web instance
	web, err := knuu.NewInstance("web")
	if err != nil {
		t.Fatalf("Error creating instance '%v':", err)
	}
	if err := web.SetImage("docker.io/nginx:latest"); err != nil {
		t.Fatalf("Error setting image '%v':", err)
	}
	if err := web.AddPortTCP(80); err != nil {
		t.Fatalf("Error adding port: %v", err)
	}
	if _, err := web.ExecuteCommand("mkdir", "-p", "/usr/share/nginx/html"); err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}
	if err := web.AddFile("resources/html/index.html", "/usr/share/nginx/html/index.html", "0:0"); err != nil {
		t.Fatalf("Error adding file '%v':", err)
	}
	if err := web.AddVolume("/usr/share/nginx/html", "1Gi"); err != nil {
		t.Fatalf("Error adding volume '%v':", err)
	}
	if err := web.Commit(); err != nil {
		t.Fatalf("Error committing instance: %v", err)
	}

	// sidecar instance
	sidecar, err := knuu.NewInstance("sidecar")
	if err != nil {
		t.Fatalf("Error creating instance '%v':", err)
	}
	err = sidecar.SetImage("alpine:latest")
	if err != nil {
		t.Fatalf("Error setting image '%v':", err)
	}
	if err := sidecar.SetCommand("sleep", "infinity"); err != nil {
		t.Fatalf("Error setting command '%v':", err)
	}
	if err := sidecar.Commit(); err != nil {
		t.Fatalf("Error committing instance: %v", err)
	}
	if err := web.AddSidecar(sidecar); err != nil {
		t.Fatalf("Error adding sidecar: %v", err)
	}

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
			t.Log("Skipping cleanup")
			return
		}

		if err = executor.Destroy(); err != nil {
			t.Fatalf("Error destroying executor: %v", err)
		}

		if err = web.Destroy(); err != nil {
			t.Fatalf("Error destroying instance: %v", err)
		}
	})

	// Test logic

	webIP, err := web.GetIP()
	if err != nil {
		t.Fatalf("Error getting IP '%v':", err)
	}

	require.Error(t, sidecar.Start())
	if err := web.Start(); err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}

	wgetExecutor, err := executor.ExecuteCommand("wget", "-q", "-O", "-", webIP)
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	assert.Equal(t, wgetExecutor, "Hello World!\n")

	wgetSidecar, err := sidecar.ExecuteCommand("wget", "-q", "-O", "-", webIP)
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	assert.Equal(t, wgetSidecar, "Hello World!\n")
}
