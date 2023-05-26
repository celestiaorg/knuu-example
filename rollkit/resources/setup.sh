#!/bin/sh

# set variables for the chain
VALIDATOR_NAME=validator1
NODE_NAME=node1
CHAIN_ID=gm
VALIDATOR_KEY_NAME=gm-key
VALIDATOR_KEY_2_NAME=gm-key-2
NODE_KEY_NAME=gm-node-key
TOKEN_AMOUNT="10000000000000000000000000stake"
STAKING_AMOUNT="1000000000stake"

# initialize the validator with the chain ID you set
gmd init ${VALIDATOR_NAME} --chain-id ${CHAIN_ID}
gmd init ${NODE_NAME} --chain-id ${CHAIN_ID} --home ${HOME}/.gm2

# add keys for key 1 and key 2 to keyring-backend test
gmd keys add ${VALIDATOR_KEY_NAME} --keyring-backend test
gmd keys add ${VALIDATOR_KEY_2_NAME} --keyring-backend test
gmd keys add ${NODE_KEY_NAME} --keyring-backend test --home ${HOME}/.gm2

# add these as genesis accounts
gmd add-genesis-account ${VALIDATOR_KEY_NAME} ${TOKEN_AMOUNT} --keyring-backend test
gmd add-genesis-account ${VALIDATOR_KEY_2_NAME} ${TOKEN_AMOUNT} --keyring-backend test

# set the staking amounts in the genesis transaction
gmd gentx ${VALIDATOR_KEY_NAME} ${STAKING_AMOUNT} --chain-id ${CHAIN_ID} --keyring-backend test

# collect genesis transactions
gmd collect-gentxs

# All should have the same genesis.json
cp ${HOME}/.gm/config/genesis.json ${HOME}/.gm2/config/genesis.json

# Generate a random namespace ID for your rollup to post blocks to and save it to a file
NAMESPACE_ID=$(echo ${RANDOM} | md5sum | head -c 16; echo)
echo ${NAMESPACE_ID} > ${HOME}/NAMESPACE_ID
