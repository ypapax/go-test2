#!/usr/bin/env bash
set -ex
cd $GOPATH/src/github.com/ypapax/go-test2/apps/go-test2
go install
go-test2 -alsologtostderr -conn $@