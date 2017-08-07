#!/usr/bin/env bash
set -ex
cd $GOPATH/src/github.com/ypapax/go-test2
./deploy.sh $@
host=$2
port=$3
./check.sh http://$host:$port