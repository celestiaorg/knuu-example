package basic

import (
	"fmt"
	"github.com/celestiaorg/knuu/pkg/knuu"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/rand"
	"os"
	"strings"
	"testing"
)

func TestExternalDns(t *testing.T) {
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

	if err := web.SetImage("docker.io/nginx:latest"); err != nil {
		t.Fatalf("Error setting image '%v':", err)
	}
	if err := web.AddPortTCP(80); err != nil {
		t.Fatalf("Error adding port '%v':", err)
	}
	if err := web.Commit(); err != nil {
		t.Fatalf("Error committing instance: %v", err)
	}

	// random prefix (lowercase RFC 1123 subdomain)
	prefix := strings.ToLower(rand.Str(10))

	if err := web.AddExternalDns(fmt.Sprintf("%s-knuu.celestia-robusta.com", prefix)); err != nil {
		t.Fatalf("Error adding external dns '%v':", err)
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

	assert.NotEmpty(t, wget)
}
