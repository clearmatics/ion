#!/bin/bash


CHAINA=127.0.0.1:8545
CHAINB=127.0.0.1:8546
SEND=0xffcf8fdee72ac11b5c542428b35eef5769c409f0
RECV=0x22d491bde2303f2f43325b2108d26f1eaba1e32b
LOCK=0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb
LINK=0xcfeb869f69431e42cdb54a4f4f105c19c080a601

## Run the event listeners
tmux new -s lista -n lista -d
tmux send-keys -t lista "python -mion etheventrelay --rpc-from $CHAINA --rpc-to $CHAINB --from-account $SEND --to-account $RECV --lock $LOCK --link $LINK" C-m

tmux new -s listb -n listb -d
tmux send-keys -t listb "python -mion etheventrelay --rpc-from $CHAINB --rpc-to $CHAINA --from-account $RECV --to-account $SEND --lock $LOCK --link $LINK" C-m
## Print the running tmux windows
tmux ls
