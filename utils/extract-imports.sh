#!/bin/sh
cat $1 | grep import | sed -e 's/import "\.\//import "/g' | cut -f 2 -d '"' | xargs echo
