# Metadata

## Before


 - title: ION
 - network: ganacheRPC-Istanbul
 - blockTime: 1s

## After


 - title: ION
 - network: ganacheRPC-Istanbul
 - blockTime: 1s
 - difference: ibft validation stores hash of sorted array of validators instead of whole set

 | name | gas before |  gas after | % gas improvement | seconds before | seconds after | % time improvement | 
 | --- | --- | --- | --- | --- | --- | --- | 
 | Deploy Ion | 606641 | 606641 | 0.000 | Not estimated | Not estimated | NaN
Deploy Ethereum Store | 1892725 | 1892725 | 0.000 | Not estimated | Not estimated | NaN
Deploy Clique Validation | 2317235 | 2317235 | 0.000 | Not estimated | Not estimated | NaN
Deploy Event verifier | 614356 | 614356 | 0.000 | Not estimated | Not estimated | NaN
Deploy Event Function | 394835 | 394835 | 0.000 | Not estimated | Not estimated | NaN
clique-Successful Register Chain | 374573 | 374573 | 0.000 | 1.063 | 1.060 | 0.282
clique-Authentic Submission Happy Path | 448375 | 448375 | 0.000 | 1.208 | 1.700 | -40.728
clique-Add Validators Through Block Submission | 497888 | 497888 | 0.000 | 1.157 | 1.482 | -28.090
ibft-Successful Register Chain | 310828 | 209874 | 32.479 | 1.050 | 1.031 | 1.810
ibft-Successful Submit block | 285422 | 233838 | 18.073 | 1.140 | 1.228 | -7.719
ibft-Submit Sequential Blocks with Additional Validator | 296160 | 215527 | 27.226 | 1.512 | 1.455 | 3.770
integration-fabric-Successful Register Chain | 120416 | 120416 | 0.000 | 1.022 | 1.089 | -6.556
integration-fabric-Successful Add Block | 1269590 | 1269590 | 0.000 | 1.723 | 1.866 | -8.299
integration-fabric-Submit Block, retrieve state and execute | 38284 | 38284 | 0.000 | 1.388 | 1.341 | 3.386
integration-clique-Successful Register Chain | 376315 | 376315 | 0.000 | 1.040 | 1.030 | 0.962
integration-clique-Successful Add Block | 537929 | 537929 | 0.000 | 1.099 | 1.891 | -72.066
integration-clique-Successful Check Proofs | 247084 | 247084 | 0.000 | 1.724 | 2.307 | -33.817
integration-clique-Submit Block, verify tx and event, and execute | 297974 | 297974 | 0.000 | 2.132 | 2.244 | -5.253
Ion contract-Successful Store Block | 41957 | 41957 | 0.000 | 1.038 | 1.136 | -9.441
storage-ethereum-Successful Check All Proofs | 247084 | 247084 | 0.000 | 2.075 | 2.828 | -36.289
storage-fabric-Successful Add Block | 1262833 | 1262833 | 0.000 | 1.531 | 2.564 | -67.472 | 
