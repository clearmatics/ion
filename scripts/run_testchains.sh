#!/bin/bash

## Run and attach network 0
tmux new -s chaina -n testrpc_migration_a -d
tmux send-keys -t chaina "npm run testrpca" C-m
tmux split-window -h -t chaina:testrpc_migration_a
tmux send-keys -t chaina "rm -rf build/contracts/*" C-m
tmux send-keys -t chaina "truffle compile" C-m
tmux send-keys -t chaina "truffle migrate --network testrpca" C-m
tmux send-keys -t chaina "geth attach http://127.0.0.1:8545" C-m

## Run and attach network 1
tmux new -s chainb -n testrpc_migration_b -d
tmux send-keys -t chainb "npm run testrpcb" C-m
tmux split-window -h -t chainb:testrpc_migration_b
tmux send-keys -t chainb "truffle migrate --network testrpcb" C-m
tmux send-keys -t chainb "geth attach http://127.0.0.1:8546" C-m

## Print the running tmux windows
tmux ls

