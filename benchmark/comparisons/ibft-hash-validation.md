# Metadata

## Before


 - title: ION
 - network: ganacheRPC-Istanbul
 - blockTime: 1s

## After


 - title: ION
 - network: ganacheRPC-Istanbul
 - blockTime: 0s
 - difference: ibft validation stores hash of sorted array of validators instead of whole set

 | name | gas before |  gas after | % gas improvement | seconds before | seconds after | % time improvement | 
 | --- | --- | --- | --- | --- | --- | --- | 
 | Deploy Ion | 606641 | 606641 | 0.000 | Not estimated | Not estimated | NaN
Deploy Ethereum Store | 1892713 | 1892725 | -0.001 | Not estimated | Not estimated | NaN
Deploy Clique Validation | 2317223 | 2317235 | -0.001 | Not estimated | Not estimated | NaN
Deploy Event verifier | 614356 | 614356 | 0.000 | Not estimated | Not estimated | NaN
Deploy Event Function | 394823 | 394835 | -0.003 | Not estimated | Not estimated | NaN
clique-Successful Register Chain | 374573 | 374561 | 0.003 | 0.213 | 0.208 | 2.347
clique-Authentic Submission Happy Path | 448363 | 448375 | -0.003 | 0.519 | 0.523 | -0.771
clique-Add Validators Through Block Submission | 497888 | 497888 | 0.000 | 0.580 | 0.437 | 24.655
ibft-Successful Register Chain | 310840 | 210011 | 32.438 | 0.168 | 0.121 | 27.976
ibft-Successful Submit block - 4 validators  | 291777 | 229373 | 21.388 | 0.550 | 0.529 | 3.818
ibft-Successful Submit block - 5 validators  | 310962 | 235986 | 24.111 | 0.478 | 0.469 | 1.883
ibft-Successful Submit block - 8 validators  | 406651 | 272282 | 33.043 | 0.769 | 0.666 | 13.394
ibft-Successful Submit block - 16 validators  | 731773 | 387105 | 47.100 | 1.194 | 0.945 | 20.854
integration-fabric-Successful Register Chain | 120404 | 120416 | -0.010 | 0.092 | 0.090 | 2.174
integration-fabric-Successful Add Block | 1269590 | 1269590 | 0.000 | 1.020 | 1.077 | -5.588
integration-fabric-Submit Block, retrieve state and execute | 38284 | 38284 | 0.000 | 0.086 | 0.090 | -4.651
integration-clique-Successful Register Chain | 376315 | 376315 | 0.000 | 0.163 | 0.122 | 25.153
integration-clique-Successful Add Block | 537929 | 537929 | 0.000 | 0.965 | 0.618 | 35.959
integration-clique-Successful Check Proofs | 247084 | 247084 | 0.000 | 1.082 | 1.064 | 1.664
integration-clique-Submit Block, verify tx and event, and execute | 297974 | 297974 | 0.000 | 1.232 | 1.318 | -6.981
Ion contract-Successful Store Block | 41957 | 41957 | 0.000 | 0.110 | 0.096 | 12.727
storage-ethereum-Successful Check All Proofs | 247084 | 247084 | 0.000 | 1.020 | 1.108 | -8.627
storage-fabric-Successful Add Block | 1262833 | 1262833 | 0.000 | 0.972 | 0.940 | 3.292 | 
