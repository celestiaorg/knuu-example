test-basic:
	go test -v ./basic -timeout 120m

test-bittwister-packetloss:
	KNUU_TIMEOUT=120m go test -v ./basic --run=TestBittwister_Packetloss -timeout 60m -count=1

test-bittwister-bandwidth:
	KNUU_TIMEOUT=120m go test -v ./basic --run=TestBittwister_Bandwidth -timeout 60m -count=1

test-bittwister-latency:
	KNUU_TIMEOUT=120m go test -v ./basic --run=TestBittwister_Latency -timeout 60m -count=1

test-bittwister-jitter:
	KNUU_TIMEOUT=120m go test -v ./basic --run=TestBittwister_Jitter -timeout 60m -count=1

test-celestia-app:
	go test -v ./celestia_app

test-celestia-node:
	go test -v ./celestia_node

test-all:
	KNUU_TIMEOUT=300m go test -v ./... -timeout 120m

.PHONY: test-all test-basic test-bittwister-packetloss test-bittwister-bandwidth test-bittwister-latency test-bittwister-jitter test-celestia-app test-celestia-node