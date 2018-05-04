#!/bin/bash
echo "==== Chain A ===="
echo "Minting"
python -mion ion mint --rpc 127.0.0.1:8545 --account 0x22d491bde2303f2f43325b2108d26f1eaba1e32b --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000
read enter

echo "Depositing"
python -mion ion deposit --rpc 127.0.0.1:8545 --account 0x22d491bde2303f2f43325b2108d26f1eaba1e32b --lock 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000 --ref stuff
read enter

echo "Fetching proof"
python -mion ion proof --rpc 127.0.0.1:8545 --account 0x22d491bde2303f2f43325b2108d26f1eaba1e32b --lock 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000 --ref stuff
read enter

echo "==== Chain B ===="
echo "Minting"
python -mion ion mint --rpc 127.0.0.1:8546 --account 0xffcf8fdee72ac11b5c542428b35eef5769c409f0 --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000
read enter

echo "Depositing"
python -mion ion deposit --rpc 127.0.0.1:8546 --account 0xffcf8fdee72ac11b5c542428b35eef5769c409f0 --lock 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000 --ref stuff
read enter

echo "Fetching proof"
python -mion ion proof --rpc 127.0.0.1:8546 --account 0xffcf8fdee72ac11b5c542428b35eef5769c409f0 --lock 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000 --ref stuff
read enter

echo "==== Withdrawing from Chain A ===="
python -mion ion withdraw --rpc-from 127.0.0.1:8546 --rpc-to 127.0.0.1:8545 --account 0xffcf8fdee72ac11b5c542428b35eef5769c409f0 --lock 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000 --ref stuff
read enter

echo "==== Withdrawing from Chain B ===="
python -mion ion withdraw --rpc-from 127.0.0.1:8545 --rpc-to 127.0.0.1:8546 --account 0x22d491bde2303f2f43325b2108d26f1eaba1e32b --lock 0xd833215cbcc3f914bd1c9ece3ee7bf8b14f841bb --tkn 0xc89ce4735882c9f0f0fe26686c53074e09b0d550 --value 5000 --ref stuff