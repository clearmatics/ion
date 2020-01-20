#!/bin/bash

# run each file test separately and his benchmark right after 
# otherwise test rpc delete old tx when too many and benchmark won't work

# TODO avoid contracts getting compiled at each round 

NETWORK=petersburgRpc
if [ ! -z "$1" ]; then 
    NETWORK="$1"
fi 

# COMMAND="$1" # compare || run&compare
COMPARE_A="$2" # before-changes benchmark file
COMPARE_B="$3" # after-changes benchmark file

runComparison() {
    
    if [ ! -z COMPARE_A ] && [ ! -z COMPARE_B ]; then
      node ./benchmark/benchmark.js compare $COMPARE_A $COMPARE_B
    fi  
}

runTests() {
    for entry in ./test/*.js
    do
        npm run test "$entry" || ( printf "\nPlease provide a script name that runs an rpc network from package.json\n"; exit )
        node ./benchmark/benchmark.js trace 
    done

    runComparison
    
}

trap "kill 0" EXIT

# run testrpc and truffle tests, benchmarking them
npm run $NETWORK > /dev/null & runComparison 