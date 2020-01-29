#!/bin/bash

# run each file test separately and his benchmark right after 
# this because test rpc delete old tx when too many 
# so benchamrk won't work if run once at the end of all the tests

# TODO avoid contracts getting compiled at each round 

COMMAND="$1" # trace, compare, toMD

runComparison() {
    if [ ! -z $COMPARE_A ] && [ ! -z $COMPARE_B ]; then
      node ./node_modules/benchmark-solidity/benchmark.js compare $COMPARE_B $COMPARE_A
    else
      printf "Please provide the two benchmark files absolute paths you intend to compare"
    fi   
}

convertToMd() {
  if [ ! -z $JSON_INPUT ] && [ ! -z $MD_OUTPUT ]; then
    node ./node_modules/benchmark-solidity/benchmark.js convertToMD "$JSON_INPUT" "$MD_OUTPUT"
    printf ""$JSON_INPUT" file converted into $MD_OUTPUT"
  else
    printf "Please provide the input and output file paths for the conversion"
  fi  
}

start() {
    for entry in ./test/*.js
    do  
        printf "Testing $entry..\n"
        npm run test "$entry" || ( printf "\nAn error has occurred running the test. Quitting\n"; exit )
        
        printf "Tracing transactions in $entry..\n"
        node ./node_modules/benchmark-solidity/benchmark.js trace $CONFIGS
    done

    printf "All tests have been benchmarked - Converting now into Markdown"

    convertToMd
}

trap "kill 0" EXIT

if [ "$COMMAND" == "start" ]; then 

  JSON_INPUT="$2" # json input path of the file to be converted
  MD_OUTPUT="$3" # markdown path of the output file
  CONFIGS="$4" # optional config object

  start

elif [ "$COMMAND" == "toMD" ]; then

  JSON_INPUT="$2"
  MD_OUTPUT="$3" 

  convertToMd

elif [ "$COMMAND" == "compare" ]; then

  COMPARE_A="$3" # after-changes benchmark file
  COMPARE_B="$2" # before-changes benchmark file 

  runComparison

else 
  printf "Command not recognised. Should be either trace, compare or toMD"
fi