#!/bin/sh
ME=`realpath $0`
`dirname $ME`/../node_modules/.bin/truffle debug $1
