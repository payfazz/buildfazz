#!/bin/sh

set -eux

cd "$(dirname "$0")"

if [ -z "$GOPATH" ]; then
	export GOPATH=~/go
fi

docker build -t a1:b2 .

dangling_docker=$(docker images -f 'dangling=true' -q)
if [ -z "$dangling_docker" ]; then
    exit 1
fi

docker rmi $dangling_docker --force