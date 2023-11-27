package basic

import (
	"os"
	"testing"
	"time"

	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/stretchr/testify/assert"
)

func TestBittwister(t *testing.T) {
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

	webLatency, err := web.CloneWithName("web-latency")
	if err != nil {
		t.Fatalf("Error cloning instance '%v':", err)
	}
	webNoLatency, err := web.CloneWithName("web-no-latency")
	if err != nil {
		t.Fatalf("Error cloning instance '%v':", err)
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

		err = webLatency.Destroy()
		if err != nil {
			t.Fatalf("Error destroying instance: %v", err)
		}

		err = webNoLatency.Destroy()
		if err != nil {
			t.Fatalf("Error destroying instance: %v", err)
		}
	})

	// Test logic

	webLatencyIP, err := webLatency.GetIP()
	if err != nil {
		t.Fatalf("Error getting IP '%v':", err)
	}

	err = webLatency.SetLatency(1000)
	if err != nil {
		t.Fatalf("Error setting latency '%v':", err)
	}

	err = webLatency.Start()
	if err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}
	err = webLatency.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for instance to be running: %v", err)
	}

	startLatency := time.Now()
	wget, err := executor.ExecuteCommand("wget", "-q", "-O", "-", webLatencyIP)
	elapsedLatency := time.Since(startLatency)
	t.Logf("Running wget with latency took %s", elapsedLatency)
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	assert.Equal(t, wget, "Hello World!\n")

	webNoLatencyIP, err := webNoLatency.GetIP()
	if err != nil {
		t.Fatalf("Error getting IP '%v':", err)
	}

	err = webNoLatency.Start()
	if err != nil {
		t.Fatalf("Error starting instance: %v", err)
	}
	err = webNoLatency.WaitInstanceIsRunning()
	if err != nil {
		t.Fatalf("Error waiting for instance to be running: %v", err)
	}

	startNoLatency := time.Now()
	wgetNoLatency, err := executor.ExecuteCommand("wget", "-q", "-O", "-", webNoLatencyIP)
	elapsedNoLatency := time.Since(startNoLatency)
	t.Logf("Running wget without latency took %s", elapsedNoLatency)
	if err != nil {
		t.Fatalf("Error executing command '%v':", err)
	}

	assert.Equal(t, wgetNoLatency, "Hello World!\n")
	assert.True(t, elapsedLatency > elapsedNoLatency, "Instance with latency should take significantly longer")
}
