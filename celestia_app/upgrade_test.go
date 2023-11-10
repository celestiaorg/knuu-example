package celestia_app

// Deactivated for now
//func TestUpgrade(t *testing.T) {
//	t.Parallel()
//	// Setup
//
//	executor, err := knuu.NewExecutor()
//	if err != nil {
//		t.Fatalf("Error creating executor: %v", err)
//	}
//
//	validator, err := Instances["validator"].Clone()
//	if err != nil {
//		t.Fatalf("Error cloning instance: %v", err)
//	}
//	err = validator.AddVolume("/home/celestia/.celestia-app", "1Gi")
//	if err != nil {
//		t.Fatalf("Error adding volume: %v", err)
//	}
//
//	t.Cleanup(func() {
//		// Cleanup
//		if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
//			t.Log("Skipping cleanup")
//			return
//		}
//
//		err = executor.Destroy()
//		if err != nil {
//			t.Fatalf("Error destroying executor: %v", err)
//		}
//
//		err = validator.Destroy()
//		if err != nil {
//			t.Fatalf("Error destroying instance: %v", err)
//		}
//	})
//
//	// Test logic
//
//	err = validator.Start()
//	if err != nil {
//		t.Fatalf("Error starting instance: %v", err)
//	}
//	err = validator.WaitInstanceIsRunning()
//	if err != nil {
//		t.Fatalf("Error waiting for instance to be running: %v", err)
//	}
//
//	// Wait until validator reaches block height 1 or more
//	err = utils.WaitForHeight(executor, validator, 1)
//	if err != nil {
//		t.Fatalf("Error waiting for height: %v", err)
//	}
//
//	t.Log("Updating validator to v0.13.1")
//
//	err = validator.SetImage("ghcr.io/celestiaorg/celestia-app:v0.13.1")
//	if err != nil {
//		t.Fatalf("Error setting image: %v", err)
//	}
//
//	// Check if blockheight is increasing, timeout after some time
//	blockHeight, err := utils.GetHeight(executor, validator)
//	if err != nil {
//		t.Fatalf("Error getting block height: %v", err)
//	}
//	err = utils.WaitForHeight(executor, validator, blockHeight+1)
//	if err != nil {
//		t.Fatalf("Error waiting for height: %v", err)
//	}
//
//}
