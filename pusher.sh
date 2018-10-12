#!/bin/sh

set -eu

docker tag $1 docker.for.mac.localhost:5000/$1

ssh -NT -L 5000:localhost:5000 core@10.0.87.78 &
ssh_pid=$!

docker push docker.for.mac.localhost:5000/$1

kill $ssh_pid
