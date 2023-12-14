test-basic:
	go test -v ./basic -timeout 60m

test-bittwister-packetloss:
	go test -v ./basic --run=TestBittwister_Packetloss -timeout 30m -count=1

test-bittwister-bandwidth:
	go test -v ./basic --run=TestBittwister_Bandwidth -timeout 30m -count=1

test-bittwister-latency:
	go test -v ./basic --run=TestBittwister_Latency -timeout 30m -count=1

test-bittwister-jitter:
	go test -v ./basic --run=TestBittwister_Jitter -timeout 30m -count=1

test-celestia-app:
	go test -v ./celestia_app

test-celestia-node:
	go test -v ./celestia_node

test-all:
	go test -v ./... -timeout 60m

.PHONY: test-all test-basic test-bittwister-packetloss test-bittwister-bandwidth test-bittwister-latency test-bittwister-jitter test-celestia-app test-celestia-node