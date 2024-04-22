package basic

import (
	"fmt"
	"os"
	"testing"

	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	t.Parallel()
	// Setup

	executor, err := knuu.NewExecutor()
	if err != nil {
		t.Fatalf("Error creating executor: %v", err)
	}

	web, err := knuu.NewInstance("web")
	if err != nil {
		t.Fatalf("Error creating instance '%v':", err)
	}
	err = web.SetImage("docker.io/nginx:latest")
	if err != nil {
		t.Fatalf("Error setting image '%v':", err)
	}
	web.AddPortTCP(80)
	_, err = web.ExecuteCommand("mkdir", "-p", "/usr/share/nginx/html")
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}
	err = web.AddFile("resources/html/index.html", "/usr/share/nginx/html/index.html", "0:0")
	if err != nil {
		t.Fatalf("Error adding file '%v':", err)
	}
	err = web.Commit()
	if err != nil {
		t.Fatalf("Error committing instance: %v", err)
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

		err = web.Destroy()
		if err != nil {
			t.Fatalf("Error destroying instance: %v", err)
		}
	})

	// Test logic

	webIP, err := web.GetIP()
	if err != nil {
		t.Fatalf("Error getting IP '%v':", err)
	}

	err = web.Start()
	if err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}
	err = web.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for instance to be running: %v", err)
	}

	wget, err := executor.ExecuteCommand("wget", "-q", "-O", "-", webIP)
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	assert.Equal(t, wget, "Hello World!\n")
}

func TestDownloadFileFromRunningInstance(t *testing.T) {
	t.Parallel()
	// Setup

	target, err := knuu.NewInstance("target")
	require.NoError(t, err, "Error creating instance")

	require.NoError(t, target.SetImage("alpine:latest"), "Error setting image")
	require.NoError(t, target.SetArgs("tail", "-f", "/dev/null"), "Error setting args") // Keep the container running
	require.NoError(t, target.Commit(), "Error committing instance")
	require.NoError(t, target.Start(), "Error starting instance")

	t.Cleanup(func() {
		// Cleanup
		if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
			t.Log("Skipping cleanup")
			return
		}

		require.NoError(t, target.Destroy(), "Error destroying instance")
	})

	// Test logic

	fileContent := "Hello World!"
	filePath := "/hello.txt"

	// Create a file in the target instance
	out, err := target.ExecuteCommand("echo", "-n", fileContent, ">", filePath)
	require.NoError(t, err, fmt.Sprintf("Error executing command: %v", out))

	gotContent, err := target.GetFileBytes(filePath)
	require.NoError(t, err, "Error getting file bytes")

	assert.Equal(t, fileContent, string(gotContent))
}
