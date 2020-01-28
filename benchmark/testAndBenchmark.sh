#!/bin/bash

# run each file test separately and his benchmark right after 
# this because test rpc delete old tx when too many 
# so benchamrk won't work if run once at the end of all the tests

# TODO avoid contracts getting compiled at each round 

COMMAND="$1" # either trace or compare
COMPARE_A="$2" # before-changes benchmark file
COMPARE_B="$3" # after-changes benchmark file
CONFIGS="$4" # configs 

runComparison() {
    
    if [ ! -z $COMPARE_A ] && [ ! -z $COMPARE_B ]; then
      node ./node_modules/benchmark-solidity/benchmark.js compare $COMPARE_A $COMPARE_B
    else
      printf "Please provide the two benchmark files absolute paths you intend to compare"
    fi   
}

runTestsAndTraceTx() {
    for entry in ./test/*.js
    do
        npm run test "$entry" || ( printf "\nAn error has occurred running the test. Quitting\n"; exit )
        node ./node_modules/benchmark-solidity/benchmark.js trace $CONFIGS
    done

    runComparison   
}

trap "kill 0" EXIT

if [[ -z $COMMAND  ||  "$COMMAND" == "trace" ]]; then 
  # default
  runTestsAndTraceTx
elif [ "$COMMAND" == "compare" ]; then
  # input files must have also been provided
  runComparison
else 
  printf "Command not recognised. Should be either trace or compare"
fi