#!/usr/bin/env bash
set -ex
onexit(){
    if [ $? -eq 0 ]; then
        echo successfull deploy
    else
        echo failed deploy
    fi
}
trap onexit EXIT
user=$1
if [ -z $user ]; then
    >&2 echo missing user
    exit 1
fi
host=$2
if [ -z $host ]; then
    >&2 echo missing host
    exit 1
fi
userHost=${user}@${host}
port=$3
if [ -z $port ]; then
    >&2 echo missing port
    exit 1
fi
connStr=$4
if [ -z $connStr ]; then
    >&2 echo missing connStr
    exit 1
fi
serviceName=go-test2
binLocal=/tmp/$serviceName
pushd $GOPATH/src/github.com/ypapax/$serviceName
pushd apps/$serviceName
go get ./...
GOOS=linux GOARCH=amd64 go build -o $binLocal
popd
remoteDir='~/go/go-test2'
ssh $userHost mkdir -p $remoteDir
zipDir=/tmp/${serviceName}-deploy
rm -rf $zipDir
mkdir -p $zipDir
zipName=${serviceName}.zip
cp $binLocal daemon.sh $zipDir
pushd $zipDir
zip $zipName *
popd
scp $zipDir/$zipName prod.sh $userHost:$remoteDir
popd
ssh $userHost ZIP_NAME=$zipName  $remoteDir/prod.sh -conn $connStr -v 4 -port $port
set +e
ssh $userHost pkill -ef go-test2
echo

while :; do
    if curl http://$host:$port/test/api/v1/temperature/avg; then
        exit 0
    fi
    sleep 1
done
