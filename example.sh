#!/bin/bash

ACC_A=0x22d491bde2303f2f43325b2108d26f1eaba1e32b
ACC_B=0xffcf8fdee72ac11b5c542428b35eef5769c409f0
TOKEN_ADDR=0x254dffcd3277c0b1660f6d42efbb754edababc2b
LOCK_ADDR=0xc89ce4735882c9f0f0fe26686c53074e09b0d550
LINK_ADDR=0xcfeb869f69431e42cdb54a4f4f105c19c080a601
PORT_A=8545
PORT_B=8546
IP_A=127.0.0.1
IP_B=127.0.0.1
API_PORT_A=5000
API_PORT_B=5001

echo "==== Chain A ===="
echo "...Minting"
python -mion ion mint --rpc $IP_A:$PORT_A --account $ACC_A --tkn $TOKEN_ADDR --value 5000
echo ""
echo "Press any key to proceed"
read enter

echo "...Depositing"
python -mion ion deposit --rpc $IP_A:$PORT_A --account $ACC_A --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
echo ""
echo "Press any key to proceed"
read enter

echo "...Fetching proof"
python -mion ion proof --lithium-port $API_PORT_A --account $ACC_A --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
echo ""
echo "Press any key to proceed"
read enter

echo "==== Chain B ===="
echo "...Minting"
python -mion ion mint --rpc $IP_B:$PORT_B --account $ACC_B --tkn $TOKEN_ADDR --value 5000
echo ""
echo "Press any key to proceed"
read enter

echo "...Depositing"
python -mion ion deposit --rpc $IP_B:$PORT_B --account $ACC_B --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
echo ""
echo "Press any key to proceed"
read enter

echo "...Fetching proof"
python -mion ion proof --lithium-port $API_PORT_B --account $ACC_B --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
echo ""
echo "Press any key to proceed"
read enter

echo "==== Withdrawing from Chain A ===="
python -mion ion withdraw --lithium-port $API_PORT_B --rpc $IP_A:$PORT_A --account $ACC_B --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
echo ""
echo "Press any key to proceed"
read enter

echo "==== Withdrawing from Chain B ===="
python -mion ion withdraw --lithium-port $API_PORT_A --rpc $IP_B:$PORT_B --account $ACC_A --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
