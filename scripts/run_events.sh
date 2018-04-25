#!/bin/bash


CHAINA=127.0.0.1:8545
CHAINB=127.0.0.1:8546
SEND=0xffcf8fdee72ac11b5c542428b35eef5769c409f0
RECV=0x22d491bde2303f2f43325b2108d26f1eaba1e32b
LOCK=0xe982e462b094850f12af94d21d470e21be9d0e9c
LINK=0xc89ce4735882c9f0f0fe26686c53074e09b0d550

## Run the event listeners
tmux new -s lista -n lista -d
tmux send-keys -t lista "python -mion etheventrelay --rpc-from $CHAINA --rpc-to $CHAINB --from-account $SEND --to-account $RECV --lock $LOCK --link $LINK" C-m

tmux new -s listb -n listb -d
tmux send-keys -t listb "python -mion etheventrelay --rpc-from $CHAINB --rpc-to $CHAINA --from-account $RECV --to-account $SEND --lock $LOCK --link $LINK" C-m
## Print the running tmux windows
tmux ls

