#!/bin/bash

DIR=$PWD
CMD=../cmd

# Kill all linotwit-* stuff
function cleanup {
	pkill linotwit
}

cd $CMD/agent
# Add `edgex-` prefix on start, so we can find the process family
exec -a linotwit-agent ./agent &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done