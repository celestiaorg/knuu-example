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
    --chain-id $CHAINID \
    --evm-address 0x966e6f22781EF6a6A82BBB4DB3df8E225DfD9488 # private key: da6ed55cb2894ac2c9c10209c09de8e8b9d109b910338d5bf3d747a7e1fc9eb9

celestia-appd collect-gentxs
