#!/bin/bash

# run each file test separately and his benchmark right after 
# otherwise test rpc delete old tx when too many and benchmark won't work

# TODO avoid contracts getting compiled at each round 

for entry in ./test/*.js
do
  npm run test "$entry"
  node ./benchmark.js trace
done

