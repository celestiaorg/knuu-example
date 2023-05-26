#!/bin/sh

# print all commands before executing them
set -x

# Read out the namespace ID generated in setup.sh
# Note: If you want to regerneate the namespace ID, you need to rebuild the container
if [ -z "${NAMESPACE_ID}" ]; then
    # If DA_BLOCK_HEIGHT not already set, fetch current height
    # This is needed if you want to test multiple nodes
    NAMESPACE_ID=$(cat ${HOME}/NAMESPACE_ID)
fi

# query the DA Layer start height, in this case we are querying
# our local devnet at port 26657, the RPC. The RPC endpoint is
# to allow users to interact with Celestia's nodes by querying
# the node's state and broadcasting transactions on the Celestia
# network. The default port is 26657.
if [ -z "${DA_BLOCK_HEIGHT}" ]; then
    # If DA_BLOCK_HEIGHT not already set, fetch current height
    # This is needed if you want to test multiple nodes
    DA_BLOCK_HEIGHT=$(curl --silent ${RPC%/}${RPC:+/}block | jq -r '.result.block.header.height')
fi

# If DA_BLOCK_HEIGHT not set manually or fetched from RPC, exit
if [ -z "{$DA_BLOCK_HEIGHT}" ]; then
    echo "DA_BLOCK_HEIGHT is empty; ensure that the DA Layer is running and accessible"
    exit 1
fi
echo "DA_BLOCK_HEIGHT: ${DA_BLOCK_HEIGHT}"

exec gmd start "--rpc.laddr=tcp://0.0.0.0:26657" --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://'${DA_IP}':26659","timeout":60000000000,"fee":6000,"gas_limit":6000000}' --rollkit.namespace_id ${NAMESPACE_ID} --rollkit.da_start_height ${DA_BLOCK_HEIGHT} ${START_ARGS}
