!#/bin/bash
geth --datadir docker_build/account/ --syncmode 'full' --port 30311 --rpc --rpcaddr '0.0.0.0' --rpcport 8545 --networkid 1515 --gasprice '0' --targetgaslimit 0xFFFFFFFFFFFF --unlock '0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d' --password docker_build/account/password-2be5ab0e43b6dc2908d5321cf318f35b80d0c10d.txt --mine
