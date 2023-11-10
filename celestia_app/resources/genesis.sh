#!/bin/sh

CHAINID="test"

# Build genesis file incl account for passed address
coins="1000000000000000utia"
celestia-appd init $CHAINID --chain-id $CHAINID
celestia-appd keys add node --keyring-backend="test"
# this won't work because some proto types are declared twice and the logs output to stdout (dependency hell involving iavl)
celestia-appd add-genesis-account $(celestia-appd keys show node -a --keyring-backend="test") $coins
celestia-appd gentx node 5000000000utia \
    --keyring-backend="test" \
    --chain-id $CHAINID

celestia-appd collect-gentxs
