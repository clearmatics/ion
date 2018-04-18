#!/bin/bash

## Run and attach network 0
tmux new -s chaina -n testrpc_migration_a -d
tmux send-keys -t chaina "npm run testrpca" C-m
tmux split-window -h -t chaina:testrpc_migration_a
tmux send-keys -t chaina "sleep 1s" C-m
tmux send-keys -t chaina "rm -rf build/contracts/*" C-m
tmux send-keys -t chaina "truffle compile" C-m
tmux send-keys -t chaina "truffle migrate --network testrpca" C-m
tmux send-keys -t chaina "sleep 1s" C-m
tmux send-keys -t chaina "geth attach http://127.0.0.1:8545" C-m

## Run and attach network 1
tmux new -s chainb -n testrpc_migration_b -d
tmux send-keys -t chainb "npm run testrpcb" C-m
tmux split-window -h -t chainb:testrpc_migration_b
tmux send-keys -t chainb "sleep 1s" C-m
tmux send-keys -t chainb "truffle migrate --network testrpcb" C-m
tmux send-keys -t chainb "sleep 1s" C-m
tmux send-keys -t chainb "geth attach http://127.0.0.1:8546" C-m

## Now deploy the token on the chain
TokenAbi=`cat abi/Token.abi`
# Define some key info
tmux send-keys -t chaina "var own = eth.accounts[0]" C-m
tmux send-keys -t chaina "var send = eth.accounts[1]" C-m
tmux send-keys -t chaina "var recv = eth.accounts[2]" C-m
tmux send-keys -t chaina "var tokenAddr = '0x9561c133dd8580860b6b7e504bc5aa500f0f06a7'" C-m
tmux send-keys -t chaina "var lockAddr = '0xe982e462b094850f12af94d21d470e21be9d0e9c'" C-m
tmux send-keys -t chaina "var linkAddr = '0xc89ce4735882c9f0f0fe26686c53074e09b0d550'" C-m

# deploy the token
tmux send-keys -t chaina "var tokenAbi = $TokenAbi" C-m
tmux send-keys -t chaina "var tokenContract = eth.contract(tokenAbi)" C-m
tmux send-keys -t chaina "var token = tokenContract.at(tokenAddr)" C-m

# Define some key info
tmux send-keys -t chainb "var own = eth.accounts[0]" C-m
tmux send-keys -t chainb "var send = eth.accounts[1]" C-m
tmux send-keys -t chainb "var recv = eth.accounts[2]" C-m
tmux send-keys -t chainb "var tokenAddr = '0x9561c133dd8580860b6b7e504bc5aa500f0f06a7'" C-m
tmux send-keys -t chainb "var lockAddr = '0xe982e462b094850f12af94d21d470e21be9d0e9c'" C-m
tmux send-keys -t chainb "var linkAddr = '0xc89ce4735882c9f0f0fe26686c53074e09b0d550'" C-m

# deploy the token
tmux send-keys -t chainb "var tokenAbi = $TokenAbi" C-m
tmux send-keys -t chainb "var tokenContract = eth.contract(tokenAbi)" C-m
tmux send-keys -t chainb "var token = tokenContract.at(tokenAddr)" C-m

## Print the running tmux windows
tmux ls

